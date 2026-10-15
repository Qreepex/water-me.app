package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qreepex/water-me-app/backend/types"

	"firebase.google.com/go/messaging"
)

const (
	NotificationCooldown = 4 * time.Hour // Don't spam users more than once per 4 hours
	FCMBatchSize         = 500           // FCM allows max 500 tokens per multicast
)

type NotificationStats struct {
	PlantsChecked       int
	NotificationsSent   int
	NotificationsFailed int
	UsersNotified       int
}

type NotificationMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type NotificationMessages struct {
	Single   []NotificationMessage `json:"single"`
	Multiple []NotificationMessage `json:"multiple"`
}

var (
	messageCache = make(map[string]*NotificationMessages)
	rng          = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// LoadNotificationMessages loads messages from JSON files
func LoadNotificationMessages(messageDir string) error {
	types := []string{"watering", "moisture", "fertilizing", "misting", "repotting"}

	for _, msgType := range types {
		filePath := filepath.Join(messageDir, msgType+".json")
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s.json: %w", msgType, err)
		}

		var messages NotificationMessages
		if err := json.Unmarshal(data, &messages); err != nil {
			return fmt.Errorf("failed to parse %s.json: %w", msgType, err)
		}

		messageCache[msgType] = &messages
		log.Printf("Loaded %d single and %d multiple messages for %s",
			len(messages.Single), len(messages.Multiple), msgType)
	}

	return nil
}

// ProcessNotifications processes all notification types and returns statistics
func ProcessNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	plantsBatchSize int,
) *NotificationStats {
	stats := &NotificationStats{}

	// Process each notification type
	processWateringNotifications(ctx, db, firebase, plantsBatchSize, stats)
	processMoistureNotifications(ctx, db, firebase, plantsBatchSize, stats)
	processFertilizingNotifications(ctx, db, firebase, plantsBatchSize, stats)
	processMistingNotifications(ctx, db, firebase, plantsBatchSize, stats)
	processRepottingNotifications(ctx, db, firebase, plantsBatchSize, stats)

	return stats
}

func processWateringNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	batchSize int,
	stats *NotificationStats,
) {
	plants, err := db.GetPlantsNeedingWatering(ctx, batchSize)
	if err != nil {
		log.Printf("Error fetching plants needing watering: %v", err)
		return
	}

	if len(plants) == 0 {
		log.Println("No plants need watering")
		return
	}

	log.Printf("Found %d plants needing watering", len(plants))
	stats.PlantsChecked += len(plants)

	userPlants := groupPlantsByUser(plants)
	for userID, userPlantList := range userPlants {
		sendNotificationsForUser(ctx, db, firebase, userID, userPlantList, "watering", stats)
	}
}

// processMoistureNotifications reminds users whose plants use the moisture
// based plan. Two things can be due: watering the plant, or taking a fresh
// measurement so the plan can be recalculated.
func processMoistureNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	batchSize int,
	stats *NotificationStats,
) {
	plants, err := db.GetMoisturePlants(ctx, batchSize)
	if err != nil {
		log.Printf("Error fetching moisture plants: %v", err)
		return
	}
	if len(plants) == 0 {
		return
	}

	now := time.Now()
	toWater := make([]types.Plant, 0, len(plants))
	toMeasure := make([]types.Plant, 0, len(plants))
	for _, plant := range plants {
		plan := ComputeWateringPlan(plant, now)
		if plan == nil {
			continue
		}
		switch plan.Status {
		case types.PlanStatusWaterNow:
			toWater = append(toWater, plant)
		case types.PlanStatusMeasure:
			// Only nudge for a measurement once it is actually due.
			if plan.MeasurementDueAt == nil || !now.Before(*plan.MeasurementDueAt) {
				toMeasure = append(toMeasure, plant)
			}
		}
	}

	stats.PlantsChecked += len(plants)

	if len(toWater) > 0 {
		log.Printf("Found %d moisture plants needing watering", len(toWater))
		for userID, userPlants := range groupPlantsByUser(toWater) {
			sendNotificationsForUser(ctx, db, firebase, userID, userPlants, "watering", stats)
		}
	}

	if len(toMeasure) > 0 {
		log.Printf("Found %d moisture plants needing a measurement", len(toMeasure))
		for userID, userPlants := range groupPlantsByUser(toMeasure) {
			sendNotificationsForUser(ctx, db, firebase, userID, userPlants, "moisture", stats)
		}
	}
}

func processFertilizingNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	batchSize int,
	stats *NotificationStats,
) {
	plants, err := db.GetPlantsNeedingFertilizer(ctx, batchSize)
	if err != nil {
		log.Printf("Error fetching plants needing fertilizer: %v", err)
		return
	}

	if len(plants) == 0 {
		return
	}

	log.Printf("Found %d plants needing fertilizer", len(plants))
	stats.PlantsChecked += len(plants)

	userPlants := groupPlantsByUser(plants)
	for userID, userPlantList := range userPlants {
		sendNotificationsForUser(ctx, db, firebase, userID, userPlantList, "fertilizing", stats)
	}
}

func processMistingNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	batchSize int,
	stats *NotificationStats,
) {
	plants, err := db.GetPlantsNeedingMisting(ctx, batchSize)
	if err != nil {
		log.Printf("Error fetching plants needing misting: %v", err)
		return
	}

	if len(plants) == 0 {
		return
	}

	log.Printf("Found %d plants needing misting", len(plants))
	stats.PlantsChecked += len(plants)

	userPlants := groupPlantsByUser(plants)
	for userID, userPlantList := range userPlants {
		sendNotificationsForUser(ctx, db, firebase, userID, userPlantList, "misting", stats)
	}
}

func processRepottingNotifications(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	batchSize int,
	stats *NotificationStats,
) {
	plants, err := db.GetPlantsNeedingRepotting(ctx, batchSize)
	if err != nil {
		log.Printf("Error fetching plants needing repotting: %v", err)
		return
	}

	if len(plants) == 0 {
		return
	}

	log.Printf("Found %d plants needing repotting", len(plants))
	stats.PlantsChecked += len(plants)

	userPlants := groupPlantsByUser(plants)
	for userID, userPlantList := range userPlants {
		sendNotificationsForUser(ctx, db, firebase, userID, userPlantList, "repotting", stats)
	}
}

func groupPlantsByUser(plants []types.Plant) map[string][]types.Plant {
	grouped := make(map[string][]types.Plant)
	for _, plant := range plants {
		grouped[plant.UserID] = append(grouped[plant.UserID], plant)
	}
	return grouped
}

func sendNotificationsForUser(
	ctx context.Context,
	db *MongoDB,
	firebase *FirebaseService,
	userID string,
	plants []types.Plant,
	notificationType string,
	stats *NotificationStats,
) {
	// Get user's notification config
	config, err := db.GetNotificationConfig(ctx, userID)
	if err != nil || config == nil {
		return
	}

	// Check if notifications are enabled
	if !config.IsEnabled {
		return
	}

	// Check notification type settings
	if !isNotificationTypeEnabled(config, notificationType) {
		return
	}

	// Check cooldown - don't spam users
	if config.LastNotificationSentAt != nil {
		timeSinceLastNotification := time.Since(*config.LastNotificationSentAt)
		if timeSinceLastNotification < NotificationCooldown {
			return
		}
	}

	// Filter out muted plants
	notifyPlants := filterMutedPlants(plants, config.MutedPlantIDs)
	if len(notifyPlants) == 0 {
		return
	}

	// Get active device tokens
	activeTokens := getActiveTokens(config.DeviceTokens)
	if len(activeTokens) == 0 {
		return
	}

	// Build notification message
	title, body := buildNotificationMessage(notifyPlants, notificationType)

	// Send notifications in batches
	failedTokens := sendNotificationBatches(
		ctx,
		firebase,
		activeTokens,
		title,
		body,
		notificationType,
		len(notifyPlants),
		stats,
	)

	// Handle failed tokens
	if len(failedTokens) > 0 {
		if err := db.MarkTokensAsInactive(ctx, userID, failedTokens); err != nil {
			log.Printf("Error marking tokens as inactive for user %s: %v", userID, err)
		}
	}

	// Update last notification sent timestamp
	if err := db.UpdateNotificationLastSent(ctx, userID); err != nil {
		log.Printf("Error updating last notification sent for user %s: %v", userID, err)
	}

	stats.UsersNotified++
	log.Printf(
		"Sent %s notifications to user %s for %d plants",
		notificationType,
		userID,
		len(notifyPlants),
	)
}

func isNotificationTypeEnabled(config *types.NotificationConfig, notificationType string) bool {
	switch notificationType {
	case "watering":
		return config.RemindWatering
	case "fertilizing":
		return config.RemindFertilize
	case "misting":
		return config.RemindMisting
	case "moisture":
		return config.RemindMoisture
	case "repotting":
		return config.RemindRepotting
	default:
		return false
	}
}

func filterMutedPlants(plants []types.Plant, mutedPlantIDs []string) []types.Plant {
	mutedPlantMap := make(map[string]bool)
	for _, plantID := range mutedPlantIDs {
		mutedPlantMap[plantID] = true
	}

	var notifyPlants []types.Plant
	for _, plant := range plants {
		if !mutedPlantMap[plant.ID] {
			notifyPlants = append(notifyPlants, plant)
		}
	}
	return notifyPlants
}

func getActiveTokens(deviceTokens []types.DeviceToken) []string {
	var activeTokens []string
	for _, device := range deviceTokens {
		if device.IsActive {
			activeTokens = append(activeTokens, device.Token)
		}
	}
	return activeTokens
}

