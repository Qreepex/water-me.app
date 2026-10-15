package services

import (
	"sort"
	"time"

	"github.com/qreepex/water-me-app/backend/types"
)

// Moisture-based watering scheduler.
//
// The idea: a soil moisture meter gives a value on a known scale. The plant
// dries out roughly linearly between two waterings. If we know the drying
// speed ("decay", moisture units lost per day) and the moisture the plant is
// happiest at ("ideal"), we can predict the moment the soil reaches the ideal
// value again and remind the user to water then.
//
// To keep the model fed we ask for a fresh measurement in a sensible cadence:
// roughly half way until the predicted watering, but never more often than once
// per day and at most once a week.
// ComputeWateringPlan derives the watering schedule for a plant.
//
// For manual plants it mirrors the interval logic the notification worker and
// the app already use (lastWatered + intervalDays). For moisture plants it
// predicts, from the recorded measurements, when the soil will reach the ideal
// moisture again and when the next measurement should be taken. The returned
// plan is never nil as long as the plant has a usable configuration.
func ComputeWateringPlan(plant types.Plant, now time.Time) *types.WateringPlan {
	mode := types.PlanModeManual
	if plant.Watering != nil && plant.Watering.Mode != nil {
		mode = *plant.Watering.Mode
	}

	if mode == types.PlanModeMoisture && plant.Watering != nil {
		return computeMoisturePlan(plant, now)
	}
	return computeManualPlan(plant, now)
}

func computeManualPlan(plant types.Plant, now time.Time) *types.WateringPlan {
	if plant.Watering == nil || plant.Watering.IntervalDays <= 0 {
		return nil
	}

	plan := &types.WateringPlan{Mode: types.PlanModeManual}
	last := plant.Watering.LastWatered
	if last == nil {
		plan.WateringDueAt = &now
		plan.Status = types.PlanStatusWaterNow
		return plan
	}

	due := last.AddDate(0, 0, plant.Watering.IntervalDays)
	plan.WateringDueAt = &due
	plan.Status = wateringStatus(due, plant.Watering.IntervalDays, now)
	return plan
}

// wateringStatus classifies a due date into the app-facing statuses.
func wateringStatus(due time.Time, intervalDays int, now time.Time) string {
	if !now.Before(due) {
		return types.PlanStatusWaterNow
	}
	// "Due soon" = within the last quarter of the interval (at least 2 days).
	soonWindow := float64(intervalDays) / 4.0
	if soonWindow < 2 {
		soonWindow = 2
	}
	if due.Sub(now) <= time.Duration(soonWindow*24)*time.Hour {
		return types.PlanStatusDueSoon
	}
	return types.PlanStatusOK
}

