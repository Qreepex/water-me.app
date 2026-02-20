package validation

import (
	"strconv"
	"strings"
	"time"

	"github.com/qreepex/water-me-app/backend/types"
)

type validationConstraints struct {
	speciesMinLength          int
	speciesMaxLength          int
	nameMinLength             int
	nameMaxLength             int
	temperatureMin            float64
	temperatureMax            float64
	humidityMin               float64
	humidityMax               float64
	notesMaxItems             int
	notesMaxItemLength        int
	photoIDsMaxItems          int
	photoIDsMaxIDLength       int
	locationRoomMaxLength     int
	locationPositionMaxLength int
	soilTypeMaxLength         int
	soilComponentsMaxItems    int
	soilComponentMaxLength    int
	npkRatioMaxLength         int
	pestTreatmentMaxLength    int
	pestNotesMaxLength        int
	mistingIntervalMin        int
	mistingIntervalMax        int
	repottingCycleMin         int
	repottingCycleMax         int
	concentrationMin          float64
	concentrationMax          float64
	targetHumidityMin         float64
	targetHumidityMax         float64
	wateringIntervalMin       int
	wateringIntervalMax       int
	fertilizingIntervalMin    int
	fertilizingIntervalMax    int
	winterWaterFactorMin      float64
	winterWaterFactorMax      float64
	minTempMin                float64
	minTempMax                float64
	growthHistoryMaxItems     int
	heightCmMin               float64
	heightCmMax               float64
	leafCountMin              int
	leafCountMax              int
	conditionMaxLength        int
}

var constraints = validationConstraints{
	speciesMinLength:          1,
	speciesMaxLength:          100,
	nameMinLength:             1,
	nameMaxLength:             100,
	temperatureMin:            -50,
	temperatureMax:            100,
	humidityMin:               0,
	humidityMax:               100,
	notesMaxItems:             100,
	notesMaxItemLength:        500,
	photoIDsMaxItems:          100,
	photoIDsMaxIDLength:       255,
	locationRoomMaxLength:     100,
	locationPositionMaxLength: 200,
	soilTypeMaxLength:         100,
	soilComponentsMaxItems:    20,
	soilComponentMaxLength:    100,
	npkRatioMaxLength:         20,
	pestTreatmentMaxLength:    200,
	pestNotesMaxLength:        500,
	mistingIntervalMin:        1,
	mistingIntervalMax:        365,
	repottingCycleMin:         1,
	repottingCycleMax:         60,
	concentrationMin:          0,
	concentrationMax:          100,
	targetHumidityMin:         0,
	targetHumidityMax:         100,
	wateringIntervalMin:       1,
	wateringIntervalMax:       365,
	fertilizingIntervalMin:    1,
	fertilizingIntervalMax:    365,
	winterWaterFactorMin:      0.1,
	winterWaterFactorMax:      2.0,
	minTempMin:                -50,
	minTempMax:                50,
	growthHistoryMaxItems:     1000,
	heightCmMin:               0.1,
	heightCmMax:               1000,
	leafCountMin:              0,
	leafCountMax:              10000,
	conditionMaxLength:        200,
}

