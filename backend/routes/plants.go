package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"plants-backend/services"
	"plants-backend/types"
	"plants-backend/util"
	"plants-backend/validation"

	"github.com/gorilla/mux"
)

// getRealIP extracts the real client IP, considering Cloudflare headers
func getRealIP(r *http.Request) string {
	// Cloudflare header
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	// Fallback to X-Forwarded-For
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// Fallback to direct connection
	return r.RemoteAddr
}

// PlantHandler registers plant-related routes
func PlantHandler(router *mux.Router, database *services.MongoDB, s3 *services.S3Service) {
	// Create rate limiter once
	rateLimiter := services.NewRateLimiter()

	// Apply rate limiting to all routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRealIP(r)
			userID, ok := getUserID(r)
			if ok && rateLimiter.IsRateLimited(userID, ip) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/api/plants", func(w http.ResponseWriter, r *http.Request) {
		getPlants(w, r, database, s3)
	}).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/plants", func(w http.ResponseWriter, r *http.Request) {
		createPlant(w, r, database)
	}).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/plants/water", func(w http.ResponseWriter, r *http.Request) {
		waterPlants(w, r, database)
	}).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/plants/slug/{slug}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]
		getPlantBySlug(w, r, database, s3, slug)
	}).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/plants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		updatePlant(w, r, database, s3, id)
	}).Methods(http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/api/plants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		deletePlant(w, r, database, id)
	}).Methods(http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/api/plants/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		getPlant(w, r, database, s3, id)
	}).Methods(http.MethodGet, http.MethodOptions)
}

func getPlants(
	w http.ResponseWriter,
	r *http.Request,
	db *services.MongoDB,
	s3 *services.S3Service,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Getting plants for %v", userID)

	plants, err := db.GetPlants(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to retrieve plants: %v", err)
		http.Error(w, "Failed to retrieve plants", http.StatusInternalServerError)
		return
	}
	if plants == nil {
		plants = []types.Plant{}
	}
	// Enrich with signed photo URLs
	for i := range plants {
		plants[i].PhotoURLs = resolvePhotoURLs(r.Context(), s3, plants[i].PhotoIDs, userID)
	}
	util.RespondJSON(w, http.StatusOK, plants)
}

func getPlantBySlug(
	w http.ResponseWriter,
	r *http.Request,
	db *services.MongoDB,
	s3 *services.S3Service,
	slug string,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Getting plant by slug '%s' for user %v", slug, userID)
	plant, err := db.GetPlantBySlug(r.Context(), userID, slug)
	if err != nil {
		log.Printf("Failed to retrieve plant by slug: %v", err)
		http.Error(w, "Failed to retrieve plant", http.StatusInternalServerError)
		return
	}
	if plant == nil {
		util.NotFound(w)
		return
	}
	plant.PhotoURLs = resolvePhotoURLs(r.Context(), s3, plant.PhotoIDs, userID)
	util.RespondJSON(w, http.StatusOK, plant)
}

func getPlant(
	w http.ResponseWriter,
	r *http.Request,
	db *services.MongoDB,
	s3 *services.S3Service,
	id string,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Getting plant by ID '%s' for user %v", id, userID)
	plant, err := db.GetPlant(r.Context(), userID, id)
	if err != nil {
		log.Printf("Failed to retrieve plant by ID: %v", err)
		http.Error(w, "Failed to retrieve plant", http.StatusInternalServerError)
		return
	}
	if plant == nil {
		util.NotFound(w)
		return
	}
	plant.PhotoURLs = resolvePhotoURLs(r.Context(), s3, plant.PhotoIDs, userID)
	util.RespondJSON(w, http.StatusOK, plant)
}

func createPlant(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req types.CreatePlantRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.BadRequest(w, err.Error(), nil)
		return
	}
	errors := validation.ValidateCreatePlantRequest(req)
	if len(errors) > 0 {
		util.BadRequest(w, "Validation failed", errors)
		return
	}
	existingPlants, err := db.GetPlants(r.Context(), userID)
	if err != nil {
		util.ServerError(w, err)
		return
	}

	// Check plant limit
	if validation.ValidatePlantLimit(len(existingPlants)) {
		util.BadRequest(w, "Plant limit exceeded", map[string]interface{}{
			"limit":   validation.MaxPlantsPerUser,
			"current": len(existingPlants),
		})
		return
	}

	plant := createPlantFromRequest(req, userID, existingPlants)
	createdPlant, err := db.CreatePlant(r.Context(), plant)
	if err != nil {
		util.ServerError(w, err)
		return
	}
	util.RespondJSON(w, http.StatusCreated, createdPlant)
}

