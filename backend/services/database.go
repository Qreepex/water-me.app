package services

import (
	"context"
	"time"

	"github.com/qreepex/water-me-app/backend/constants"
	"github.com/qreepex/water-me-app/backend/types"

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
	unsetDoc := bson.M{}

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
		if isEmptyLocation(*update.Location) {
			unsetDoc["location"] = ""
		} else {
			updateDoc["location"] = *update.Location
		}
	}
	if update.Watering != nil {
		if isEmptyWatering(*update.Watering) {
			unsetDoc["watering"] = ""
		} else {
			updateDoc["watering"] = *update.Watering
		}
	}
	if update.Fertilizing != nil {
		if isEmptyFertilizing(*update.Fertilizing) {
			unsetDoc["fertilizing"] = ""
		} else {
			updateDoc["fertilizing"] = *update.Fertilizing
		}
	}
	if update.Humidity != nil {
		if isEmptyHumidity(*update.Humidity) {
			unsetDoc["humidity"] = ""
		} else {
			updateDoc["humidity"] = *update.Humidity
		}
	}
	if update.Soil != nil {
		if isEmptySoil(*update.Soil) {
			unsetDoc["soil"] = ""
		} else {
			updateDoc["soil"] = *update.Soil
		}
	}
	if update.Seasonality != nil {
		if isEmptySeasonality(*update.Seasonality) {
			unsetDoc["seasonality"] = ""
		} else {
			updateDoc["seasonality"] = *update.Seasonality
		}
	}
	if update.PestHistory != nil {
		if len(*update.PestHistory) == 0 {
			unsetDoc["pestHistory"] = ""
		} else {
			updateDoc["pestHistory"] = *update.PestHistory
		}
	}
	if update.Flags != nil {
		if len(*update.Flags) == 0 {
			unsetDoc["flags"] = ""
		} else {
			updateDoc["flags"] = *update.Flags
		}
	}
	if update.Notes != nil {
		if len(*update.Notes) == 0 {
			unsetDoc["notes"] = ""
		} else {
			updateDoc["notes"] = *update.Notes
		}
	}
	if update.PhotoIDs != nil {
		if len(*update.PhotoIDs) == 0 {
			unsetDoc["photoIds"] = ""
		} else {
			updateDoc["photoIds"] = *update.PhotoIDs
		}
	}
	if update.GrowthHistory != nil {
		if len(*update.GrowthHistory) == 0 {
			unsetDoc["growthHistory"] = ""
		} else {
			updateDoc["growthHistory"] = *update.GrowthHistory
		}
	}

	// Always update the updatedAt timestamp
	updateDoc["updatedAt"] = time.Now()

	updateOps := bson.M{}
	if len(updateDoc) > 0 {
		updateOps["$set"] = updateDoc
	}
	if len(unsetDoc) > 0 {
		updateOps["$unset"] = unsetDoc
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID, "userId": userID},
		updateOps,
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

func isEmptyLocation(loc types.Location) bool {
	return loc.Room == "" && loc.Position == "" && !loc.IsOutdoors
}

func isEmptyWatering(w types.WateringConfig) bool {
	return w.IntervalDays == 0 && w.Method == "" && w.WaterType == "" && w.LastWatered == nil
}

func isEmptyFertilizing(f types.FertilizerConfig) bool {
	return f.Type == "" &&
		f.IntervalDays == 0 &&
		f.NPKRatio == "" &&
		f.ConcentrationPercent == 0 &&
		f.LastFertilized == nil &&
		!f.ActiveInWinter
}

func isEmptyHumidity(h types.HumidityConfig) bool {
	return !h.RequiresMisting &&
		h.MistingIntervalDays == 0 &&
		h.LastMisted == nil &&
		!h.RequiresHumidifier &&
		h.TargetHumidityPct == 0
}

func isEmptySoil(s types.SoilConfig) bool {
	return s.Type == "" && len(s.Components) == 0 && s.LastRepotted == nil && s.RepottingCycle == 0
}

func isEmptySeasonality(s types.SeasonalAdjustments) bool {
	return !s.WinterRestPeriod && s.WinterWaterFactor == 0 && s.MinTempCelsius == 0
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

// GetPlantsNeedingWatering returns plants that need watering (nextWateringDate <= now)
// Processes in batches for scalability. Limit controls batch size (recommended: 1000)
func (m *MongoDB) GetPlantsNeedingWatering(ctx context.Context, limit int) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	now := time.Now()

	// Query plants where watering is configured and lastWatered + intervalDays <= now
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"watering":              bson.M{"$exists": true, "$ne": nil},
			"watering.intervalDays": bson.M{"$gt": 0},
		}},
		bson.M{"$addFields": bson.M{
			"nextWateringDate": bson.M{
				"$dateAdd": bson.M{
					"startDate": "$watering.lastWatered",
					"unit":      "day",
					"amount":    "$watering.intervalDays",
				},
			},
		}},
		bson.M{"$match": bson.M{
			"$or": bson.A{
				bson.M{"watering.lastWatered": nil}, // Never watered
				bson.M{"nextWateringDate": bson.M{"$lte": now}},
			},
		}},
		bson.M{"$limit": limit},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plants []types.Plant
	if err := cursor.All(ctx, &plants); err != nil {
		return nil, err
	}

	return plants, nil
}