// ValidateCreatePlantRequest validates a CreatePlantRequest.
func ValidateCreatePlantRequest(req types.CreatePlantRequest) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	// Required fields
	if strings.TrimSpace(req.Name) == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "name",
				Message: "Name is required and must be a non-empty string",
			},
		)
	} else if len(strings.TrimSpace(req.Name)) > constraints.nameMaxLength {
		errors = append(errors, types.ValidationError{Field: "name", Message: "Name must be 100 characters or less"})
	}

	if req.Species != "" && len(strings.TrimSpace(req.Species)) > constraints.speciesMaxLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "species",
				Message: "Species must be 100 characters or less",
			},
		)
	}

	// Sunlight is optional, but if provided must be valid
	if req.Sunlight != nil && !isSunlightRequirement(*req.Sunlight) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "sunlight",
				Message: "Sunlight must be one of: Full Sun, Indirect Sun, Partial Shade, Partial to Full Shade, Full Shade",
			},
		)
	}

	if req.PreferedTemperature != nil {
		if *req.PreferedTemperature < constraints.temperatureMin ||
			*req.PreferedTemperature > constraints.temperatureMax {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "preferedTemperature",
					Message: "PreferredTemperature must be between -50 and 100",
				},
			)
		}
	}

	// Location validation
	if req.Location != nil {
		errors = append(errors, validateLocation(*req.Location)...)
	}

	// Watering validation (optional)
	if req.Watering != nil {
		errors = append(errors, validateWateringConfig(*req.Watering)...)
	}

	// Fertilizing validation
	if req.Fertilizing != nil {
		errors = append(errors, validateFertilizerConfig(*req.Fertilizing)...)
	}

	// Humidity validation
	if req.Humidity != nil {
		errors = append(errors, validateHumidityConfig(*req.Humidity)...)
	}

	// Soil validation
	if req.Soil != nil {
		errors = append(errors, validateSoilConfig(*req.Soil)...)
	}

	// Seasonality validation
	if req.Seasonality != nil {
		errors = append(errors, validateSeasonalAdjustments(*req.Seasonality)...)
	}

	// Pest history validation
	for i, pest := range req.PestHistory {
		errors = append(errors, validatePestInfection(pest, i)...)
	}

	// Flags validation
	for _, flag := range req.Flags {
		if !IsPlantFlag(flag) {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "flags",
					Message: "Flags must be one of: No Draught, Remove Brown Leaves, High Humidity Required, Sensitive Roots",
				},
			)
			break
		}
	}

	// Notes validation
	if len(req.Notes) > constraints.notesMaxItems {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "notes",
				Message: "Notes array must contain 100 items or less",
			},
		)
	}
	for _, note := range req.Notes {
		trimmed := strings.TrimSpace(note)
		if trimmed == "" {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "notes",
					Message: "All notes must be non-empty strings",
				},
			)
			break
		}
		if len(trimmed) > constraints.notesMaxItemLength {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "notes",
					Message: "Each note must be 500 characters or less",
				},
			)
			break
		}
	}

	// PhotoIDs validation
	if len(req.PhotoIDs) > constraints.photoIDsMaxItems {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "photoIds",
				Message: "PhotoIds array must contain 100 items or less",
			},
		)
	}
	for _, id := range req.PhotoIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "photoIds",
					Message: "Each photo ID must be a non-empty string",
				},
			)
			break
		}
		if !strings.HasPrefix(trimmed, "data:") && len(trimmed) > constraints.photoIDsMaxIDLength {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "photoIds",
					Message: "Non-data photo IDs must be 255 characters or less",
				},
			)
			break
		}
	}

	// GrowthHistory validation
	if len(req.GrowthHistory) > constraints.growthHistoryMaxItems {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "growthHistory",
				Message: "GrowthHistory array must contain 1000 items or less",
			},
		)
	}
	for i, log := range req.GrowthHistory {
		errors = append(errors, validateGrowthLog(log, i)...)
	}

	return errors
}

