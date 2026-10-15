package services

import (
	"math"
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
const (
	// Until we have a trend we simply ask for the next value after this many days.
	defaultMeasurementIntervalDays = 2.0
	// Never ask for a measurement more often than once per day ...
	minMeasurementIntervalDays = 1.0
	// ... and at most once a week.
	maxMeasurementIntervalDays = 7.0
	// A rise of at least this fraction of the full scale is treated as a
	// watering event (the user watered the plant => new drying cycle).
	wateringEventRiseRatio = 0.15
	// Gaps longer than this split the data into independent drying cycles.
	maxCycleGapDays = 21.0
)

// moistureScaleBounds returns the min/max value of a moisture scale.
func moistureScaleBounds(scale types.MoistureScale) (float64, float64) {
	switch scale {
	case types.ScalePercent:
		return 0, 100
	default:
		// Analog meters usually go from 1 (dry) to 10 (soaking wet).
		return 1, 10
	}
}

// normalizeMoisture converts a raw reading to a 0..1 fraction of its scale.
func normalizeMoisture(value float64, scale types.MoistureScale) float64 {
	min, max := moistureScaleBounds(scale)
	if max <= min {
		return 0
	}
	return clamp01((value - min) / (max - min))
}

// denormalizeMoisture converts a fraction back into raw scale units.
func denormalizeMoisture(fraction float64, scale types.MoistureScale) float64 {
	min, max := moistureScaleBounds(scale)
	return min + clamp01(fraction)*(max-min)
}

// denormalizeDelta converts a normalized *difference* (e.g. a drying speed)
// back into raw scale units. Unlike denormalizeMoisture it must not add the
// scale offset, because a delta has no absolute position on the scale.
func denormalizeDelta(fraction float64, scale types.MoistureScale) float64 {
	min, max := moistureScaleBounds(scale)
	return fraction * (max - min)
}

// moistureSample is a normalized reading.
type moistureSample struct {
	at    time.Time
	value float64 // 0..1
	raw   float64 // original value
}

// moistureCycle is a set of consecutive readings between two waterings.
type moistureCycle struct {
	samples []moistureSample
}

// slopePerDay returns the least-squares drying slope (normalized value per day).
// It returns ok=false when the cycle cannot produce a meaningful trend.
func (c moistureCycle) slopePerDay() (float64, bool) {
	if len(c.samples) < 2 {
		return 0, false
	}
	base := c.samples[0].at
	var sumX, sumY, sumXX, sumXY float64
	n := float64(len(c.samples))
	for _, s := range c.samples {
		x := s.at.Sub(base).Hours() / 24.0
		y := s.value
		sumX += x
		sumY += y
		sumXX += x * x
		sumXY += x * y
	}
	denominator := n*sumXX - sumX*sumX
	if math.Abs(denominator) < 1e-9 {
		return 0, false
	}
	slope := (n*sumXY - sumX*sumY) / denominator
	if math.IsNaN(slope) || math.IsInf(slope, 0) {
		return 0, false
	}
	return slope, true
}

// buildCycles splits readings into drying cycles: a new cycle starts after a
// watering (detected as a significant rise) or after a long gap.
func buildCycles(samples []moistureSample) []moistureCycle {
	if len(samples) == 0 {
		return nil
	}

	cycles := make([]moistureCycle, 0)
	current := moistureCycle{samples: []moistureSample{samples[0]}}

	for i := 1; i < len(samples); i++ {
		prev := samples[i-1]
		next := samples[i]
		gapDays := next.at.Sub(prev.at).Hours() / 24.0
		rise := next.value - prev.value

		if gapDays > maxCycleGapDays || rise >= wateringEventRiseRatio {
			cycles = append(cycles, current)
			current = moistureCycle{samples: []moistureSample{next}}
			continue
		}
		current.samples = append(current.samples, next)
	}
	cycles = append(cycles, current)
	return cycles
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

func clamp01(v float64) float64 {
	return math.Min(1, math.Max(0, v))
}

func clampFloat(v, min, max float64) float64 {
	return math.Min(max, math.Max(min, v))
}

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
		// Misconfigured: we cannot compute anything useful yet.
		plan.NeedsFirstMeasurement = true
		plan.Status = types.PlanStatusMeasure
		plan.MeasurementDueAt = &now
		return plan
	}

	ideal := normalizeMoisture(*watering.IdealMoisture, scale)

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
		plan.NeedsFirstMeasurement = true
		plan.MeasurementDueAt = &now
		plan.Status = types.PlanStatusMeasure
		return plan
	}

	latest := samples[len(samples)-1]
	plan.CurrentMoisture = &latest.raw
	plan.MeasuredAt = &latest.at
	estimatedFull := denormalizeMoisture(maxSampleValue(samples), scale)
	plan.EstimatedFullMoisture = &estimatedFull

	decay := estimateDecayPerDay(samples)
	if decay != nil {
		perDay := denormalizeDelta(*decay, scale)
		plan.DecayPerDay = &perDay
	}

	// Predict when the soil reaches the ideal value again.
	if watering.LastWatered != nil && latest.at.Before(*watering.LastWatered) {
		// Watered after the last measurement: restart from the estimated
		// saturation value right after watering.
		if decay != nil {
			full := maxSampleValue(samples)
			days := daysUntilValue(full, ideal, *decay)
			due := watering.LastWatered.Add(durationDays(days))
			plan.WateringDueAt = &due
		}
	} else if latest.value <= ideal {
		due := now
		plan.WateringDueAt = &due
	} else if decay != nil {
		days := (latest.value - ideal) / *decay
		due := latest.at.Add(durationDays(days))
		plan.WateringDueAt = &due
	}

	// Decide when the next measurement should be taken.
	nextMeasurement := nextMeasurementAt(samples, plan.WateringDueAt, decay)
	plan.MeasurementDueAt = nextMeasurement

	switch {
	case plan.WateringDueAt == nil:
		// Not enough data to predict the watering: we need another value.
		plan.Status = types.PlanStatusMeasure
	case !now.Before(*plan.WateringDueAt):
		plan.Status = types.PlanStatusWaterNow
	case nextMeasurement != nil && !now.Before(*nextMeasurement):
		plan.Status = types.PlanStatusMeasure
	case plan.WateringDueAt.Sub(now) <= durationDays(dueSoonDays(*plan.WateringDueAt, now)):
		plan.Status = types.PlanStatusDueSoon
	default:
		plan.Status = types.PlanStatusOK
	}

	return plan
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

