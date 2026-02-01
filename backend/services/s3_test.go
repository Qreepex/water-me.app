package services

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// init loads environment variables from .env file before running tests
func init() {
	// Load .env file - it's ok if it doesn't exist in some test environments
	_ = godotenv.Load("../.env")
	_ = godotenv.Load()
}

// TestSetupBucketPolicy verifies that the bucket policy can be set (Ionos S3 supports policies)
func TestSetupBucketPolicy(t *testing.T) {
	ctx := context.Background()

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Skip("AWS credentials not set")
	}

	s3Service, err := NewS3Service(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 service: %v", err)
	}

	// Attempt to set the bucket policy
	err = s3Service.SetupBucketPolicy(ctx)
	if err != nil {
		t.Fatalf("Failed to set bucket policy: %v", err)
	}

	t.Log("Bucket policy set successfully")
}

// TestBucketSizeLimit verifies that the S3 bucket policy rejects uploads larger than MaxUploadBytes
// Ionos S3 supports bucket policies, so oversized uploads should be rejected at S3 level.
func TestBucketSizeLimit(t *testing.T) {
	ctx := context.Background()

	// Check if AWS credentials are configured
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Skip("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables not set")
	}

	// Initialize S3 service
	s3Service, err := NewS3Service(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 service: %v", err)
	}

	// Test user ID for scoped uploads
	userID := "test-user-bucket-limit"

	// Generate presign URL for a 3MB file
	fileSize := int64(3 * 1024 * 1024) // 3MB
	contentType := "image/jpeg"
	filename := "test-large-file.jpg"

	key := s3Service.GenerateObjectKey(userID, filename)
	presignURL, headers, err := s3Service.PresignPutURL(ctx, key, contentType, userID)
	if err != nil {
		t.Fatalf("Failed to generate presign URL: %v", err)
	}

	// Create a 3MB dummy file
	fileData := bytes.Repeat([]byte("a"), int(fileSize))

	// Attempt to upload the 3MB file using the presign URL
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		presignURL,
		bytes.NewReader(fileData),
	)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set required headers from presign response
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.ContentLength = fileSize

	// Execute the upload request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute upload request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	body, _ := io.ReadAll(resp.Body)

	// Verify that the upload was rejected (should be 403 Forbidden or 400 Bad Request)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		t.Errorf(
			"Expected upload to be rejected, but got status %d. Response: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusBadRequest {
		t.Logf(
			"Upload rejected with status %d (expected 403 or 400): %s",
			resp.StatusCode,
			string(body),
		)
	}

	t.Logf("Upload correctly rejected by bucket policy. Status: %d", resp.StatusCode)

	// Cleanup: Try to delete the presigned key from S3 if it somehow got uploaded
	// (This shouldn't succeed due to the bucket policy, but we try anyway)
	_ = s3Service.DeleteObject(ctx, key)
}

// TestValidUploadSize verifies that uploads within the size limit are accepted by S3
func TestValidUploadSize(t *testing.T) {
	ctx := context.Background()

	// Check if AWS credentials are configured
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Skip("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables not set")
	}

	s3Service, err := NewS3Service(ctx)
	if err != nil {
		t.Fatalf("Failed to create S3 service: %v", err)
	}

	userID := "test-user-valid-upload"

	// Create a 1MB file (within 2MB limit)
	fileSize := int64(1 * 1024 * 1024) // 1MB
	contentType := "image/jpeg"
	filename := "test-valid-file.jpg"

	key := s3Service.GenerateObjectKey(userID, filename)
	presignURL, headers, err := s3Service.PresignPutURL(ctx, key, contentType, userID)
	if err != nil {
		t.Fatalf("Failed to generate presign URL: %v", err)
	}

	// Create a 1MB dummy file
	fileData := bytes.Repeat([]byte("b"), int(fileSize))

	// Upload the file
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		presignURL,
		bytes.NewReader(fileData),
	)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.ContentLength = fileSize

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute upload request: %v", err)
	}
	defer resp.Body.Close()

	// Verify upload was accepted (2xx status)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf(
			"Expected upload to succeed, but got status %d. Response: %s",
			resp.StatusCode,
			string(body),
		)
	} else {
		t.Logf("Upload successful. Status: %d", resp.StatusCode)
	}

	// Cleanup
	_ = s3Service.DeleteObject(ctx, key)
}
