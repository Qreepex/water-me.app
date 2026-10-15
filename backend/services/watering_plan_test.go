package services

import (
	"testing"
	"time"

	"github.com/qreepex/water-me-app/backend/types"
)

func tenScale(v float64) *float64 { return &v }
func planMode(m types.WateringPlanMode) *types.WateringPlanMode {
	return &m
}

func daysAgo(now time.Time, days float64) time.Time {
	return now.Add(-time.Duration(days * 24 * float64(time.Hour)))
}

func TestManualPlanUsesInterval(t *testing.T) {
	now := time.Date(2026, time.March, 1, 12, 0, 0, 0, time.UTC)
	lastWatered := now.AddDate(0, 0, -3)
	plant := types.Plant{
		Watering: &types.WateringConfig{
			IntervalDays: 7,
			LastWatered:  &lastWatered,
		},
	}

	plan := ComputeWateringPlan(plant, now)
	if plan == nil {
		t.Fatal("expected a plan")
	}
	if plan.Mode != types.PlanModeManual {
		t.Fatalf("mode = %q, want Manual", plan.Mode)
	}
	if plan.WateringDueAt == nil {
		t.Fatal("expected wateringDueAt")
	}
	want := lastWatered.AddDate(0, 0, 7)
	if !plan.WateringDueAt.Equal(want) {
		t.Fatalf("wateringDueAt = %v, want %v", plan.WateringDueAt, want)
	}
	if plan.Status != types.PlanStatusOK {
		t.Fatalf("status = %q, want ok", plan.Status)
	}
}

func TestManualPlanOverdue(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)
	lastWatered := now.AddDate(0, 0, -8)
	plant := types.Plant{
		Watering: &types.WateringConfig{IntervalDays: 7, LastWatered: &lastWatered},
	}
	plan := ComputeWateringPlan(plant, now)
	if plan.Status != types.PlanStatusWaterNow {
		t.Fatalf("status = %q, want water_now", plan.Status)
	}
}

func TestManualPlanNoIntervalReturnsNil(t *testing.T) {
	plant := types.Plant{Watering: &types.WateringConfig{}}
	if plan := ComputeWateringPlan(plant, time.Now()); plan != nil {
		t.Fatalf("expected nil plan, got %+v", plan)
	}
}

func TestMoisturePlanNeedsFirstMeasurement(t *testing.T) {
	now := time.Date(2026, time.March, 1, 12, 0, 0, 0, time.UTC)
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(5),
			MoistureScale: scalePtr(types.ScaleTen),
		},
	}
	plan := ComputeWateringPlan(plant, now)
	if plan == nil {
		t.Fatal("expected a plan")
	}
	if !plan.NeedsFirstMeasurement {
		t.Error("expected needsFirstMeasurement = true")
	}
	if plan.Status != types.PlanStatusMeasure {
		t.Errorf("status = %q, want measure", plan.Status)
	}
	if plan.MeasurementDueAt == nil {
		t.Error("expected measurementDueAt to be set")
	}
}

func scalePtr(s types.MoistureScale) *types.MoistureScale { return &s }

// Reproduces the scenario from the feature request: measurements 6, 5, 4 on
// consecutive days with an ideal value of 5 => the plant is due right away.
func TestMoisturePlanScenarioFromRequest(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	logs := []types.MoistureLog{
		{ID: "1", Value: 6, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
		{ID: "2", Value: 5, Scale: types.ScaleTen, Date: daysAgo(now, 1)},
		{ID: "3", Value: 4, Scale: types.ScaleTen, Date: now},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(5),
			MoistureScale: scalePtr(types.ScaleTen),
		},
		MoistureLogs: logs,
	}

	plan := ComputeWateringPlan(plant, now)
	if plan == nil {
		t.Fatal("expected a plan")
	}
	if plan.CurrentMoisture == nil || *plan.CurrentMoisture != 4 {
		t.Fatalf("currentMoisture = %v, want 4", plan.CurrentMoisture)
	}
	if plan.DecayPerDay == nil {
		t.Fatal("expected a decay estimate")
	}
	// 6 -> 4 over two days is roughly one unit per day.
	if got := *plan.DecayPerDay; got < 0.8 || got > 1.2 {
		t.Fatalf("decayPerDay = %v, want ~1", got)
	}
	// ideal 5 with the default band of 2 means: water until 7, water again at 3.
	// At 4 the soil is still inside the band, so it is not due yet.
	if plan.Status != types.PlanStatusDueSoon {
		t.Fatalf("status = %q, want due_soon (still inside the band)", plan.Status)
	}
	if plan.WaterBelow == nil || *plan.WaterBelow != 3 {
		t.Fatalf("waterBelow = %v, want 3", plan.WaterBelow)
	}
	if plan.WaterUntil == nil || *plan.WaterUntil != 7 {
		t.Fatalf("waterUntil = %v, want 7", plan.WaterUntil)
	}
	if plan.WateringDueAt == nil || !plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want in the future (one unit/day to 3)", plan.WateringDueAt)
	}
}

