package services

import (
	"context"
	"strconv"
	"time"

	"github.com/qreepex/water-me-app/backend/constants"
	"github.com/qreepex/water-me-app/backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	if update.MoistureLogs != nil {
		if len(*update.MoistureLogs) == 0 {
			unsetDoc["moistureLogs"] = ""
		} else {
			updateDoc["moistureLogs"] = *update.MoistureLogs
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
	return w.IntervalDays == 0 &&
		w.Method == "" &&
		w.WaterType == "" &&
		w.LastWatered == nil &&
		w.Mode == nil &&
		w.IdealMoisture == nil &&
		w.MoistureScale == nil
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

// AddMoistureLog appends a measurement to a plant and keeps the per-plant list
// capped so the document cannot grow without bound.
func (m *MongoDB) AddMoistureLog(
	ctx context.Context,
	id string,
	userID string,
	log types.MoistureLog,
) (*types.Plant, bool, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, false, types.ErrNoDocuments
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID, "userId": userID},
		bson.M{
			"$push": bson.M{
				"moistureLogs": bson.M{
					"$each":  bson.A{log},
					"$slice": -constants.MaxMoistureLogsPerPlant,
				},
			},
			"$set": bson.M{"updatedAt": now},
		},
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

// DeleteMoistureLog removes a single measurement by id.
func (m *MongoDB) DeleteMoistureLog(ctx context.Context, id string, userID string, logID string) (*types.Plant, bool, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, false, types.ErrNoDocuments
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID, "userId": userID},
		bson.M{
			"$pull": bson.M{"moistureLogs": bson.M{"id": logID}},
			"$set":  bson.M{"updatedAt": now},
		},
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

// WaterPlants marks the given plants as watered and optionally records the
// moisture value(s) measured right before watering.
//
// moistures maps plant ID to the measured value; entries are optional so the
// classic "just water" flow keeps working without a meter.
// WaterPlants records a watering for each plant.
//
// Per your model: watering must NOT fabricate a moisture measurement (a soil
// reading right after watering tells you nothing about the drying rate and
// corrupts the fit). Instead a WateringEvent is appended, carrying the model's
// estimate of the moisture just before and just after watering. If the caller
// supplied a real reading (the user measured while watering) that value is used
// for MoistureAtWatering, and the post-watering value is re-derived from the
// fitted curve rather than guessed.
//
// A real measurement taken shortly after watering is picked up separately (see
// ApplyPostWateringMeasurements) and overwrites the prediction.
func (m *MongoDB) WaterPlants(
	ctx context.Context,
	userID string,
	plantIDs []string,
	moistures map[string]types.MoistureReading,
) (int64, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return 0, types.ErrNoDocuments
	}

	type waterTarget struct {
		id       string
		objectID primitive.ObjectID
	}
	targets := make([]waterTarget, 0, len(plantIDs))
	for _, id := range plantIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		targets = append(targets, waterTarget{id: id, objectID: objectID})
	}
	if len(targets) == 0 {
		return 0, nil
	}

	now := time.Now()
	var modified int64
	for _, target := range targets {
		plant, err := m.GetPlant(ctx, userID, target.id)
		if err != nil {
			return modified, err
		}
		if plant == nil {
			continue
		}

		// Derive the watering's moisture values from the model *before* the
		// watering is recorded, so it describes the state at that moment.
		before, after := predictWateringMoistures(*plant, now)

		// A reading supplied with the request is real data: prefer it.
		if reading, ok := moistures[target.id]; ok {
			before = &reading.Value
		}

		event := types.WateringEvent{
			ID:                    primitive.NewObjectID().Hex(),
			Date:                  now,
			MoistureAtWatering:    before,
			MoistureAfterWatering: after,
			Predicted:             len(moistures) == 0,
		}

		result, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": target.objectID, "userId": userID},
			bson.M{
				"$set":  bson.M{"watering.lastWatered": now, "updatedAt": now},
				"$push": bson.M{"watering.history": event},
			},
		)
		if err != nil {
			return modified, err
		}
		modified += result.ModifiedCount
	}

	return modified, nil
}

// predictWateringMoistures estimates the soil moisture just before and just
// after a watering, using the plant's existing model. Both may be nil when
// there is not enough history to say anything.
func predictWateringMoistures(plant types.Plant, at time.Time) (*float64, *float64) {
	plan := ComputeWateringPlan(plant, at)
	if plan == nil {
		return nil, nil
	}
	before := plan.PredictedMoisture
	// Right before watering, the measured value is the best estimate.
	if plan.CurrentMoisture != nil &&
		plan.WateredSinceLastMeasurement == false &&
		plan.MeasuredAt != nil {
		before = plan.CurrentMoisture
	}
	if plan.EstimatedFullMoisture == nil {
		return before, before
	}
	// Watering refills towards the plant's own saturation level.
	return before, plan.EstimatedFullMoisture
}

// ApplyPostWateringMeasurements refines the most recent watering event when the
// user takes a real reading shortly after watering. The reading replaces the
// predicted post-watering value, which is the best observation of saturation.
func (m *MongoDB) ApplyPostWateringMeasurements(
	ctx context.Context,
	userID string,
	plantID string,
) error {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return types.ErrNoDocuments
	}
	plant, err := m.GetPlant(ctx, userID, plantID)
	if err != nil || plant == nil || plant.Watering == nil {
		return err
	}
	return m.applyPostWateringMeasurement(ctx, collection, userID, *plant)
}

// applyPostWateringMeasurement is the shared implementation.
func (m *MongoDB) applyPostWateringMeasurement(
	ctx context.Context,
	collection *mongo.Collection,
	userID string,
	plant types.Plant,
) error {
	watering := plant.Watering
	if watering == nil || len(watering.History) == 0 {
		return nil
	}
	last := watering.History[len(watering.History)-1]
	window := durationDays(PostWateringMeasurementWindowDays)
	best := -1
	var bestValue float64
	for i, log := range plant.MoistureLogs {
		diff := log.Date.Sub(last.Date)
		if diff < 0 || diff > window {
			continue
		}
		// Use the earliest reading in the window: closest to saturation.
		if best == -1 || log.Date.Before(plant.MoistureLogs[best].Date) {
			best = i
			bestValue = log.Value
		}
	}
	if best == -1 {
		return nil
	}
	index := len(watering.History) - 1
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": plant.ID, "userId": userID},
		bson.M{"$set": bson.M{
			buildHistoryIndexField("moistureAfterWatering", index): bestValue,
			buildHistoryIndexField("predicted", index):             false,
			"updatedAt": time.Now(),
		}},
	)
	return err
}

func buildHistoryIndexField(name string, index int) string {
	return "watering.history." + strconv.Itoa(index) + "." + name
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

// GetMoisturePlants returns plants that use the moisture-based watering plan.
// The actual schedule is derived in Go (see ComputeWateringPlan), so this is
// only a cheap pre-filter.
func (m *MongoDB) GetMoisturePlants(ctx context.Context, limit int) ([]types.Plant, error) {
	collection := m.GetCollection(constants.MongoDBCollections.Plants)
	if collection == nil {
		return nil, types.ErrNoDocuments
	}

	opts := options.Find().SetLimit(int64(limit))
	cursor, err := collection.Find(
		ctx,
		bson.M{"watering.mode": string(types.PlanModeMoisture)},
		opts,
	)
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