func sendNotificationBatches(
	ctx context.Context,
	firebase *FirebaseService,
	tokens []string,
	title, body, notificationType string,
	plantCount int,
	stats *NotificationStats,
) []string {
	var failedTokens []string

	// Send notifications in batches of 500 (FCM limit)
	for i := 0; i < len(tokens); i += FCMBatchSize {
		end := min(i+FCMBatchSize, len(tokens))

		tokenBatch := tokens[i:end]

		// Send multicast notification
		data := map[string]string{
			"type":       notificationType,
			"plantCount": fmt.Sprintf("%d", plantCount),
		}

		response, err := firebase.SendMulticastNotification(ctx, tokenBatch, title, body, data)
		if err != nil {
			log.Printf("Error sending notification batch: %v", err)
			stats.NotificationsFailed += len(tokenBatch)
			continue
		}

		// Track success/failure
		stats.NotificationsSent += response.SuccessCount
		stats.NotificationsFailed += response.FailureCount

		// Collect failed tokens
		if response.FailureCount > 0 {
			failed := extractFailedTokens(tokenBatch, response)
			failedTokens = append(failedTokens, failed...)
		}
	}

	return failedTokens
}

func extractFailedTokens(tokens []string, response *messaging.BatchResponse) []string {
	var failed []string
	for i, sendResponse := range response.Responses {
		if !sendResponse.Success {
			if i < len(tokens) {
				failed = append(failed, tokens[i])
				// Log the error for debugging
				if sendResponse.Error != nil {
					log.Printf(
						"Token failed: %s, Error: %v",
						tokens[i][:20]+"...",
						sendResponse.Error,
					)
				}
			}
		}
	}
	return failed
}

func buildNotificationMessage(plants []types.Plant, notificationType string) (string, string) {
	messages, ok := messageCache[notificationType]
	if !ok || messages == nil {
		// Fallback to default messages if cache is empty
		return buildDefaultMessage(plants, notificationType)
	}

	count := len(plants)

	if count == 1 {
		// Single plant - choose random message from single variants
		if len(messages.Single) == 0 {
			return buildDefaultMessage(plants, notificationType)
		}
		msg := messages.Single[rng.Intn(len(messages.Single))]
		title := strings.ReplaceAll(msg.Title, "{plantName}", plants[0].Name)
		body := strings.ReplaceAll(msg.Body, "{plantName}", plants[0].Name)
		return title, body
	}

	// Multiple plants - choose random message from multiple variants
	if len(messages.Multiple) == 0 {
		return buildDefaultMessage(plants, notificationType)
	}
	msg := messages.Multiple[rng.Intn(len(messages.Multiple))]

	// Build plant names list
	plantNames := formatPlantNames(plants)

	title := strings.ReplaceAll(msg.Title, "{count}", fmt.Sprintf("%d", count))
	title = strings.ReplaceAll(title, "{plantNames}", plantNames)
	title = strings.ReplaceAll(title, "{remaining}", fmt.Sprintf("%d", count-2))

	body := strings.ReplaceAll(msg.Body, "{count}", fmt.Sprintf("%d", count))
	body = strings.ReplaceAll(body, "{plantNames}", plantNames)
	body = strings.ReplaceAll(body, "{remaining}", fmt.Sprintf("%d", count-2))

	return title, body
}

func formatPlantNames(plants []types.Plant) string {
	if len(plants) <= 3 {
		var names strings.Builder
		for i, plant := range plants {
			if i > 0 {
				names.WriteString(", ")
			}
			names.WriteString(plant.Name)
		}
		return names.String()
	}
	// For more than 3 plants, show first 2
	return fmt.Sprintf("%s, %s", plants[0].Name, plants[1].Name)
}

func buildDefaultMessage(plants []types.Plant, notificationType string) (string, string) {
	count := len(plants)
	var title, body string

	switch notificationType {
	case "watering":
		if count == 1 {
			title = "Time to water your plant! 💧"
			body = fmt.Sprintf("%s needs water", plants[0].Name)
		} else {
			title = fmt.Sprintf("%d plants need water! 💧", count)
			body = formatPlantNames(plants)
		}

	case "fertilizing":
		if count == 1 {
			title = "Time to fertilize! 🌱"
			body = fmt.Sprintf("%s needs fertilizer", plants[0].Name)
		} else {
			title = fmt.Sprintf("%d plants need fertilizer! 🌱", count)
			body = "Don't forget to feed your plants"
		}

	case "misting":
		if count == 1 {
			title = "Misting time! 💦"
			body = fmt.Sprintf("%s needs misting", plants[0].Name)
		} else {
			title = fmt.Sprintf("%d plants need misting! 💦", count)
			body = "Keep your plants humid and happy"
		}

	case "repotting":
		if count == 1 {
			title = "Repotting reminder! 🪴"
			body = fmt.Sprintf("%s might need repotting", plants[0].Name)
		} else {
			title = fmt.Sprintf("%d plants might need repotting! 🪴", count)
			body = "Check if your plants need fresh soil"
		}

	case "moisture":
		if count == 1 {
			title = "Measure your soil moisture 📏"
			body = fmt.Sprintf("Take a new reading for %s", plants[0].Name)
		} else {
			title = fmt.Sprintf("Measure %d plants 📏", count)
			body = "New readings keep the watering plan accurate"
		}
	}

	return title, body
}