// GetPlantsNeedingFertilizer returns plants that need fertilizing
func (m *MongoDB) GetPlantsNeedingFertilizer(
	ctx context.Context,
	limit int,
) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	now := time.Now()

	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"fertilizing":              bson.M{"$exists": true, "$ne": nil},
			"fertilizing.intervalDays": bson.M{"$gt": 0},
		}},
		bson.M{"$addFields": bson.M{
			"nextFertilizingDate": bson.M{
				"$dateAdd": bson.M{
					"startDate": "$fertilizing.lastFertilized",
					"unit":      "day",
					"amount":    "$fertilizing.intervalDays",
				},
			},
		}},
		bson.M{"$match": bson.M{
			"$or": bson.A{
				bson.M{"fertilizing.lastFertilized": nil},
				bson.M{"nextFertilizingDate": bson.M{"$lte": now}},
			},
		}},
		bson.M{"$limit": limit},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plants []types.Plant
	if err := cursor.All(ctx, &plants); err != nil {
		return nil, err
	}

	return plants, nil
}

// GetPlantsNeedingMisting returns plants that need misting
func (m *MongoDB) GetPlantsNeedingMisting(ctx context.Context, limit int) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	now := time.Now()

	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"humidity":                     bson.M{"$exists": true, "$ne": nil},
			"humidity.requiresMisting":     true,
			"humidity.mistingIntervalDays": bson.M{"$gt": 0},
		}},
		bson.M{"$addFields": bson.M{
			"nextMistingDate": bson.M{
				"$dateAdd": bson.M{
					"startDate": "$humidity.lastMisted",
					"unit":      "day",
					"amount":    "$humidity.mistingIntervalDays",
				},
			},
		}},
		bson.M{"$match": bson.M{
			"$or": bson.A{
				bson.M{"humidity.lastMisted": nil},
				bson.M{"nextMistingDate": bson.M{"$lte": now}},
			},
		}},
		bson.M{"$limit": limit},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plants []types.Plant
	if err := cursor.All(ctx, &plants); err != nil {
		return nil, err
	}

	return plants, nil
}

// GetPlantsNeedingRepotting returns plants that need repotting
func (m *MongoDB) GetPlantsNeedingRepotting(ctx context.Context, limit int) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	now := time.Now()

	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"soil":                bson.M{"$exists": true, "$ne": nil},
			"soil.repottingCycle": bson.M{"$gt": 0},
		}},
		bson.M{"$addFields": bson.M{
			"nextRepottingDate": bson.M{
				"$dateAdd": bson.M{
					"startDate": "$soil.lastRepotted",
					"unit":      "day",
					"amount":    "$soil.repottingCycle",
				},
			},
		}},
		bson.M{"$match": bson.M{
			"$or": bson.A{
				bson.M{"soil.lastRepotted": nil},
				bson.M{"nextRepottingDate": bson.M{"$lte": now}},
			},
		}},
		bson.M{"$limit": limit},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plants []types.Plant
	if err := cursor.All(ctx, &plants); err != nil {
		return nil, err
	}

	return plants, nil
}

// GetAllNotificationConfigs returns all notification configs in batches
func (m *MongoDB) GetAllNotificationConfigs(
	ctx context.Context,
	limit int,
	skip int,
) ([]types.NotificationConfig, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []types.NotificationConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// UpdateNotificationLastSent updates the last notification sent timestamp
func (m *MongoDB) UpdateNotificationLastSent(ctx context.Context, userID string) error {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return types.ErrNoDocuments
	}

	now := time.Now()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{"$set": bson.M{"lastNotificationSentAt": now}},
	)

	return err
}

// MarkTokensAsInactive marks the specified FCM tokens as inactive
func (m *MongoDB) MarkTokensAsInactive(ctx context.Context, userID string, tokens []string) error {
	collection := m.GetCollection(constants.MongoDBCollections.Notifications)
	if collection == nil {
		return types.ErrNoDocuments
	}

	// Update all device tokens matching the failed tokens
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{
			"$set": bson.M{
				"deviceTokens.$[elem].isActive": false,
			},
		},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"elem.token": bson.M{"$in": tokens}},
			},
		}),
	)

	return err
}