// Once the soil drops to the lower edge of the band, watering is due.
func TestMoisturePlanWatersAtLowerBandEdge(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(5),
			MoistureScale: scalePtr(types.ScaleTen),
		},
		MoistureLogs: []types.MoistureLog{
			{ID: "1", Value: 7, Scale: types.ScaleTen, Date: daysAgo(now, 4)},
			{ID: "2", Value: 4, Scale: types.ScaleTen, Date: daysAgo(now, 1)},
			{ID: "3", Value: 3, Scale: types.ScaleTen, Date: now},
		},
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.Status != types.PlanStatusWaterNow {
		t.Fatalf("status = %q, want water_now at the lower band edge", plan.Status)
	}
}

// A custom band widens or narrows the hysteresis window.
func TestMoisturePlanCustomBand(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	newPlant := func(band float64) types.Plant {
		return types.Plant{
			Watering: &types.WateringConfig{
				Mode:          planMode(types.PlanModeMoisture),
				IdealMoisture: tenScale(5),
				MoistureBand:  &band,
				MoistureScale: scalePtr(types.ScaleTen),
			},
			MoistureLogs: []types.MoistureLog{
				{ID: "1", Value: 7, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
				{ID: "2", Value: 5, Scale: types.ScaleTen, Date: daysAgo(now, 1)},
				{ID: "3", Value: 4, Scale: types.ScaleTen, Date: now},
			},
		}
	}

	// band 1 => trigger at 4, which is reached right now.
	narrow := ComputeWateringPlan(newPlant(1), now)
	if narrow.WaterBelow == nil || *narrow.WaterBelow != 4 {
		t.Fatalf("band 1: waterBelow = %v, want 4", narrow.WaterBelow)
	}
	if narrow.Status != types.PlanStatusWaterNow {
		t.Fatalf("band 1: status = %q, want water_now", narrow.Status)
	}

	// band 2 => trigger at 3, still 1 unit away.
	wide := ComputeWateringPlan(newPlant(2), now)
	if wide.WaterBelow == nil || *wide.WaterBelow != 3 {
		t.Fatalf("band 2: waterBelow = %v, want 3", wide.WaterBelow)
	}
	if wide.Status == types.PlanStatusWaterNow {
		t.Fatalf("band 2: status = %q, want not water_now", wide.Status)
	}
}

// Above the band the plan predicts the moment the soil falls to the lower edge.
func TestMoisturePlanPredictsFutureWatering(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	logs := []types.MoistureLog{
		{ID: "1", Value: 9.5, Scale: types.ScaleTen, Date: daysAgo(now, 4)},
		{ID: "2", Value: 8.5, Scale: types.ScaleTen, Date: daysAgo(now, 3)},
		{ID: "3", Value: 7.5, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(3),
			MoistureScale: scalePtr(types.ScaleTen),
		},
		MoistureLogs: logs,
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.WateringDueAt == nil {
		t.Fatal("expected a predicted watering date")
	}
	if !plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want in the future", plan.WateringDueAt)
	}
	// Ideal 3 with the default band of 2 triggers at 1, so 7.5 -> 1 at one
	// unit/day is ~6.5 days after the last reading.
	if plan.WaterBelow == nil || *plan.WaterBelow != 1 {
		t.Fatalf("waterBelow = %v, want 1", plan.WaterBelow)
	}
	days := plan.WateringDueAt.Sub(logs[2].Date).Hours() / 24
	if days < 6 || days > 7 {
		t.Fatalf("predicted watering %.1f days after the last reading, want ~6.5", days)
	}
	if !plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want in the future", plan.WateringDueAt)
	}
	if plan.MeasurementDueAt == nil {
		t.Fatal("expected a next measurement date")
	}
	if !plan.MeasurementDueAt.After(logs[2].Date) {
		t.Fatalf("measurementDueAt = %v, want after last reading", plan.MeasurementDueAt)
	}
	if plan.MeasurementDueAt.After(*plan.WateringDueAt) {
		t.Error("measurement should be requested before watering")
	}
	if plan.Status != types.PlanStatusOK && plan.Status != types.PlanStatusDueSoon {
		t.Fatalf("status = %q, want ok or due_soon", plan.Status)
	}
}

