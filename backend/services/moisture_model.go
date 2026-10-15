package services

import (
	"math"
	"sort"
	"time"

	"github.com/qreepex/water-me-app/backend/types"
)

// Per-plant moisture model.
//
// The plant is modelled as a simple physical system: soil starts near
// saturation right after watering and then dries out at a roughly constant
// rate. Two numbers describe it completely:
//
//	saturation (0..1)  – moisture right after watering
//	decayPerDay (0..1) – moisture lost per day
//
// With those, the moisture at any moment is  saturation - decayPerDay * days,
// which answers both questions the app cares about: "how wet is it now?" and
// "when does it reach the lower edge of the target band?".
//
// Fitting uses every cycle ever recorded, not just the newest one, so a plant
// that was calibrated over 2-3 waterings keeps predicting even if the user
// stops measuring for weeks. New measurements always feed back in.
//
// All internal maths is done in normalized 0..1 units so the same code works
// for 1-10 meters and percentage sensors.

const (
	// A cycle with fewer usable points than this cannot contribute a rate.
	minCyclePoints = 2
	// Readings closer together than this cannot produce a per-day rate.
	minTrendSpanDays = 0.125
	// A gap larger than this ends a cycle even without a watering.
	maxCycleGapDays = 21.0
	// A rise of at least this much of the scale is treated as a watering.
	wateringEventRiseRatio = 0.15
	// Cycles older than this decay in weight: a plant's water needs change with
	// season and repotting, so recent behaviour matters more.
	cycleHalfLifeDays = 120.0

	// Cadence: while the model is still being calibrated we ask comparatively
	// often; once it is confident we only ask for a rare drift check.
	defaultMeasurementIntervalDays = 2.0
	minMeasurementIntervalDays     = 1.0
	firstMeasurementIntervalDays   = 3.0
	// Once the model is ready, confine the user's effort to a rare confirmation
	// so the prediction can be corrected if reality has drifted.
	driftCheckIntervalDays = 30.0

	// Confidence thresholds.
	// Two solid cycles are enough for a first usable prediction.
	minCyclesForPrediction = 2
	// Below this the plan asks for more data instead of predicting.
	minConfidenceToPredict = 0.45

	// A real measurement taken within this window after watering is used as the
	// observed saturation, replacing the model's prediction for that event.
	PostWateringMeasurementWindowDays = 0.5
)

// moistureSample is a single normalized reading.
type moistureSample struct {
	at    time.Time
	value float64 // 0..1
	raw   float64 // original scale units
}

// moistureCycle is the span between two waterings, with its readings.
type moistureCycle struct {
	// start is when the cycle began (the watering, or the first reading).
	start time.Time
	// wateringAt is set when the cycle start is a recorded watering.
	wateringAt *time.Time
	samples    []moistureSample
	weight     float64
}

// cycleContribution is one cycle's contribution to the pooled fit.
type cycleContribution struct {
	rate    float64 // positive drying rate (0..1/day)
	inter   float64 // saturation estimate at cycle start (0..1)
	weight  float64
	at      time.Time
	nPoints int
}

// moistureModel is the fitted per-plant model.
type moistureModel struct {
	ready       bool
	saturation  float64 // 0..1 right after watering
	decayPerDay float64 // 0..1 lost per day
	cyclesUsed  int
	confidence  float64 // 0..1
	// fittedAt is the end of the newest cycle the fit is anchored at.
	fittedAt time.Time
	// residual is the typical error of the fit, in 0..1 units.
	residual float64
}

