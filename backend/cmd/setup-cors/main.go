package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qreepex/water-me-app/backend/services"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()

	s3Svc, err := services.NewS3Service(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 service: %v", err)
	}

	// Add your frontend origins here
	allowedOrigins := []string{
		"http://localhost:5173",                       // Vite dev server
		"http://localhost:4173",                       // Vite preview
		"https://water-me-webapp.qreepex.workers.dev", // Production domain (NO trailing slash!)
	}

	// Add any additional origins from environment
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		allowedOrigins = append(allowedOrigins, extra)
	}

	fmt.Printf("Setting up CORS for bucket: %s\n", s3Svc.Bucket)
	fmt.Printf("Allowed origins: %v\n", allowedOrigins)

	if err := s3Svc.SetupCORS(ctx, allowedOrigins); err != nil {
		log.Fatalf("Failed to set up CORS: %v", err)
	}

	fmt.Println("✓ CORS configuration successfully applied!")
}
