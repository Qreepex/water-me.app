package services

import (
	"context"
	"time"

	"plants-backend/constants"
	"plants-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) GetPlants(ctx context.Context, userID string) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	var plants []types.Plant
	cursor, err := collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &plants)
	if err != nil {
		return nil, err
	}

	return plants, nil
}

func (m *MongoDB) CreatePlant(ctx context.Context, plantInput types.Plant) (*types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	result, err := collection.InsertOne(ctx, plantInput)
	if err != nil {
		return nil, err
	}

	// Set the ID from the MongoDB InsertedID
	plantInput.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &plantInput, nil
}

func (m *MongoDB) UpdatePlant(
	ctx context.Context,
	id string,
	userID string,
	update types.UpdatePlantRequest,
) (*types.Plant, bool, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, false, types.ErrNoDocuments
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, false, err
	}

	updateDoc := bson.M{}

	if update.Name != nil {
		updateDoc["name"] = *update.Name
	}
	if update.Species != nil {
		updateDoc["species"] = *update.Species
	}
	if update.IsToxic != nil {
		updateDoc["isToxic"] = *update.IsToxic
	}
	if update.Sunlight != nil {
		updateDoc["sunlight"] = *update.Sunlight
	}
	if update.PreferedTemperature != nil {
		updateDoc["preferedTemperature"] = *update.PreferedTemperature
	}
	if update.Location != nil {
		updateDoc["location"] = *update.Location
	}
	if update.Watering != nil {
		updateDoc["watering"] = *update.Watering
	}
	if update.Fertilizing != nil {
		updateDoc["fertilizing"] = *update.Fertilizing
	}
	if update.Humidity != nil {
		updateDoc["humidity"] = *update.Humidity
	}
	if update.Soil != nil {
		updateDoc["soil"] = *update.Soil
	}
	if update.Seasonality != nil {
		updateDoc["seasonality"] = *update.Seasonality
	}
	if update.PestHistory != nil {
		updateDoc["pestHistory"] = *update.PestHistory
	}
	if update.Flags != nil {
		updateDoc["flags"] = *update.Flags
	}
	if update.Notes != nil {
		updateDoc["notes"] = *update.Notes
	}
	if update.PhotoIDs != nil {
		updateDoc["photoIds"] = *update.PhotoIDs
	}
	if update.GrowthHistory != nil {
		updateDoc["growthHistory"] = *update.GrowthHistory
	}

	// Always update the updatedAt timestamp
	updateDoc["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID, "userId": userID},
		bson.M{"$set": updateDoc},
		opts,
	)

	var plant types.Plant
	if err := result.Decode(&plant); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &plant, true, nil
}

func (m *MongoDB) DeletePlant(ctx context.Context, id string, userID string) (bool, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return false, types.ErrNoDocuments
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID, "userId": userID})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (m *MongoDB) GetPlantBySlug(
	ctx context.Context,
	userID string,
	slug string,
) (*types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	var plant types.Plant
	err := collection.FindOne(ctx, bson.M{"userId": userID, "slug": slug}).Decode(&plant)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	return &plant, nil
}

func (m *MongoDB) GetPlant(
	ctx context.Context,
	userID string,
	id string,
) (*types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var plant types.Plant
	err = collection.FindOne(ctx, bson.M{"userId": userID, "_id": objectID}).Decode(&plant)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	return &plant, nil
}

func (m *MongoDB) WaterPlants(
	ctx context.Context,
	userID string,
	plantIDs []string,
) (int64, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return 0, types.ErrNoDocuments
	}

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(plantIDs))
	for _, id := range plantIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return 0, nil
	}

	now := time.Now()
	result, err := collection.UpdateMany(
		ctx,
		bson.M{
			"_id":    bson.M{"$in": objectIDs},
			"userId": userID,
		},
		bson.M{
			"$set": bson.M{
				"watering.lastWatered": now,
				"updatedAt":            now,
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// --- Uploads ---

// GetUserUploadCount returns number of uploads for a user
func (m *MongoDB) GetUserUploadCount(ctx context.Context, userID string) (int64, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Uploads)
	if collection == nil {
		return 0, types.ErrNoDocuments
	}
	return collection.CountDocuments(ctx, bson.M{"userId": userID})
}

// RegisterUpload records a newly uploaded object key for a user
func (m *MongoDB) RegisterUpload(ctx context.Context, userID string, key string, size int64) error {
	collection := m.GetCollection(constants.MongoDBCollections.Uploads)
	if collection == nil {
		return types.ErrNoDocuments
	}
	doc := types.Upload{
		UserID:    userID,
		Key:       key,
		SizeBytes: size,
		CreatedAt: time.Now(),
	}
	_, err := collection.InsertOne(ctx, doc)
	return err
}

// NotificationConfig methods

func (m *MongoDB) GetNotificationConfig(
	ctx context.Context,
	userID string,
) (*types.NotificationConfig, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	var config types.NotificationConfig
	err := collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&config)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, types.ErrNoDocuments
		}
		return nil, err
	}

	return &config, nil
}

func (m *MongoDB) CreateNotificationConfig(
	ctx context.Context,
	config types.NotificationConfig,
) (*types.NotificationConfig, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	_, err := collection.InsertOne(ctx, config)
	return &config, err
}

func (m *MongoDB) UpdateNotificationConfig(
	ctx context.Context,
	config types.NotificationConfig,
) (*types.NotificationConfig, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"userId": config.UserID},
		bson.M{"$set": config},
		opts,
	)

	var updatedConfig types.NotificationConfig
	if err := result.Decode(&updatedConfig); err != nil {
		return nil, err
	}

	return &updatedConfig, nil
}

func (m *MongoDB) DeleteNotificationConfig(ctx context.Context, userID string) (bool, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return false, types.ErrNoDocuments
	}

	result, err := collection.DeleteOne(ctx, bson.M{"userId": userID})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// CountActiveUsers returns the count of unique users who have plants
func (m *MongoDB) CountActiveUsers(ctx context.Context) (int64, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return 0, types.ErrNoDocuments
	}

	// Count distinct userIds in plants collection
	results, err := collection.Distinct(ctx, "userId", bson.M{})
	if err != nil {
		return 0, err
	}

	return int64(len(results)), nil
}

// CountPlants returns the total count of all plants
func (m *MongoDB) CountPlants(ctx context.Context) (int64, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return 0, types.ErrNoDocuments
	}

	count, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