func computeMoisturePlan(plant types.Plant, now time.Time) *types.WateringPlan {
	watering := plant.Watering
	scale := types.ScaleTen
	if watering.MoistureScale != nil {
		scale = *watering.MoistureScale
	}
	plan := &types.WateringPlan{
		Mode:          types.PlanModeMoisture,
		IdealMoisture: watering.IdealMoisture,
		MoistureScale: scale,
	}

	if watering.IdealMoisture == nil {
		// Misconfigured: without a target there is nothing to compute.
		plan.NeedsFirstMeasurement = true
		plan.Status = types.PlanStatusMeasure
		plan.MeasurementDueAt = &now
		return plan
	}

	ideal := normalizeMoisture(*watering.IdealMoisture, scale)

	// The configured value is the centre of a target band, not a floor: the soil
	// is watered once it falls `band` below the target and refilled to about
	// `band` above it.
	bandRaw := DefaultMoistureBand
	if watering.MoistureBand != nil && *watering.MoistureBand > 0 {
		bandRaw = *watering.MoistureBand
	}
	band := normalizeBand(bandRaw, scale)
	if band > 1 {
		band = 1
	}
	waterBelow := clamp01(ideal - band)
	waterUntil := clamp01(ideal + band)
	plan.WaterBelow = floatPtr(denormalizeMoisture(waterBelow, scale))
	plan.WaterUntil = floatPtr(denormalizeMoisture(waterUntil, scale))

	// ---- readings ----
	samples := make([]moistureSample, 0, len(plant.MoistureLogs))
	for _, log := range plant.MoistureLogs {
		logScale := log.Scale
		if logScale == "" {
			logScale = scale
		}
		samples = append(samples, moistureSample{
			at:    log.Date,
			value: normalizeMoisture(log.Value, logScale),
			raw:   log.Value,
		})
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i].at.Before(samples[j].at) })
	plan.SampleCount = len(samples)

	if len(samples) == 0 {
		// Nothing recorded yet: ask for the first value now.
		plan.NeedsFirstMeasurement = true
		plan.MeasurementDueAt = &now
		plan.Status = types.PlanStatusMeasure
		return plan
	}

	latest := samples[len(samples)-1]
	plan.CurrentMoisture = &latest.raw
	plan.MeasuredAt = &latest.at

	// ---- fit the per-plant model over every recorded cycle ----
	model := fitModel(samples, effectiveWaterings(watering), now)
	plan.CyclesUsed = model.cyclesUsed
	plan.Confidence = model.confidence
	plan.ModelReady = model.ready

	decay := denormalizeDelta(model.decayPerDay, scale)
	if model.decayPerDay > 1e-9 {
		plan.DecayPerDay = &decay
	}

	// The model's own estimate of the current moisture, used whenever the user
	// has not measured recently.
	if model.cyclesUsed > 0 {
		if model.ready {
			predicted := denormalizeMoisture(model.moistureAt(now), scale)
			plan.PredictedMoisture = &predicted
		}
		full := denormalizeMoisture(model.saturationAt(), scale)
		plan.EstimatedFullMoisture = &full
	} else {
		full := denormalizeMoisture(maxSampleValue(samples), scale)
		plan.EstimatedFullMoisture = &full
	}

	// Is the newest reading already stale (watered afterwards)?
	plan.WateredSinceLastMeasurement = watering.LastWatered != nil &&
		watering.LastWatered.After(latest.at)

	// ---- when is the next watering due? ----
	if model.cyclesUsed > 0 {
		// A watering *resets* the soil to saturation, so drying restarts there.
		// Without a recorded watering we simply continue the measured curve from
		// the last reading.
		var days float64
		var ok bool
		var anchor time.Time
		if watering.LastWatered != nil && watering.LastWatered.After(latest.at) {
			anchor = *watering.LastWatered
			days, ok = model.daysUntilFromSaturation(anchor, waterBelow)
		} else {
			anchor = latest.at
			days, ok = model.daysUntil(anchor, waterBelow)
		}
		if ok {
			due := anchor.Add(durationDays(days))
			plan.WateringDueAt = &due
		}
	}

	// The model is only "ready" once enough cycles agree; until then the plan
	// keeps asking for readings at a fairly high cadence.
	plan.NeedsFirstMeasurement = !model.ready

	// ---- when should the user measure again? ----
	plan.MeasurementDueAt, plan.DriftCheckDueAt = measurementSchedule(samples, model, waterBelow, now)

	switch {
	case plan.WateringDueAt != nil && !now.Before(*plan.WateringDueAt):
		plan.Status = types.PlanStatusWaterNow
	case plan.MeasurementDueAt != nil && !now.Before(*plan.MeasurementDueAt):
		plan.Status = types.PlanStatusMeasure
	case plan.WateringDueAt != nil &&
		plan.WateringDueAt.Sub(now) <= durationDays(dueSoonDays(*plan.WateringDueAt, now)):
		plan.Status = types.PlanStatusDueSoon
	default:
		plan.Status = types.PlanStatusOK
	}

	return plan
}

// measurementSchedule decides when the next reading is wanted.
//
// While calibrating the app asks fairly often (the data is the bottleneck).
// Once the model is confident the user should not have to measure at all, so
// only a rare drift check remains, to notice if reality moved away from the
// prediction (season change, repotting, a forgotten watering).
func measurementSchedule(
	samples []moistureSample,
	model moistureModel,
	waterBelow float64,
	now time.Time,
) (next *time.Time, drift *time.Time) {
	latest := samples[len(samples)-1]
	calendarYear := 365.0

	if !model.ready {
		// --- calibration phase ---
		days := defaultMeasurementIntervalDays
		if model.cyclesUsed > 0 {
			// We know the rate: aim for half way to the next watering so the
			// point lands on the drying curve.
			if toWater, ok := model.daysUntil(latest.at, waterBelow); ok {
				days = toWater / 2.0
			}
		} else if len(samples) == 1 {
			days = firstMeasurementIntervalDays
		}
		days = clampFloat(days, minMeasurementIntervalDays, calendarYear)
		at := latest.at.Add(durationDays(days))
		return &at, nil
	}

	// --- confident phase: rare drift check only ---
	driftAt := now.Add(durationDays(driftCheckIntervalDays))
	return &driftAt, &driftAt
}

// dueSoonDays returns how far ahead of a watering date we start showing the
// "due soon" state (20% of the remaining time, at least 2 days).
func dueSoonDays(due time.Time, now time.Time) float64 {
	remaining := due.Sub(now).Hours() / 24.0
	window := remaining * 0.2
	if window < 2 {
		window = 2
	}
	return window
}
