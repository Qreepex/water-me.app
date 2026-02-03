package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/qreepex/water-me-app/backend/middlewares"
	"github.com/qreepex/water-me-app/backend/routes"
	"github.com/qreepex/water-me-app/backend/services"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	connString := getenv("DATABASE_URL", "mongodb://localhost:27017/plants")
	mongoUser := getenv("MONGODB_USERNAME", "test2")
	mongoPassword := getenv("MONGODB_PASSWORD", "test")
	mongoDatabase := getenv("MONGODB_DATABASE", "plants")

	db, err := services.Connect(connString, mongoDatabase, mongoUser, mongoPassword)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	firebase, err := services.NewFirebaseService()
	if err != nil {
		log.Fatalf("failed to initialize firebase: %v", err)
	}

	// Protected S3 & plant routes
	s3svc, err := services.NewS3Service(ctx)
	if err != nil {
		log.Fatalf("failed to init s3: %v", err)
	}

	// Configure S3 bucket CORS for browser uploads
	// This allows the frontend to make direct PUT requests to S3
	allowedOrigins := []string{
		"https://localhost",
		"http://localhost",
		"https://app.water-me.app", // Allows all origins; restrict in production
		"https://water-me.app",
		"https://my.water-me.app",
	}
	if err := s3svc.SetupCORS(ctx, allowedOrigins); err != nil {
		log.Printf("warning: failed to setup S3 CORS: %v", err)
	} else {
		log.Println("S3 bucket CORS configured successfully")
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "*"}),
		handlers.ExposedHeaders([]string{"Authorization", "Content-Type", "ETag"}),
		handlers.AllowedMethods(
			[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()
	r.Use(cors)

	routes.RegisterRoutes(r, db, s3svc)

	r.Use(middlewares.AuthMiddleware(firebase))

	// Start background cleanup worker for orphaned uploads
	go startCleanupWorker(db, s3svc)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// startCleanupWorker runs a background job to clean up orphaned uploads every 30 minutes
func startCleanupWorker(db *services.MongoDB, s3 *services.S3Service) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	log.Println("Orphaned upload cleanup worker started (runs every 30 minutes)")

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		uploadSvc := services.NewUploadService(db, s3)
		count, err := uploadSvc.CleanupOrphanedUploads(ctx, 1*time.Hour)
		cancel()

		if err != nil {
			log.Printf("Cleanup worker error: %v", err)
		} else if count > 0 {
			log.Printf("Cleaned up %d orphaned uploads", count)
		}
	}
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