// slopePerDay returns the least-squares drying slope (normalized per day).
func (c moistureCycle) slopePerDay() (float64, bool) {
	if len(c.samples) < minCyclePoints {
		return 0, false
	}
	span := c.samples[len(c.samples)-1].at.Sub(c.samples[0].at).Hours() / 24.0
	if span < minTrendSpanDays {
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

// interceptAt returns the fitted value at the cycle *start* (the watering, or
// the first reading when there is none). Extrapolating back to the real start
// is what gives a usable saturation estimate: the first reading may be taken
// hours or days after watering, by which time the soil has already dried.
func (c moistureCycle) interceptAt() (float64, bool) {
	slope, ok := c.slopePerDay()
	if !ok {
		return 0, false
	}
	origin := c.samples[0].at
	var sumY, sumX float64
	for _, s := range c.samples {
		sumX += s.at.Sub(origin).Hours() / 24.0
		sumY += s.value
	}
	n := float64(len(c.samples))
	interceptAtFirstReading := sumY/n - slope*(sumX/n)

	// Shift from the first reading back to the cycle start.
	days := origin.Sub(c.start).Hours() / 24.0
	if days <= 0 {
		return interceptAtFirstReading, true
	}
	return interceptAtFirstReading - slope*days, true
}

// buildCycles splits readings into drying cycles.
//
// A cycle starts at the most recent watering at or before its readings, so the
// cycle's origin is the true start of drying. When no watering is recorded the
// cycle degrades gracefully: it starts at the first reading after a detected
// event (a clear rise, or a gap long enough to imply a watering happened).
func buildCycles(samples []moistureSample, waterings []types.WateringEvent) []moistureCycle {
	if len(samples) == 0 {
		return nil
	}

	starts := make([]time.Time, 0, len(waterings))
	for _, w := range waterings {
		if !w.Date.IsZero() {
			starts = append(starts, w.Date)
		}
	}
	sort.Slice(starts, func(i, j int) bool { return starts[i].Before(starts[j]) })

	// latestWateringBefore returns the newest watering at or before `at`.
	latestWateringBefore := func(at time.Time) *time.Time {
		var found *time.Time
		for i := range starts {
			if starts[i].After(at) {
				break
			}
			t := starts[i]
			found = &t
		}
		return found
	}

	cycles := make([]moistureCycle, 0)
	current := moistureCycle{
		samples:    []moistureSample{samples[0]},
		wateringAt: latestWateringBefore(samples[0].at),
	}

	for i := 1; i < len(samples); i++ {
		prev := samples[i-1]
		next := samples[i]
		gapDays := next.at.Sub(prev.at).Hours() / 24.0
		rise := next.value - prev.value

		// A watering between the two readings starts a new cycle there.
		newWatering := false
		w := latestWateringBefore(next.at)
		if w != nil && w.After(prev.at) {
			newWatering = true
		}
		// Without a recorded watering, infer the event from a clear rise or a
		// gap long enough that drying certainly restarted.
		inferred := w == nil && (rise >= wateringEventRiseRatio || gapDays > maxCycleGapDays)

		if newWatering || inferred {
			cycles = append(cycles, current)
			anchor := w
			if inferred {
				t := next.at
				anchor = &t
			}
			current = moistureCycle{samples: []moistureSample{next}, wateringAt: anchor}
			continue
		}
		current.samples = append(current.samples, next)
	}
	cycles = append(cycles, current)

	// Set origins and recency weights.
	newest := samples[len(samples)-1].at
	for i := range cycles {
		if cycles[i].wateringAt != nil {
			cycles[i].start = *cycles[i].wateringAt
		} else {
			cycles[i].start = cycles[i].samples[0].at
		}
		ageDays := newest.Sub(cycles[i].start).Hours() / 24.0
		if ageDays < 0 {
			ageDays = 0
		}
		cycles[i].weight = math.Pow(0.5, ageDays/cycleHalfLifeDays)
	}
	return cycles
}

// effectiveWaterings returns the waterings to anchor cycles on.
//
// Plants created before the history existed only carry `lastWatered`, so that
// timestamp is folded in as a synthetic event. Without it their first cycle
// would be anchored at the first reading instead of at the watering, which
// shifts the saturation estimate and skews the predicted date.
func effectiveWaterings(watering *types.WateringConfig) []types.WateringEvent {
	if watering == nil {
		return nil
	}
	out := make([]types.WateringEvent, 0, len(watering.History)+1)
	out = append(out, watering.History...)
	if watering.LastWatered != nil {
		known := false
		for _, w := range watering.History {
			if w.Date.Equal(*watering.LastWatered) {
				known = true
				break
			}
		}
		if !known {
			out = append(out, types.WateringEvent{Date: *watering.LastWatered})
		}
	}
	return out
}

// fitModel pools every usable cycle into a single model.
//
// The drying rate is a weighted median of the per-cycle rates (robust against
// one bad cycle), and saturation is a weighted mean of the per-cycle
// intercepts. Remaining spread between cycles drives the confidence estimate.
func fitModel(samples []moistureSample, waterings []types.WateringEvent, now time.Time) moistureModel {
	cycles := buildCycles(samples, waterings)

	parts := make([]cycleContribution, 0, len(cycles))
	for _, c := range cycles {
		rate, ok := c.slopePerDay()
		if !ok {
			continue
		}
		inter, ok := c.interceptAt()
		if !ok {
			continue
		}
		// Only drying cycles (negative slope) describe the plant.
		if rate >= 0 {
			continue
		}
		// Weight by recency and by how much data the cycle contains, so a cycle
		// with many readings counts more than a two-point one.
		w := c.weight * math.Sqrt(float64(len(c.samples)))
		parts = append(parts, cycleContribution{
			rate:    -rate,
			inter:   clamp01(inter),
			weight:  w,
			at:      c.start,
			nPoints: len(c.samples),
		})
	}

	model := moistureModel{fittedAt: now, cyclesUsed: len(parts)}
	if len(parts) == 0 {
		return model
	}

	// ---- weighted median for the rate (robust to outliers) ----
	rates := make([]float64, 0, len(parts))
	weights := make([]float64, 0, len(parts))
	for _, p := range parts {
		rates = append(rates, p.rate)
		weights = append(weights, p.weight)
	}
	model.decayPerDay = weightedMedian(rates, weights)

	// ---- weighted mean for saturation ----
	var sumW, sumIW float64
	for _, p := range parts {
		sumW += p.weight
		sumIW += p.inter * p.weight
	}
	if sumW > 0 {
		model.saturation = clamp01(sumIW / sumW)
	}

	// ---- residual: typical deviation of the pooled cycles from the model ----
	var errW float64
	for _, p := range parts {
		d := math.Abs(p.rate - model.decayPerDay)
		errW += d * p.weight
	}
	if sumW > 0 {
		model.residual = errW / sumW
	}

	model.confidence = confidenceFor(parts, model)
	model.ready = len(parts) >= minCyclesForPrediction && model.confidence >= minConfidenceToPredict
	model.fittedAt = parts[len(parts)-1].at
	return model
}

// confidenceFor scores the model 0..1 from cycle count, agreement between the
// cycles and how much data each cycle carries.
func confidenceFor(parts []cycleContribution, model moistureModel) float64 {
	// 1. Amount of evidence: 2 cycles = 0.5, 3 = 0.75, 4+ = ~1.
	countScore := math.Min(1, float64(len(parts))/4.0)
	if len(parts) >= minCyclesForPrediction {
		// Guarantee a usable score as soon as the minimum is reached.
		countScore = math.Max(countScore, 0.5)
	}

	// 2. Agreement: cycles that all dry at a similar rate are trustworthy.
	// `residual` is in 0..1 scale units per day, so relate it to the rate.
	agreement := 1.0
	if model.decayPerDay > 1e-6 {
		rel := model.residual / model.decayPerDay
		agreement = clamp01(1 - rel)
	}

	// 3. Data richness: more readings per cycle = better fit.
	avgPoints := 0.0
	for _, p := range parts {
		avgPoints += float64(p.nPoints)
	}
	avgPoints /= float64(len(parts))
	richness := math.Min(1, avgPoints/4.0)

	return clamp01(0.5*countScore + 0.3*agreement + 0.2*richness)
}

// weightedMedian returns the value at which the cumulative weight reaches half.
func weightedMedian(values, weights []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	idx := make([]int, len(values))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return values[idx[a]] < values[idx[b]] })
	total := 0.0
	for _, w := range weights {
		total += w
	}
	acc := 0.0
	for _, i := range idx {
		acc += weights[i]
		if acc >= total/2 {
			return values[i]
		}
	}
	return values[idx[len(idx)-1]]
}

