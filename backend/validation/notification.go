package validation

import (
	"regexp"
	"strings"

	"github.com/qreepex/water-me-app/backend/types"
)

type notificationConstraints struct {
	preferredTimeFormat string
	quietHoursFormat    string
	batchingDaysMin     int
	batchingDaysMax     int
	mutedPlantsMaxItems int
}

var defaultNotificationConstraints = notificationConstraints{
	preferredTimeFormat: `^([01]\d|2[0-3]):([0-5]\d)$`, // HH:mm format
	quietHoursFormat:    `^([01]\d|2[0-3]):([0-5]\d)$`, // HH:mm format
	batchingDaysMin:     0,
	batchingDaysMax:     30,
	mutedPlantsMaxItems: 100,
}

// ValidateNotificationConfig validates a NotificationConfig for creation or update.
func ValidateNotificationConfig(config types.NotificationConfig) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	// Validate PreferredTime format (HH:mm)
	preferredTime := strings.TrimSpace(config.PreferredTime)
	if preferredTime == "" {
		errors = append(errors, types.ValidationError{
			Field:   "preferredTime",
			Message: "PreferredTime is required and must be in HH:mm format (e.g., 08:30)",
		})
	} else {
		matched, _ := regexp.MatchString(defaultNotificationConstraints.preferredTimeFormat, preferredTime)
		if !matched {
			errors = append(errors, types.ValidationError{
				Field:   "preferredTime",
				Message: "PreferredTime must be in HH:mm format (e.g., 08:30)",
			})
		}
	}

	// Validate QuietHours if provided
	if config.QuietHours != nil {
		startMatched, _ := regexp.MatchString(
			defaultNotificationConstraints.quietHoursFormat,
			strings.TrimSpace(config.QuietHours.Start),
		)
		if !startMatched {
			errors = append(errors, types.ValidationError{
				Field:   "quietHours.start",
				Message: "QuietHours start time must be in HH:mm format (e.g., 22:00)",
			})
		}

		endMatched, _ := regexp.MatchString(
			defaultNotificationConstraints.quietHoursFormat,
			strings.TrimSpace(config.QuietHours.End),
		)
		if !endMatched {
			errors = append(errors, types.ValidationError{
				Field:   "quietHours.end",
				Message: "QuietHours end time must be in HH:mm format (e.g., 07:00)",
			})
		}
	}

	// Validate BatchingDays
	if config.BatchingDays < defaultNotificationConstraints.batchingDaysMin ||
		config.BatchingDays > defaultNotificationConstraints.batchingDaysMax {
		errors = append(errors, types.ValidationError{
			Field:   "batchingDays",
			Message: "BatchingDays must be between 0 and 30",
		})
	}

	// Validate MutedPlantIDs
	if len(config.MutedPlantIDs) > defaultNotificationConstraints.mutedPlantsMaxItems {
		errors = append(errors, types.ValidationError{
			Field:   "mutedPlantIds",
			Message: "MutedPlantIds array must contain 1000 items or less",
		})
	}

	for _, plantID := range config.MutedPlantIDs {
		if strings.TrimSpace(plantID) == "" {
			errors = append(errors, types.ValidationError{
				Field:   "mutedPlantIds",
				Message: "All plant IDs must be non-empty strings",
			})
			break
		}
	}

	return errors
}

// ValidateUpdateNotificationConfig validates a partial update to NotificationConfig.
func ValidateUpdateNotificationConfig(config types.NotificationConfig) []types.ValidationError {
	// For now, we use the same validation as create since all fields can be updated
	// If you need different rules for updates, implement them here
	return ValidateNotificationConfig(config)
}
