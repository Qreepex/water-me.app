package routes

import (
	"net/http"

	"plants-backend/constants"
	"plants-backend/services"
	"plants-backend/util"

	"github.com/gorilla/mux"
)

// UploadHandler registers upload-related routes
func UploadHandler(router *mux.Router, database *services.MongoDB, s3 *services.S3Service) {
	router.HandleFunc("/api/uploads/presign", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Filename    string `json:"filename"`
			ContentType string `json:"contentType"`
			SizeBytes   int64  `json:"sizeBytes"`
		}
		if err := util.DecodeJSON(r, &req); err != nil {
			util.BadRequest(w, err.Error(), nil)
			return
		}
		if req.Filename == "" || req.ContentType == "" || req.SizeBytes <= 0 {
			util.BadRequest(w, "filename, contentType and sizeBytes are required", nil)
			return
		}

		if !constants.AllowedImageContentTypes[req.ContentType] {
			util.BadRequest(
				w,
				"Unsupported content type",
				map[string]interface{}{"allowed": constants.AllowedImageContentTypes},
			)
			return
		}

		if req.SizeBytes > constants.MaxUploadBytes {
			util.BadRequest(
				w,
				"File too large (max 2MB)",
				map[string]interface{}{"maxBytes": constants.MaxUploadBytes},
			)
			return
		}

		count, err := database.GetUserUploadCount(r.Context(), userID)
		if err != nil {
			util.ServerError(w, err)
			return
		}
		if count >= 20 {
			util.BadRequest(w, "Upload limit reached (20 per user)", nil)
			return
		}

		key := s3.GenerateObjectKey(userID, req.Filename)
		url, headers, err := s3.PresignPutURL(r.Context(), key, req.ContentType, userID)
		if err != nil {
			util.ServerError(w, err)
			return
		}
		util.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"key":     key,
			"url":     url,
			"headers": headers,
		})
	}).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/uploads/register", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Key string `json:"key"`
		}
		if err := util.DecodeJSON(r, &req); err != nil {
			util.BadRequest(w, err.Error(), nil)
			return
		}
		if req.Key == "" || !services.KeyBelongsToUser(req.Key, userID) {
			util.BadRequest(w, "Invalid or unauthorized key", nil)
			return
		}

		size, contentType, err := s3.HeadObjectInfo(r.Context(), req.Key)
		if err != nil {
			util.BadRequest(w, "Uploaded object not found", nil)
			return
		}

		if size > constants.MaxUploadBytes {
			util.BadRequest(
				w,
				"Uploaded file exceeds 2MB limit",
				map[string]interface{}{"maxBytes": constants.MaxUploadBytes},
			)
			return
		}

		if !constants.AllowedImageContentTypes[contentType] {
			util.BadRequest(
				w,
				"Uploaded file type not allowed",
				map[string]interface{}{
					"allowed":     constants.AllowedImageContentTypes,
					"contentType": contentType,
				},
			)
			return
		}

		if err := database.RegisterUpload(r.Context(), userID, req.Key, size); err != nil {
			util.ServerError(w, err)
			return
		}
		util.RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
	}).Methods(http.MethodPost, http.MethodOptions)

	// DELETE /api/uploads/{key} - Delete an upload and remove from S3
	router.HandleFunc("/api/uploads/{key}", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := getUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		key := mux.Vars(r)["key"]
		if key == "" || !services.KeyBelongsToUser(key, userID) {
			util.BadRequest(w, "Invalid or unauthorized key", nil)
			return
		}

		uploadSvc := services.NewUploadService(database, s3)
		if err := uploadSvc.DeleteUpload(r.Context(), key, userID); err != nil {
			util.ServerError(w, err)
			return
		}

		util.RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
	}).Methods(http.MethodDelete, http.MethodOptions)
}
