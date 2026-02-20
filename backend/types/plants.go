package types

import (
	"time"
)

// --- ENUMS ---

type SunlightRequirement string

const (
	SunlightFullSun            SunlightRequirement = "Full Sun"
	SunlightIndirectSun        SunlightRequirement = "Indirect Sun"
	SunlightPartialShade       SunlightRequirement = "Partial Shade"
	SunlightPartialToFullShade SunlightRequirement = "Partial to Full Shade"
	SunlightFullShade          SunlightRequirement = "Full Shade"
)

type WateringMethod string

const (
	MethodTopWatering    WateringMethod = "Top"
	MethodBottomWatering WateringMethod = "Bottom"
	MethodSoaking        WateringMethod = "Soaking"
	MethodSelfWatering   WateringMethod = "Self"
	MethodMistingOnly    WateringMethod = "MistingOnly"
)

type WaterType string

const (
	WaterTap          WaterType = "Tap"
	WaterFiltered     WaterType = "Filtered"
	WaterRain         WaterType = "Rain"
	WaterDistilled    WaterType = "Distilled"
	WaterStaleTap     WaterType = "StaleTap"
	WaterLowLimestone WaterType = "LowLimestone"
)

type FertilizerType string

const (
	FertilizerLiquid     FertilizerType = "Liquid"
	FertilizerSticks     FertilizerType = "Sticks"
	FertilizerGranulate  FertilizerType = "Granulate"
	FertilizerLongTerm   FertilizerType = "LongTerm"
	FertilizerOrganic    FertilizerType = "Organic"
	FertilizerHydroponic FertilizerType = "Hydroponic"
)

type PestType string

const (
	PestSpiderMites PestType = "Spider Mites"
	PestAphids      PestType = "Aphids"
	PestThrips      PestType = "Thrips"
	PestMealybugs   PestType = "Mealybugs"
	PestScale       PestType = "Scale"
	PestFungusGnats PestType = "Fungus Gnats"
	PestRootRot     PestType = "Root Rot"
)

type PestStatus string

const (
	PestStatusActive   PestStatus = "Active"
	PestStatusTreated  PestStatus = "Treated"
	PestStatusResolved PestStatus = "Resolved"
)

type PlantFlag string

const (
	FlagNoDraught         PlantFlag = "No Draught"
	FlagRemoveBrownLeaves PlantFlag = "Remove Brown Leaves"
	FlagHighHumidity      PlantFlag = "High Humidity Required"
	FlagSensitiveRoots    PlantFlag = "Sensitive Roots"
)

type HealthStatus string

const (
	HealthExcellent HealthStatus = "Excellent"
	HealthGood      HealthStatus = "Good"
	HealthFair      HealthStatus = "Fair"
	HealthPoor      HealthStatus = "Poor"
	HealthDormant   HealthStatus = "Dormant"
)

// --- SUB-STRUCTURES ---

type Location struct {
	Room       string `json:"room"       bson:"room"`
	Position   string `json:"position"   bson:"position"`
	IsOutdoors bool   `json:"isOutdoors" bson:"isOutdoors"`
}

type WateringConfig struct {
	IntervalDays int            `json:"intervalDays" bson:"intervalDays"`
	Method       WateringMethod `json:"method"       bson:"method"`
	WaterType    WaterType      `json:"waterType"    bson:"waterType"`
	LastWatered  *time.Time     `json:"lastWatered"  bson:"lastWatered"`
}

type HumidityConfig struct {
	RequiresMisting     bool       `json:"requiresMisting"     bson:"requiresMisting"`
	MistingIntervalDays int        `json:"mistingIntervalDays" bson:"mistingIntervalDays"`
	LastMisted          *time.Time `json:"lastMisted"          bson:"lastMisted"`
	RequiresHumidifier  bool       `json:"requiresHumidifier"  bson:"requiresHumidifier"`
	TargetHumidityPct   float64    `json:"targetHumidityPct"   bson:"targetHumidityPct"`
}

type FertilizerConfig struct {
	Type                 FertilizerType `json:"type"                 bson:"type"`
	IntervalDays         int            `json:"intervalDays"         bson:"intervalDays"`
	NPKRatio             string         `json:"npkRatio"             bson:"npkRatio"`
	ConcentrationPercent float64        `json:"concentrationPercent" bson:"concentrationPercent"`
	LastFertilized       *time.Time     `json:"lastFertilized"       bson:"lastFertilized"`
	ActiveInWinter       bool           `json:"activeInWinter"       bson:"activeInWinter"`
}

type SoilConfig struct {
	Type           string     `json:"type"           bson:"type"`
	Components     []string   `json:"components"     bson:"components"`
	LastRepotted   *time.Time `json:"lastRepotted"   bson:"lastRepotted"`
	RepottingCycle int        `json:"repottingCycle" bson:"repottingCycle"`
}

type PestInfection struct {
	ID         string     `json:"id"                   bson:"id"`
	Pest       PestType   `json:"pest"                 bson:"pest"`
	DetectedAt time.Time  `json:"detectedAt"           bson:"detectedAt"`
	ResolvedAt *time.Time `json:"resolvedAt,omitempty" bson:"resolvedAt,omitempty"`
	Status     PestStatus `json:"status"               bson:"status"`
	Treatment  string     `json:"treatment"            bson:"treatment"`
	Notes      string     `json:"notes"                bson:"notes"`
}

