package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/qreepex/water-me-app/backend/services"
	"github.com/qreepex/water-me-app/backend/types"
	"github.com/qreepex/water-me-app/backend/util"
	"github.com/qreepex/water-me-app/backend/validation"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func NotificationHandler(router *mux.Router, database *services.MongoDB) {
	router.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		getNotificationConfig(w, r, database)
	}).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		upsertNotificationConfig(w, r, database)
	}).Methods(http.MethodPut, http.MethodOptions)

	router.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		deleteNotificationConfig(w, r, database)
	}).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/api/notifications/tokens", func(w http.ResponseWriter, r *http.Request) {
		registerDeviceToken(w, r, database)
	}).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/notifications/tokens/{deviceId}", func(w http.ResponseWriter, r *http.Request) {
		deleteDeviceToken(w, r, database)
	}).
		Methods(http.MethodDelete, http.MethodOptions)
}

func getNotificationConfig(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Getting notification config for user %v", userID)

	config, err := db.GetNotificationConfig(r.Context(), userID)
	if err != nil {
		if err == types.ErrNoDocuments {
			// Return default config if none exists
			defaultConfig := createDefaultConfig(userID)
			util.RespondJSON(w, http.StatusOK, defaultConfig)
			return
		}
		log.Printf("Failed to retrieve notification config: %v", err)
		http.Error(w, "Failed to retrieve notification config", http.StatusInternalServerError)
		return
	}

	util.RespondJSON(w, http.StatusOK, config)
}

func upsertNotificationConfig(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var config types.NotificationConfig
	if err := util.DecodeJSON(r, &config); err != nil {
		util.BadRequest(w, err.Error(), nil)
		return
	}

	// Validate the config
	errors := validation.ValidateNotificationConfig(config)
	if len(errors) > 0 {
		util.BadRequest(w, "Validation failed", errors)
		return
	}

	// Verify that all muted plant IDs belong to the user
	if len(config.MutedPlantIDs) > 0 {
		userPlants, err := db.GetPlants(r.Context(), userID)
		if err != nil {
			util.ServerError(w, err)
			return
		}

		// Create a map of valid plant IDs for this user
		validPlantIDs := make(map[string]bool)
		for _, plant := range userPlants {
			validPlantIDs[plant.ID] = true
		}

		// Check that all muted plant IDs are valid
		invalidIDs := []string{}
		for _, plantID := range config.MutedPlantIDs {
			if !validPlantIDs[plantID] {
				invalidIDs = append(invalidIDs, plantID)
			}
		}

		if len(invalidIDs) > 0 {
			util.BadRequest(w, "Invalid plant IDs", []types.ValidationError{
				{
					Field:   "mutedPlantIds",
					Message: "One or more plant IDs do not belong to this user or do not exist",
				},
			})
			return
		}
	}

	// Set user ID and timestamps
	config.UserID = userID
	config.UpdatedAt = time.Now()

	// Check if config exists
	existing, err := db.GetNotificationConfig(r.Context(), userID)
	if err != nil && err != types.ErrNoDocuments {
		util.ServerError(w, err)
		return
	}

	if existing != nil {
		// Update existing config
		config.ID = existing.ID
		updatedConfig, err := db.UpdateNotificationConfig(r.Context(), config)
		if err != nil {
			util.ServerError(w, err)
			return
		}
		util.RespondJSON(w, http.StatusOK, updatedConfig)
	} else {
		// Create new config
		id, err := gonanoid.New()
		if err != nil {
			log.Printf("Failed to generate notification config ID: %v", err)
			http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
			return
		}
		config.ID = id

		createdConfig, err := db.CreateNotificationConfig(r.Context(), config)
		if err != nil {
			util.ServerError(w, err)
			return
		}
		util.RespondJSON(w, http.StatusCreated, createdConfig)
	}
}