// A big rise in moisture means the user watered; the cycle restarts.
func TestMoisturePlanDetectsWateringEvent(t *testing.T) {
	now := time.Date(2026, time.March, 6, 12, 0, 0, 0, time.UTC)
	logs := []types.MoistureLog{
		{ID: "1", Value: 8, Scale: types.ScaleTen, Date: daysAgo(now, 5)},
		{ID: "2", Value: 6, Scale: types.ScaleTen, Date: daysAgo(now, 4)},
		{ID: "3", Value: 3, Scale: types.ScaleTen, Date: daysAgo(now, 3)},
		// Watered, then measured again.
		{ID: "4", Value: 9, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
		{ID: "5", Value: 8, Scale: types.ScaleTen, Date: daysAgo(now, 1)},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(4),
			MoistureScale: scalePtr(types.ScaleTen),
		},
		MoistureLogs: logs,
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.DecayPerDay == nil {
		t.Fatal("expected a decay estimate")
	}
	if got := *plan.DecayPerDay; got <= 0 {
		t.Fatalf("decayPerDay = %v, want > 0", got)
	}
	if plan.SampleCount != len(logs) {
		t.Fatalf("sampleCount = %d, want %d", plan.SampleCount, len(logs))
	}
	// Latest value is far above ideal => watering should be in the future.
	if plan.WateringDueAt == nil || !plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want in the future", plan.WateringDueAt)
	}
	if plan.Status == types.PlanStatusWaterNow {
		t.Error("plant should not be due right after watering")
	}
}

func TestMoisturePlanPercentScale(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	logs := []types.MoistureLog{
		{ID: "1", Value: 60, Scale: types.ScalePercent, Date: daysAgo(now, 2)},
		{ID: "2", Value: 45, Scale: types.ScalePercent, Date: daysAgo(now, 1)},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(30),
			MoistureScale: scalePtr(types.ScalePercent),
		},
		MoistureLogs: logs,
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.MoistureScale != types.ScalePercent {
		t.Fatalf("scale = %q, want Percent", plan.MoistureScale)
	}
	if plan.WateringDueAt == nil {
		t.Fatal("expected a watering date")
	}
	// 60 -> 45 is 15%/day; from 45% it takes one more day to reach 30%.
	days := plan.WateringDueAt.Sub(logs[1].Date).Hours() / 24
	if days < 0.5 || days > 1.5 {
		t.Fatalf("predicted watering %.2f days after the last reading, want ~1", days)
	}
}

func TestMoisturePlanAfterWateringRestartsCycle(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	lastWatered := daysAgo(now, 0.5)
	logs := []types.MoistureLog{
		{ID: "1", Value: 8, Scale: types.ScaleTen, Date: daysAgo(now, 3)},
		{ID: "2", Value: 6, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(4),
			MoistureScale: scalePtr(types.ScaleTen),
			LastWatered:   &lastWatered,
		},
		MoistureLogs: logs,
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.WateringDueAt == nil {
		t.Fatal("expected a predicted watering date")
	}
	if !plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want in the future", plan.WateringDueAt)
	}
}

func TestBuildCyclesSplitsOnRise(t *testing.T) {
	now := time.Now()
	samples := []moistureSample{
		{at: daysAgo(now, 4), value: 0.8},
		{at: daysAgo(now, 3), value: 0.6},
		{at: daysAgo(now, 2), value: 0.9}, // watered
		{at: daysAgo(now, 1), value: 0.7},
	}
	cycles := buildCycles(samples, nil)
	if len(cycles) != 2 {
		t.Fatalf("got %d cycles, want 2", len(cycles))
	}
}

func TestNormalizeMoisture(t *testing.T) {
	cases := []struct {
		value float64
		scale types.MoistureScale
		want  float64
	}{
		{1, types.ScaleTen, 0},
		{10, types.ScaleTen, 1},
		{5.5, types.ScaleTen, 0.5},
		{0, types.ScalePercent, 0},
		{100, types.ScalePercent, 1},
		{25, types.ScalePercent, 0.25},
		{999, types.ScaleTen, 1}, // clamped
	}
	for _, c := range cases {
		if got := normalizeMoisture(c.value, c.scale); got != c.want {
			t.Errorf("normalizeMoisture(%v, %v) = %v, want %v", c.value, c.scale, got, c.want)
		}
	}
}