// ValidateUpdatePlantRequest validates an UpdatePlantRequest for PATCH operations.
func ValidateUpdatePlantRequest(req types.UpdatePlantRequest) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			errors = append(
				errors,
				types.ValidationError{Field: "name", Message: "Name must be a non-empty string"},
			)
		} else if len(name) > constraints.nameMaxLength {
			errors = append(errors, types.ValidationError{Field: "name", Message: "Name must be 100 characters or less"})
		}
	}

	if req.Species != nil {
		species := strings.TrimSpace(*req.Species)
		if species == "" {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "species",
					Message: "Species must be a non-empty string",
				},
			)
		} else if len(species) > constraints.speciesMaxLength {
			errors = append(errors, types.ValidationError{Field: "species", Message: "Species must be 100 characters or less"})
		}
	}

	if req.Sunlight != nil && !isSunlightRequirement(*req.Sunlight) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "sunlight",
				Message: "Sunlight must be one of: Full Sun, Indirect Sun, Partial Shade, Partial to Full Shade, Full Shade",
			},
		)
	}

	if req.PreferedTemperature != nil {
		temp := *req.PreferedTemperature
		if temp < constraints.temperatureMin || temp > constraints.temperatureMax {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "preferedTemperature",
					Message: "PreferredTemperature must be between -50 and 100",
				},
			)
		}
	}

	if req.Location != nil {
		errors = append(errors, validateLocation(*req.Location)...)
	}

	if req.Watering != nil {
		errors = append(errors, validateWateringConfig(*req.Watering)...)
	}

	if req.Fertilizing != nil {
		errors = append(errors, validateFertilizerConfig(*req.Fertilizing)...)
	}

	if req.Humidity != nil {
		errors = append(errors, validateHumidityConfig(*req.Humidity)...)
	}

	if req.Soil != nil {
		errors = append(errors, validateSoilConfig(*req.Soil)...)
	}

	if req.Seasonality != nil {
		errors = append(errors, validateSeasonalAdjustments(*req.Seasonality)...)
	}

	if req.PestHistory != nil {
		for i, pest := range *req.PestHistory {
			errors = append(errors, validatePestInfection(pest, i)...)
		}
	}

	if req.Flags != nil {
		for _, flag := range *req.Flags {
			if !IsPlantFlag(flag) {
				errors = append(
					errors,
					types.ValidationError{
						Field:   "flags",
						Message: "Flags must be one of: No Draught, Remove Brown Leaves, High Humidity Required, Sensitive Roots",
					},
				)
				break
			}
		}
	}

	if req.Notes != nil {
		notes := *req.Notes
		if len(notes) > constraints.notesMaxItems {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "notes",
					Message: "Notes array must contain 100 items or less",
				},
			)
		}
		for _, note := range notes {
			trimmed := strings.TrimSpace(note)
			if trimmed == "" {
				errors = append(
					errors,
					types.ValidationError{
						Field:   "notes",
						Message: "All notes must be non-empty strings",
					},
				)
				break
			}
			if len(trimmed) > constraints.notesMaxItemLength {
				errors = append(
					errors,
					types.ValidationError{
						Field:   "notes",
						Message: "Each note must be 500 characters or less",
					},
				)
				break
			}
		}
	}

	if req.PhotoIDs != nil {
		ids := *req.PhotoIDs
		if len(ids) > constraints.photoIDsMaxItems {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "photoIds",
					Message: "PhotoIds array must contain 100 items or less",
				},
			)
		}
		for _, id := range ids {
			trimmed := strings.TrimSpace(id)
			if trimmed == "" {
				errors = append(
					errors,
					types.ValidationError{
						Field:   "photoIds",
						Message: "Each photo ID must be a non-empty string",
					},
				)
				break
			}
			if !strings.HasPrefix(trimmed, "data:") &&
				len(trimmed) > constraints.photoIDsMaxIDLength {
				errors = append(
					errors,
					types.ValidationError{
						Field:   "photoIds",
						Message: "Non-data photo IDs must be 255 characters or less",
					},
				)
				break
			}
		}
	}

	if req.GrowthHistory != nil {
		logs := *req.GrowthHistory
		if len(logs) > constraints.growthHistoryMaxItems {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "growthHistory",
					Message: "GrowthHistory array must contain 1000 items or less",
				},
			)
		}
		for i, log := range logs {
			errors = append(errors, validateGrowthLog(log, i)...)
		}
	}

	return errors
}

func validateLocation(loc types.Location) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	room := strings.TrimSpace(loc.Room)
	if len(room) > 0 && len(room) > constraints.locationRoomMaxLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "location.room",
				Message: "Location room must be 100 characters or less",
			},
		)
	}

	position := strings.TrimSpace(loc.Position)
	if len(position) > 0 && len(position) > constraints.locationPositionMaxLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "location.position",
				Message: "Location position must be 200 characters or less",
			},
		)
	}

	return errors
}

func validateWateringConfig(wc types.WateringConfig) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	if wc.IntervalDays < constraints.wateringIntervalMin ||
		wc.IntervalDays > constraints.wateringIntervalMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "watering.intervalDays",
				Message: "Watering interval must be between 1 and 365 days",
			},
		)
	}

	if !isWateringMethod(wc.Method) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "watering.method",
				Message: "Watering method must be one of: Top, Bottom, Soaking, Self, MistingOnly",
			},
		)
	}

	if !isWaterType(wc.WaterType) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "watering.waterType",
				Message: "Water type must be one of: Tap, Filtered, Rain, Distilled, StaleTap, LowLimestone",
			},
		)
	}

	return errors
}

