package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"plants-backend/constants"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type S3Service struct {
	Client    *s3.Client
	Presigner *s3.PresignClient
	Bucket    string
	URLExpire time.Duration
	PublicURL string // optional CDN/base URL for viewing
}

func NewS3Service(ctx context.Context) (*S3Service, error) {
	bucket := getenv("S3_BUCKET", "plants-app-images")
	region := getenv("AWS_REGION", "de")
	endpoint := getenv("S3_ENDPOINT", "https://s3.eu-central-4.ionoscloud.com")

	// Create credentials provider from environment variables
	creds := credentials.NewStaticCredentialsProvider(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	)

	// Load AWS config with explicit credentials and no EC2 IMDS fallback
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	// Create S3 client with custom endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // Ionos S3 uses path-style URLs
	})

	return &S3Service{
		Client:    client,
		Presigner: s3.NewPresignClient(client, s3.WithPresignExpires(1*time.Hour)),
		Bucket:    bucket,
		URLExpire: 1 * time.Hour,
		PublicURL: os.Getenv("S3_PUBLIC_URL"),
	}, nil
}

// SetupCORS configures CORS rules for the bucket to allow browser uploads.
// Call this once during setup or when CORS rules need to be updated.
func (s *S3Service) SetupCORS(ctx context.Context, allowedOrigins []string) error {
	corsRules := []types.CORSRule{
		{
			AllowedOrigins: allowedOrigins,
			AllowedMethods: []string{"GET", "PUT", "POST", "HEAD"},
			AllowedHeaders: []string{
				"*", // Allow all headers including Content-Type, x-amz-*, etc.
			},
			ExposeHeaders: []string{
				"ETag",
				"x-amz-request-id",
			},
			MaxAgeSeconds: aws.Int32(3600),
		},
	}

	_, err := s.Client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
		Bucket: &s.Bucket,
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: corsRules,
		},
	})
	if err != nil {
		return fmt.Errorf("put bucket cors: %w", err)
	}
	return nil
}

// SetupBucketPolicy configures a bucket policy to deny uploads larger than MaxUploadBytes.
// This provides server-side protection against oversized uploads.
func (s *S3Service) SetupBucketPolicy(ctx context.Context) error {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "DenyLargeUploads",
				"Effect":    "Deny",
				"Principal": "*",
				"Action":    "s3:PutObject",
				"Resource":  fmt.Sprintf("arn:aws:s3:::%s/*", s.Bucket),
				"Condition": map[string]interface{}{
					"NumericGreaterThan": map[string]interface{}{
						"s3:content-length": constants.MaxUploadBytes,
					},
				},
			},
		},
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("marshal bucket policy: %w", err)
	}

	policyStr := string(policyJSON)
	_, err = s.Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &s.Bucket,
		Policy: &policyStr,
	})
	if err != nil {
		return fmt.Errorf("put bucket policy: %w", err)
	}
	return nil
}

// GenerateObjectKey builds a unique, user-scoped object key.
func (s *S3Service) GenerateObjectKey(userID, filename string) string {
	id := uuid.New().String()
	return fmt.Sprintf("users/%s/%s_%s", userID, id, sanitizeFilename(filename))
}

// KeyBelongsToUser checks if an object key is scoped under the user's prefix.
func KeyBelongsToUser(key, userID string) bool {
	prefix := fmt.Sprintf("users/%s/", userID)
	return strings.HasPrefix(key, prefix)
}

// PresignPutURL generates a pre-signed PUT URL for direct upload.
// Size limits are enforced via bucket policy (Ionos S3 supports policies).
func (s *S3Service) PresignPutURL(
	ctx context.Context,
	key, contentType string,
	userID string,
) (string, map[string]string, error) {
	params := &s3.PutObjectInput{
		Bucket:      &s.Bucket,
		Key:         &key,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPrivate,
		Metadata:    map[string]string{"user": userID},
	}
	req, err := s.Presigner.PresignPutObject(ctx, params)
	if err != nil {
		return "", nil, fmt.Errorf("presign put: %w", err)
	}
	// Headers client MUST include - must match exactly what's in the signature
	headers := map[string]string{
		"Content-Type":    contentType,
		"x-amz-acl":       "private",
		"x-amz-meta-user": userID,
	}
	return req.URL, headers, nil
}

// PresignGetURL generates a short-lived URL to view a private object.
func (s *S3Service) PresignGetURL(ctx context.Context, key string) (string, error) {
	params := &s3.GetObjectInput{Bucket: &s.Bucket, Key: &key}
	req, err := s.Presigner.PresignGetObject(ctx, params)
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}
	return req.URL, nil
}

// HeadObjectInfo returns object size and content type, verifying existence.
func (s *S3Service) HeadObjectInfo(ctx context.Context, key string) (int64, string, error) {
	out, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &s.Bucket, Key: &key})
	if err != nil {
		return 0, "", fmt.Errorf("head object: %w", err)
	}
	size := aws.ToInt64(out.ContentLength)
	ctype := aws.ToString(out.ContentType)
	return size, ctype, nil
}

// DeleteObject removes an object from S3. Returns error if object doesn't exist or deletion fails.
func (s *S3Service) DeleteObject(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

// DeleteObjects removes multiple objects from S3 in a batch.
func (s *S3Service) DeleteObjects(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	objects := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		k := key
		objects[i] = types.ObjectIdentifier{Key: &k}
	}
	_, err := s.Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: &s.Bucket,
		Delete: &types.Delete{Objects: objects},
	})
	if err != nil {
		return fmt.Errorf("delete objects: %w", err)
	}
	return nil
}

// GetUserUsage returns total bytes and object count for a user's prefix.
func (s *S3Service) GetUserUsage(
	ctx context.Context,
	userID string,
) (total int64, count int, err error) {
	prefix := fmt.Sprintf("users/%s/", userID)
	var token *string
	for {
		out, err := s.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            &s.Bucket,
			Prefix:            &prefix,
			ContinuationToken: token,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("list objects: %w", err)
		}
		for _, obj := range out.Contents {
			total += aws.ToInt64(obj.Size)
			count++
		}
		if aws.ToBool(out.IsTruncated) && out.NextContinuationToken != nil {
			token = out.NextContinuationToken
		} else {
			break
		}
	}
	return total, count, nil
}

func sanitizeFilename(name string) string {
	// Basic sanitization; could expand if needed
	return url.PathEscape(name)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