func TestNextMeasurementIsBounded(t *testing.T) {
	now := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)
	logs := []types.MoistureLog{
		{ID: "1", Value: 9, Scale: types.ScaleTen, Date: daysAgo(now, 3)},
		{ID: "2", Value: 6, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
		{ID: "3", Value: 3, Scale: types.ScaleTen, Date: daysAgo(now, 1)},
	}
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(2),
			MoistureScale: scalePtr(types.ScaleTen),
		},
		MoistureLogs: logs,
	}
	plan := ComputeWateringPlan(plant, now)
	if plan.MeasurementDueAt == nil {
		t.Fatal("expected a next measurement date")
	}
	lastReading := logs[2].Date
	days := plan.MeasurementDueAt.Sub(lastReading).Hours() / 24
	if days < minMeasurementIntervalDays-1e-9 || days > driftCheckIntervalDays+1e-9 {
		t.Fatalf("measurement in %.2f days, want within [%.0f, %.0f]",
			days, minMeasurementIntervalDays, driftCheckIntervalDays)
	}
}

// Regression: after the very first reading the app must not immediately ask for
// another one. The next measurement has to be at least a day away, otherwise the
// drying rate would be derived from readings minutes apart.
func TestFirstMeasurementSchedulesNextAtLeastADayLater(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)
	lastWatered := now
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(5),
			MoistureScale: scalePtr(types.ScaleTen),
			LastWatered:   &lastWatered,
		},
		MoistureLogs: []types.MoistureLog{
			{ID: "1", Value: 6, Scale: types.ScaleTen, Date: now},
		},
	}

	plan := ComputeWateringPlan(plant, now)
	if plan == nil {
		t.Fatal("expected a plan")
	}
	if plan.Status == types.PlanStatusMeasure {
		t.Errorf("status = measure right after the first reading, want no prompt")
	}
	if plan.WateringDueAt != nil {
		t.Errorf("wateringDueAt = %v, want nil (no trend yet)", plan.WateringDueAt)
	}
	if plan.MeasurementDueAt == nil {
		t.Fatal("expected a next measurement date")
	}
	gap := plan.MeasurementDueAt.Sub(now)
	if gap < 24*time.Hour {
		t.Fatalf("next measurement in %v, want at least 24h", gap)
	}
}

// Regression: readings taken far too close together must not produce a drying
// rate, no matter how large the drop between them.
func TestReadingsMinutesApartDoNotProduceDecay(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)
	for _, gap := range []time.Duration{0, time.Second, time.Minute, time.Hour} {
		lastWatered := now
		plant := types.Plant{
			Watering: &types.WateringConfig{
				Mode:          planMode(types.PlanModeMoisture),
				IdealMoisture: tenScale(5),
				MoistureScale: scalePtr(types.ScaleTen),
				LastWatered:   &lastWatered,
			},
			MoistureLogs: []types.MoistureLog{
				{ID: "1", Value: 8, Scale: types.ScaleTen, Date: now.Add(-gap)},
				{ID: "2", Value: 4, Scale: types.ScaleTen, Date: now},
			},
		}

		plan := ComputeWateringPlan(plant, now)
		if plan.DecayPerDay != nil {
			t.Errorf("gap %v: decayPerDay = %v, want nil", gap, *plan.DecayPerDay)
		}
		if plan.Status == types.PlanStatusWaterNow {
			t.Errorf("gap %v: plant reported as due although no rate is known", gap)
		}
		if plan.MeasurementDueAt == nil || plan.MeasurementDueAt.Sub(now) < 24*time.Hour {
			t.Errorf("gap %v: next measurement not at least 24h away", gap)
		}
	}
}

// A measurement due date must never be in the past for a freshly measured plant.
func TestMeasurementDueAtIsInTheFutureAfterMeasuring(t *testing.T) {
	now := time.Date(2026, time.March, 10, 12, 0, 0, 0, time.UTC)
	lastWatered := daysAgo(now, 0.5)
	plant := types.Plant{
		Watering: &types.WateringConfig{
			Mode:          planMode(types.PlanModeMoisture),
			IdealMoisture: tenScale(4),
			MoistureScale: scalePtr(types.ScaleTen),
			LastWatered:   &lastWatered,
		},
		MoistureLogs: []types.MoistureLog{
			{ID: "1", Value: 9, Scale: types.ScaleTen, Date: daysAgo(now, 2)},
			{ID: "2", Value: 8, Scale: types.ScaleTen, Date: now},
		},
	}

	plan := ComputeWateringPlan(plant, now)
	if plan.MeasurementDueAt == nil {
		t.Fatal("expected a next measurement date")
	}
	if !plan.MeasurementDueAt.After(now) {
		t.Fatalf("measurementDueAt = %v, want in the future", plan.MeasurementDueAt)
	}
}