// moistureAt evaluates the model at a moment.
func (m moistureModel) moistureAt(at time.Time) float64 {
	days := at.Sub(m.fittedAt).Hours() / 24.0
	return clamp01(m.saturation - m.decayPerDay*days)
}

// daysUntil evaluates when the model reaches `target` after `from`.
// Returns 0 if already at or below it, and ok=false when the model has no rate.
func (m moistureModel) daysUntil(from time.Time, target float64) (float64, bool) {
	if m.decayPerDay <= 1e-9 {
		return 0, false
	}
	current := m.moistureAt(from)
	if current <= target {
		return 0, true
	}
	return (current - target) / m.decayPerDay, true
}

// saturationAt gives the moisture right after a watering at `at`, restricted to
// the values the model can physically reach.
func (m moistureModel) saturationAt() float64 {
	return clamp01(m.saturation)
}

// daysUntilFromSaturation is daysUntil for a just-watered plant: the soil starts
// at saturation rather than on the extrapolated curve, because watering refills
// it (the old line described the previous, drier cycle).
func (m moistureModel) daysUntilFromSaturation(at time.Time, target float64) (float64, bool) {
	if m.decayPerDay <= 1e-9 {
		return 0, false
	}
	if m.saturation <= target {
		return 0, true
	}
	return (m.saturation - target) / m.decayPerDay, true
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

func maxSampleValue(samples []moistureSample) float64 {
	if len(samples) == 0 {
		return 0
	}
	max := samples[0].value
	for _, s := range samples {
		if s.value > max {
			max = s.value
		}
	}
	return max
}

func durationDays(days float64) time.Duration {
	if days < 0 {
		days = 0
	}
	return time.Duration(days * 24 * float64(time.Hour))
}

// daysUntilValue is a small helper for raw (unnormalized) arithmetic.
func daysUntilValue(from, to, perDay float64) float64 {
	if perDay <= 0 {
		return 0
	}
	days := (from - to) / perDay
	if days < 0 {
		return 0
	}
	return days
}

func floatPtr(v float64) *float64 { return &v }

// --- scale helpers -------------------------------------------------------

// moistureScaleBounds returns the min/max value of a moisture scale.
func moistureScaleBounds(scale types.MoistureScale) (float64, float64) {
	if scale == types.ScalePercent {
		return 0, 100
	}
	// Analog meters usually go from 1 (dry) to 10 (soaking wet).
	return 1, 10
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

// normalizeBand converts a band width in scale units to a 0..1 fraction.
func normalizeBand(width float64, scale types.MoistureScale) float64 {
	min, max := moistureScaleBounds(scale)
	if max <= min {
		return 0
	}
	return width / (max - min)
}

// DefaultMoistureBand is half the width of the target moisture band, in scale
// units. `idealMoisture` is the centre of that band: water once the soil falls
// to ideal-Band, and water to about ideal+Band.
const DefaultMoistureBand = 2.0