func validateFertilizerConfig(fc types.FertilizerConfig) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	if !isFertilizerType(fc.Type) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "fertilizing.type",
				Message: "Fertilizer type must be one of: Liquid, Sticks, Granulate, LongTerm, Organic, Hydroponic",
			},
		)
	}

	if fc.IntervalDays < constraints.fertilizingIntervalMin ||
		fc.IntervalDays > constraints.fertilizingIntervalMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "fertilizing.intervalDays",
				Message: "Fertilizing interval must be between 1 and 365 days",
			},
		)
	}

	npk := strings.TrimSpace(fc.NPKRatio)
	if npk == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "fertilizing.npkRatio",
				Message: "NPK ratio must be a non-empty string",
			},
		)
	} else if len(npk) > constraints.npkRatioMaxLength {
		errors = append(errors, types.ValidationError{Field: "fertilizing.npkRatio", Message: "NPK ratio must be 20 characters or less"})
	}

	if fc.ConcentrationPercent < constraints.concentrationMin ||
		fc.ConcentrationPercent > constraints.concentrationMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "fertilizing.concentrationPercent",
				Message: "Concentration must be between 0 and 100",
			},
		)
	}

	return errors
}

func validateHumidityConfig(hc types.HumidityConfig) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	if hc.RequiresMisting {
		if hc.MistingIntervalDays < constraints.mistingIntervalMin ||
			hc.MistingIntervalDays > constraints.mistingIntervalMax {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "humidity.mistingIntervalDays",
					Message: "Misting interval must be between 1 and 365 days",
				},
			)
		}
	}

	if hc.TargetHumidityPct < constraints.targetHumidityMin ||
		hc.TargetHumidityPct > constraints.targetHumidityMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "humidity.targetHumidityPct",
				Message: "Target humidity must be between 0 and 100",
			},
		)
	}

	return errors
}

func validateSoilConfig(sc types.SoilConfig) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	soilType := strings.TrimSpace(sc.Type)
	if soilType == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "soil.type",
				Message: "Soil type must be a non-empty string",
			},
		)
	} else if len(soilType) > constraints.soilTypeMaxLength {
		errors = append(errors, types.ValidationError{Field: "soil.type", Message: "Soil type must be 100 characters or less"})
	}

	if len(sc.Components) > constraints.soilComponentsMaxItems {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "soil.components",
				Message: "Soil components must contain 20 items or less",
			},
		)
	}
	for _, component := range sc.Components {
		trimmed := strings.TrimSpace(component)
		if trimmed == "" {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "soil.components",
					Message: "All soil components must be non-empty strings",
				},
			)
			break
		}
		if len(trimmed) > constraints.soilComponentMaxLength {
			errors = append(
				errors,
				types.ValidationError{
					Field:   "soil.components",
					Message: "Each soil component must be 100 characters or less",
				},
			)
			break
		}
	}

	if sc.RepottingCycle < constraints.repottingCycleMin ||
		sc.RepottingCycle > constraints.repottingCycleMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "soil.repottingCycle",
				Message: "Repotting cycle must be between 1 and 60 months",
			},
		)
	}

	return errors
}

func validateSeasonalAdjustments(sa types.SeasonalAdjustments) []types.ValidationError {
	errors := make([]types.ValidationError, 0)

	if sa.WinterWaterFactor < constraints.winterWaterFactorMin ||
		sa.WinterWaterFactor > constraints.winterWaterFactorMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "seasonality.winterWaterFactor",
				Message: "Winter water factor must be between 0.1 and 2.0",
			},
		)
	}

	if sa.MinTempCelsius < constraints.minTempMin || sa.MinTempCelsius > constraints.minTempMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   "seasonality.minTempCelsius",
				Message: "Minimum temperature must be between -50 and 50",
			},
		)
	}

	return errors
}

func validatePestInfection(pi types.PestInfection, index int) []types.ValidationError {
	errors := make([]types.ValidationError, 0)
	prefix := "pestHistory[" + strconv.Itoa(index) + "]"

	if strings.TrimSpace(pi.ID) == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".id",
				Message: "Pest infection ID must be a non-empty string",
			},
		)
	}

	if !isPestType(pi.Pest) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".pest",
				Message: "Pest type must be one of: Spider Mites, Aphids, Thrips, Mealybugs, Scale, Fungus Gnats, Root Rot",
			},
		)
	}

	if !isPestStatus(pi.Status) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".status",
				Message: "Pest status must be one of: Active, Treated, Resolved",
			},
		)
	}

	treatment := strings.TrimSpace(pi.Treatment)
	if treatment == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".treatment",
				Message: "Treatment must be a non-empty string",
			},
		)
	} else if len(treatment) > constraints.pestTreatmentMaxLength {
		errors = append(errors, types.ValidationError{Field: prefix + ".treatment", Message: "Treatment must be 200 characters or less"})
	}

	notes := strings.TrimSpace(pi.Notes)
	if len(notes) > constraints.pestNotesMaxLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".notes",
				Message: "Notes must be 500 characters or less",
			},
		)
	}

	return errors
}

