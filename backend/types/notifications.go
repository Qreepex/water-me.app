package types

import "time"

type QuietHours struct {
	Start string `json:"start" bson:"start"` // e.g., "22:00"
	End   string `json:"end"   bson:"end"`   // e.g., "07:00"
}

type DeviceToken struct {
	Token      string    `json:"token"      bson:"token"`
	DeviceID   string    `json:"deviceId"   bson:"deviceId"`   // Unique device identifier
	DeviceType string    `json:"deviceType" bson:"deviceType"` // "android", "ios", "web"
	AddedAt    time.Time `json:"addedAt"    bson:"addedAt"`
	LastUsedAt time.Time `json:"lastUsedAt" bson:"lastUsedAt"`
	IsActive   bool      `json:"isActive"   bson:"isActive"` // Can be set to false if token becomes invalid
}

type NotificationConfig struct {
	ID     string `json:"id"     bson:"_id"`
	UserID string `json:"userId" bson:"userId"`

	// FCM Device Tokens for push notifications
	DeviceTokens []DeviceToken `json:"deviceTokens" bson:"deviceTokens"`

	// Globale Einstellungen
	IsEnabled bool `json:"isEnabled" bson:"isEnabled"`

	// Timing & Batching (Zusammenfassung)
	PreferredTime string      `json:"preferredTime"        bson:"preferredTime"` // z.B. "08:30" (HH:mm)
	QuietHours    *QuietHours `json:"quietHours,omitempty" bson:"quietHours,omitempty"`

	// Die "Zusammenfassungs"-Logik
	BatchingDays int  `json:"batchingDays" bson:"batchingDays"` // 0 = Sofort, 1 = Täglich sammeln, 2 = Alle 2 Tage
	GroupByType  bool `json:"groupByType"  bson:"groupByType"`  // Gießen und Düngen in einer Nachricht?

	// Spezifische Filter und Ausnahmen
	MutedPlantIDs   []string `json:"mutedPlantIds"   bson:"mutedPlantIds"` // Diese Pflanzen schicken NIE Benachrichtigungen
	RemindWatering  bool     `json:"remindWatering"  bson:"remindWatering"`
	RemindFertilize bool     `json:"remindFertilize" bson:"remindFertilize"`
	RemindRepotting bool     `json:"remindRepotting" bson:"remindRepotting"`
	RemindMisting   bool     `json:"remindMisting"   bson:"remindMisting"`
	// RemindMoisture asks the user for a new soil moisture measurement so the
	// moisture-based watering plan can be recalculated.
	RemindMoisture bool `json:"remindMoisture" bson:"remindMoisture"`

	// Notification tracking
	LastNotificationSentAt *time.Time `json:"lastNotificationSentAt,omitempty" bson:"lastNotificationSentAt,omitempty"`

	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type NotificationBatch struct {
	UserID       string     `json:"userId"       bson:"userId"`
	DeviceTokens []string   `json:"deviceTokens" bson:"deviceTokens"`
	Plants       []string   `json:"plants"       bson:"plants"` // Plant names
	Type         string     `json:"type"         bson:"type"`   // "watering", "fertilizing", "repotting", "misting"
	ScheduledFor time.Time  `json:"scheduledFor" bson:"scheduledFor"`
	SentAt       *time.Time `json:"sentAt"       bson:"sentAt,omitempty"`
	FailedAt     *time.Time `json:"failedAt"     bson:"failedAt,omitempty"`
	ErrorMessage string     `json:"errorMessage" bson:"errorMessage,omitempty"`
	RetryCount   int        `json:"retryCount"   bson:"retryCount"`
	CreatedAt    time.Time  `json:"createdAt"    bson:"createdAt"`
}
