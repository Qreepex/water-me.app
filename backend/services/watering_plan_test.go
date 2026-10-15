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
	if plan.Status != types.PlanStatusWaterNow {
		t.Fatalf("status = %q, want water_now", plan.Status)
	}
	if plan.WateringDueAt == nil || plan.WateringDueAt.After(now) {
		t.Fatalf("wateringDueAt = %v, want <= now", plan.WateringDueAt)
	}
}

// Above the ideal value the plan predicts the moment the soil reaches it.
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
	// 7.5 -> 3 at one unit/day => 4.5 days after the last reading.
	days := plan.WateringDueAt.Sub(logs[2].Date).Hours() / 24
	if days < 4 || days > 5 {
		t.Fatalf("predicted watering %.1f days after the last reading, want ~4.5", days)
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
	cycles := buildCycles(samples)
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
	if days < minMeasurementIntervalDays-1e-9 || days > maxMeasurementIntervalDays+1e-9 {
		t.Fatalf("measurement in %.2f days, want within [%.0f, %.0f]",
			days, minMeasurementIntervalDays, maxMeasurementIntervalDays)
	}
}