func validateGrowthLog(gl types.GrowthLog, index int) []types.ValidationError {
	errors := make([]types.ValidationError, 0)
	prefix := "growthHistory[" + strconv.Itoa(index) + "]"

	if strings.TrimSpace(gl.ID) == "" {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".id",
				Message: "Growth log ID must be a non-empty string",
			},
		)
	}

	if gl.HeightCm < constraints.heightCmMin || gl.HeightCm > constraints.heightCmMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".heightCm",
				Message: "Height must be between 0.1 and 1000 cm",
			},
		)
	}

	if gl.LeafCount < constraints.leafCountMin || gl.LeafCount > constraints.leafCountMax {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".leafCount",
				Message: "Leaf count must be between 0 and 10000",
			},
		)
	}

	if !isHealthStatus(gl.Health) {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".health",
				Message: "Health must be one of: Excellent, Good, Fair, Poor, Dormant",
			},
		)
	}

	condition := strings.TrimSpace(gl.Condition)
	if len(condition) > constraints.conditionMaxLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".condition",
				Message: "Condition must be 200 characters or less",
			},
		)
	}

	photoID := strings.TrimSpace(gl.PhotoID)
	if photoID != "" && !strings.HasPrefix(photoID, "data:") &&
		len(photoID) > constraints.photoIDsMaxIDLength {
		errors = append(
			errors,
			types.ValidationError{
				Field:   prefix + ".photoId",
				Message: "Non-data photo ID must be 255 characters or less",
			},
		)
	}

	return errors
}

// Helper functions for enum validation
func IsPlantFlag(flag types.PlantFlag) bool {
	switch flag {
	case types.FlagNoDraught,
		types.FlagRemoveBrownLeaves,
		types.FlagHighHumidity,
		types.FlagSensitiveRoots:
		return true
	default:
		return false
	}
}

func isSunlightRequirement(val types.SunlightRequirement) bool {
	switch val {
	case types.SunlightFullSun,
		types.SunlightIndirectSun,
		types.SunlightPartialShade,
		types.SunlightPartialToFullShade,
		types.SunlightFullShade:
		return true
	default:
		return false
	}
}

func isWateringMethod(method types.WateringMethod) bool {
	switch method {
	case types.MethodTopWatering,
		types.MethodBottomWatering,
		types.MethodSoaking,
		types.MethodSelfWatering,
		types.MethodMistingOnly:
		return true
	default:
		return false
	}
}

func isWaterType(wt types.WaterType) bool {
	switch wt {
	case types.WaterTap,
		types.WaterFiltered,
		types.WaterRain,
		types.WaterDistilled,
		types.WaterStaleTap,
		types.WaterLowLimestone:
		return true
	default:
		return false
	}
}

func isFertilizerType(ft types.FertilizerType) bool {
	switch ft {
	case types.FertilizerLiquid,
		types.FertilizerSticks,
		types.FertilizerGranulate,
		types.FertilizerLongTerm,
		types.FertilizerOrganic,
		types.FertilizerHydroponic:
		return true
	default:
		return false
	}
}

func isPestType(pt types.PestType) bool {
	switch pt {
	case types.PestSpiderMites,
		types.PestAphids,
		types.PestThrips,
		types.PestMealybugs,
		types.PestScale,
		types.PestFungusGnats,
		types.PestRootRot:
		return true
	default:
		return false
	}
}

func isPestStatus(ps types.PestStatus) bool {
	switch ps {
	case types.PestStatusActive, types.PestStatusTreated, types.PestStatusResolved:
		return true
	default:
		return false
	}
}

func isHealthStatus(hs types.HealthStatus) bool {
	switch hs {
	case types.HealthExcellent,
		types.HealthGood,
		types.HealthFair,
		types.HealthPoor,
		types.HealthDormant:
		return true
	default:
		return false
	}
}

func isValidISODate(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return false
	}
	return t.UTC().Format(time.RFC3339Nano) == value
}