type SeasonalAdjustments struct {
	WinterRestPeriod  bool    `json:"winterRestPeriod"  bson:"winterRestPeriod"`
	WinterWaterFactor float64 `json:"winterWaterFactor" bson:"winterWaterFactor"`
	MinTempCelsius    float64 `json:"minTempCelsius"    bson:"minTempCelsius"`
}

type GrowthLog struct {
	ID        string       `json:"id"        bson:"id"`
	Date      time.Time    `json:"date"      bson:"date"`
	HeightCm  float64      `json:"heightCm"  bson:"heightCm"`
	LeafCount int          `json:"leafCount" bson:"leafCount"`
	Health    HealthStatus `json:"health"    bson:"health"`
	Condition string       `json:"condition" bson:"condition"`
	PhotoID   string       `json:"photoId"   bson:"photoId"`
}

// --- MAIN STRUCT ---

type Plant struct {
	ID      string `json:"id"      bson:"_id,omitempty"`
	UserID  string `json:"userId"  bson:"userId"`
	Slug    string `json:"slug"    bson:"slug"`
	Name    string `json:"name"    bson:"name"`
	Species string `json:"species" bson:"species"`
	IsToxic bool   `json:"isToxic" bson:"isToxic"`

	Sunlight            *SunlightRequirement `json:"sunlight"  bson:"sunlight,omitempty"`
	PreferedTemperature *float64             `json:"preferedTemperature" bson:"preferedTemperature,omitempty"`
	Location            *Location            `json:"location"  bson:"location,omitempty"`

	Watering    *WateringConfig      `json:"watering"    bson:"watering,omitempty"`
	Fertilizing *FertilizerConfig    `json:"fertilizing" bson:"fertilizing,omitempty"`
	Humidity    *HumidityConfig      `json:"humidity"    bson:"humidity,omitempty"`
	Soil        *SoilConfig          `json:"soil"        bson:"soil,omitempty"`
	Seasonality *SeasonalAdjustments `json:"seasonality" bson:"seasonality,omitempty"`

	PestHistory   []PestInfection `json:"pestHistory"         bson:"pestHistory,omitempty"`
	Flags         []PlantFlag     `json:"flags"               bson:"flags,omitempty"`
	Notes         []string        `json:"notes"               bson:"notes,omitempty"`
	PhotoIDs      []string        `json:"photoIds"            bson:"photoIds,omitempty"`
	PhotoURLs     []string        `json:"photoUrls" bson:"-"`
	GrowthHistory []GrowthLog     `json:"growthHistory"       bson:"growthHistory,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// CreatePlantRequest is the request body for creating a new plant.
type CreatePlantRequest struct {
	Name                string               `json:"name"`
	Species             string               `json:"species,omitempty"`
	IsToxic             bool                 `json:"isToxic"`
	Sunlight            *SunlightRequirement `json:"sunlight,omitempty"`
	PreferedTemperature *float64             `json:"preferedTemperature,omitempty"`
	Location            *Location            `json:"location,omitempty"`
	Watering            *WateringConfig      `json:"watering,omitempty"`
	Fertilizing         *FertilizerConfig    `json:"fertilizing,omitempty"`
	Humidity            *HumidityConfig      `json:"humidity,omitempty"`
	Soil                *SoilConfig          `json:"soil,omitempty"`
	Seasonality         *SeasonalAdjustments `json:"seasonality,omitempty"`
	PestHistory         []PestInfection      `json:"pestHistory"`
	Flags               []PlantFlag          `json:"flags"`
	Notes               []string             `json:"notes"`
	PhotoIDs            []string             `json:"photoIds"`
	GrowthHistory       []GrowthLog          `json:"growthHistory"`
}

// UpdatePlantRequest is for PATCH operations with optional fields.
type UpdatePlantRequest struct {
	Name                *string              `json:"name,omitempty"`
	Species             *string              `json:"species,omitempty"`
	IsToxic             *bool                `json:"isToxic,omitempty"`
	Sunlight            *SunlightRequirement `json:"sunlight,omitempty"`
	PreferedTemperature *float64             `json:"preferedTemperature,omitempty"`
	Location            *Location            `json:"location,omitempty"`
	Watering            *WateringConfig      `json:"watering,omitempty"`
	Fertilizing         *FertilizerConfig    `json:"fertilizing,omitempty"`
	Humidity            *HumidityConfig      `json:"humidity,omitempty"`
	Soil                *SoilConfig          `json:"soil,omitempty"`
	Seasonality         *SeasonalAdjustments `json:"seasonality,omitempty"`
	PestHistory         *[]PestInfection     `json:"pestHistory,omitempty"`
	Flags               *[]PlantFlag         `json:"flags,omitempty"`
	Notes               *[]string            `json:"notes,omitempty"`
	PhotoIDs            *[]string            `json:"photoIds,omitempty"`
	GrowthHistory       *[]GrowthLog         `json:"growthHistory,omitempty"`
}