func updatePlant(
	w http.ResponseWriter,
	r *http.Request,
	db *services.MongoDB,
	s3 *services.S3Service,
	id string,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req types.UpdatePlantRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.BadRequest(w, err.Error(), nil)
		return
	}
	errors := validation.ValidateUpdatePlantRequest(req)
	if len(errors) > 0 {
		util.BadRequest(w, "Validation failed", errors)
		return
	}
	plant, found, err := db.UpdatePlant(r.Context(), id, userID, req)
	if err != nil {
		util.ServerError(w, err)
		return
	}
	if !found {
		util.NotFound(w)
		return
	}
	plant.PhotoURLs = resolvePhotoURLs(r.Context(), s3, plant.PhotoIDs, userID)
	util.RespondJSON(w, http.StatusOK, plant)
}

func deletePlant(w http.ResponseWriter, r *http.Request, db *services.MongoDB, id string) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	deleted, err := db.DeletePlant(r.Context(), id, userID)
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

func waterPlants(w http.ResponseWriter, r *http.Request, db *services.MongoDB) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		PlantIDs []string `json:"plantIds"`
	}
	if err := util.DecodeJSON(r, &req); err != nil {
		util.BadRequest(w, err.Error(), nil)
		return
	}
	if len(req.PlantIDs) == 0 {
		util.BadRequest(w, "At least one plant ID is required", nil)
		return
	}
	_, err := db.WaterPlants(r.Context(), userID, req.PlantIDs)
	if err != nil {
		util.ServerError(w, err)
		return
	}
	util.RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// generateUniqueSlug creates a unique slug for the plant within the user's collection
func generateUniqueSlug(
	name string,
	location *types.Location,
	existingPlants []types.Plant,
) string {
	baseSlug := slugify(name)
	if baseSlug == "" {
		baseSlug = "plant"
	}
	if !slugExists(baseSlug, existingPlants) {
		return baseSlug
	}
	var locationPart string
	if location != nil {
		locationPart = slugify(location.Room)
		if locationPart == "" {
			locationPart = slugify(location.Position)
		}
	}
	if locationPart != "" {
		slugWithLocation := baseSlug + "-" + locationPart
		if !slugExists(slugWithLocation, existingPlants) {
			return slugWithLocation
		}
	}
	counter := 1
	for {
		numberedSlug := fmt.Sprintf("%s-%d", baseSlug, counter)
		if !slugExists(numberedSlug, existingPlants) {
			return numberedSlug
		}
		counter++
	}
}

// slugExists checks if a slug already exists in the user's plants
func slugExists(slug string, plants []types.Plant) bool {
	for _, plant := range plants {
		if plant.Slug == slug {
			return true
		}
	}
	return false
}

func createPlantFromRequest(
	req types.CreatePlantRequest,
	userID string,
	existingPlants []types.Plant,
) types.Plant {
	now := time.Now()
	slug := generateUniqueSlug(req.Name, req.Location, existingPlants)
	plant := types.Plant{
		UserID:              userID,
		Slug:                slug,
		Name:                req.Name,
		Species:             req.Species,
		IsToxic:             req.IsToxic,
		Sunlight:            req.Sunlight,
		PreferedTemperature: req.PreferedTemperature,
		Location:            req.Location,
		Watering:            req.Watering,
		Fertilizing:         req.Fertilizing,
		Humidity:            req.Humidity,
		Soil:                req.Soil,
		Seasonality:         req.Seasonality,
		PestHistory:         req.PestHistory,
		Flags:               req.Flags,
		Notes:               req.Notes,
		PhotoIDs:            req.PhotoIDs,
		GrowthHistory:       req.GrowthHistory,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	return plant
}

// resolvePhotoURLs presigns GET URLs for user's images
func resolvePhotoURLs(
	ctx context.Context,
	s3 *services.S3Service,
	keys []string,
	userID string,
) []string {
	urls := make([]string, 0, len(keys))
	for _, k := range keys {
		if !services.KeyBelongsToUser(k, userID) {
			continue
		}
		url, err := s3.PresignGetURL(ctx, k)
		if err == nil && url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}
