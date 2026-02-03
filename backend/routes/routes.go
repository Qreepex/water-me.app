package routes

import (
	"net/http"

	"github.com/qreepex/water-me-app/backend/constants"
	"github.com/qreepex/water-me-app/backend/services"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, database *services.MongoDB, s3service *services.S3Service) {
	PlantHandler(router, database, s3service)
	UploadHandler(router, database, s3service)
	NotificationHandler(router, database)
	StatsHandler(router, database)
}

func getUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(constants.UserIdKey).(string)
	return userID, ok
}