func deleteNotificationConfig(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deleted, err := db.DeleteNotificationConfig(r.Context(), userID)
	if err != nil {
		util.ServerError(w, err)
		return
	}

	if !deleted {
		util.NotFound(w)
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func createDefaultConfig(userID string) types.NotificationConfig {
	return types.NotificationConfig{
		ID:              "",
		UserID:          userID,
		DeviceTokens:    []types.DeviceToken{},
		IsEnabled:       true,
		PreferredTime:   "08:00",
		QuietHours:      nil,
		BatchingDays:    1,
		GroupByType:     true,
		MutedPlantIDs:   []string{},
		RemindWatering:  true,
		RemindFertilize: true,
		RemindRepotting: true,
		RemindMisting:   true,
		RemindMoisture:  true,
		UpdatedAt:       time.Now(),
	}
}

func registerDeviceToken(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		Token      string `json:"token"`
		DeviceID   string `json:"deviceId"`
		DeviceType string `json:"deviceType"`
	}

	if err := util.DecodeJSON(r, &request); err != nil {
		util.BadRequest(w, err.Error(), nil)
		return
	}

	if request.Token == "" || request.DeviceID == "" {
		util.BadRequest(w, "token and deviceId are required", nil)
		return
	}

	// Get or create notification config
	config, err := db.GetNotificationConfig(r.Context(), userID)
	if err != nil && err != types.ErrNoDocuments {
		util.ServerError(w, err)
		return
	}

	if config == nil {
		// Create default config
		defaultConfig := createDefaultConfig(userID)
		id, err := gonanoid.New()
		if err != nil {
			log.Printf("Failed to generate notification config ID: %v", err)
			http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
			return
		}
		defaultConfig.ID = id
		config, err = db.CreateNotificationConfig(r.Context(), defaultConfig)
		if err != nil {
			util.ServerError(w, err)
			return
		}
	}

	// Check if device already exists
	now := time.Now()
	deviceExists := false
	for i, device := range config.DeviceTokens {
		if device.DeviceID == request.DeviceID {
			// Update existing device token
			config.DeviceTokens[i].Token = request.Token
			config.DeviceTokens[i].DeviceType = request.DeviceType
			config.DeviceTokens[i].LastUsedAt = now
			config.DeviceTokens[i].IsActive = true
			deviceExists = true
			break
		}
	}

	if !deviceExists {
		// Add new device token
		newDevice := types.DeviceToken{
			Token:      request.Token,
			DeviceID:   request.DeviceID,
			DeviceType: request.DeviceType,
			AddedAt:    now,
			LastUsedAt: now,
			IsActive:   true,
		}
		config.DeviceTokens = append(config.DeviceTokens, newDevice)
	}

	config.UpdatedAt = now

	// Update config in database
	updatedConfig, err := db.UpdateNotificationConfig(r.Context(), *config)
	if err != nil {
		util.ServerError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, updatedConfig)
	log.Printf("Registered device token for user %s, device %s", userID, request.DeviceID)
}

func deleteDeviceToken(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	deviceID := vars["deviceId"]

	if deviceID == "" {
		util.BadRequest(w, "deviceId is required", nil)
		return
	}

	// Get notification config
	config, err := db.GetNotificationConfig(r.Context(), userID)
	if err != nil {
		if err == types.ErrNoDocuments {
			util.NotFound(w)
			return
		}
		util.ServerError(w, err)
		return
	}

	// Remove device token
	found := false
	newTokens := make([]types.DeviceToken, 0, len(config.DeviceTokens))
	for _, device := range config.DeviceTokens {
		if device.DeviceID != deviceID {
			newTokens = append(newTokens, device)
		} else {
			found = true
		}
	}

	if !found {
		util.NotFound(w)
		return
	}

	config.DeviceTokens = newTokens
	config.UpdatedAt = time.Now()

	// Update config in database
	updatedConfig, err := db.UpdateNotificationConfig(r.Context(), *config)
	if err != nil {
		util.ServerError(w, err)
		return
	}

	util.RespondJSON(w, http.StatusOK, updatedConfig)
	log.Printf("Removed device token for user %s, device %s", userID, deviceID)
}
