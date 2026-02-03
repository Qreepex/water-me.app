package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/qreepex/water-me-app/backend/services"

	"github.com/gorilla/mux"
)

// StatsResponse represents the stats response
type StatsResponse struct {
	Users     int64 `json:"users"`
	Plants    int64 `json:"plants"`
	Reminders int64 `json:"reminders"`
}

// StatsHandler registers stats-related routes
func StatsHandler(router *mux.Router, database *services.MongoDB) {
	router.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		getStats(w, r, database)
	}).Methods(http.MethodGet, http.MethodOptions)
}

var (
	cachedStats          *StatsResponse
	cachedStatsTimestamp time.Time
)

func getStats(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	log.Printf("Getting stats")

	if cachedStats != nil && time.Since(cachedStatsTimestamp) < 5*time.Minute {
		log.Printf("Returning cached stats")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedStats); err != nil {
			log.Printf("Failed to encode cached stats: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	ctx := r.Context()

	// Get count of unique users
	users, err := db.CountActiveUsers(ctx)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	// Get count of plants
	plants, err := db.CountPlants(ctx)
	if err != nil {
		log.Printf("Failed to count plants: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	stats := StatsResponse{
		Users:     users,
		Plants:    plants,
		Reminders: 0, // Placeholder for now
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Failed to encode stats: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	cachedStats = &stats
	cachedStatsTimestamp = time.Now()
}
