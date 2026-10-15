package constants

var MongoDBCollections = struct {
	Plants        string
	Notifications string
	Uploads       string
}{
	Plants:        "plants",
	Notifications: "notifications",
	Uploads:       "uploads",
}

const UserIdKey = "userID"

// MaxMoistureLogsPerPlant caps the moisture history stored per plant.
const MaxMoistureLogsPerPlant = 2000

// Max allowed upload size (bytes)
const MaxUploadBytes int64 = 2 * 1024 * 1024

// Allowed image MIME types
var AllowedImageContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}