func maxSampleValue(samples []moistureSample) float64 {
	max := samples[0].value
	for _, s := range samples {
		if s.value > max {
			max = s.value
		}
	}
	return max
}

func daysUntilValue(from, to, decayPerDay float64) float64 {
	if decayPerDay <= 0 {
		return 0
	}
	days := (from - to) / decayPerDay
	if days < 0 {
		return 0
	}
	return days
}

func durationDays(days float64) time.Duration {
	if days < 0 {
		days = 0
	}
	return time.Duration(days * 24 * float64(time.Hour))
}

// estimateDecayPerDay returns the median drying speed (normalized value per
// day, positive) across all drying cycles. Returns nil when no trend exists.
func estimateDecayPerDay(samples []moistureSample) *float64 {
	cycles := buildCycles(samples)
	slopes := make([]float64, 0, len(cycles))
	for _, cycle := range cycles {
		slope, ok := cycle.slopePerDay()
		if !ok {
			continue
		}
		if slope < 0 {
			slopes = append(slopes, -slope)
		}
	}
	if len(slopes) == 0 {
		return nil
	}
	decay := median(slopes)
	if decay <= 1e-6 {
		return nil
	}
	return &decay
}

// nextMeasurementAt picks the next reminder date for a fresh reading.
func nextMeasurementAt(
	samples []moistureSample,
	wateringDueAt *time.Time,
	decay *float64,
) *time.Time {
	latest := samples[len(samples)-1]

	interval := defaultMeasurementIntervalDays
	if decay != nil && wateringDueAt != nil {
		daysUntilWatering := wateringDueAt.Sub(latest.at).Hours() / 24.0
		if daysUntilWatering > 0 {
			// Sample twice on the way down to the ideal value.
			interval = daysUntilWatering / 2.0
		} else {
			interval = minMeasurementIntervalDays
		}
	}
	interval = clampFloat(interval, minMeasurementIntervalDays, maxMeasurementIntervalDays)

	next := latest.at.Add(durationDays(interval))
	// Never ask later than the predicted watering moment.
	if wateringDueAt != nil && next.After(*wateringDueAt) {
		next = *wateringDueAt
	}
	// Readings are per day, so don't schedule before "tomorrow".
	earliest := latest.at.Add(durationDays(minMeasurementIntervalDays))
	if next.Before(earliest) {
		next = earliest
	}
	return &next
}
