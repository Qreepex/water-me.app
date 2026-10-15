import type { Plant } from '$lib/types/api';
import { getMoistureBounds, getPlan, getWateringMode, hasWateringConfig } from './plant';

export type TaskKind = 'water' | 'measure';

export interface CareTask {
	kind: TaskKind;
	plant: Plant;
	/** Sort key: lower = more urgent */
	urgency: number;
	/** i18n key for the short reason */
	reasonKey: string;
}

export interface TaskList {
	/** Watering tasks, most urgent first. */
	water: CareTask[];
	/** Moisture measurement tasks, most urgent first. */
	measure: CareTask[];
	/** Plants with a plan that need nothing today. */
	relaxed: Plant[];
	/** Plants without any usable watering plan yet. */
	unconfigured: Plant[];
	total: number;
}

function daysUntil(timestamp?: string | null): number | null {
	if (!timestamp) return null;
	return (new Date(timestamp).getTime() - Date.now()) / (1000 * 60 * 60 * 24);
}

/**
 * Build the "what should I do today" list.
 *
 * Watering tasks come from the server-computed plan (`wateringPlan.status`).
 * Measurement tasks come from `measurementDueAt`; a plant whose plan is still
 * being learned (`needsFirstMeasurement`) is always a measurement task.
 */
export function buildTaskList(plants: Plant[], now = Date.now()): TaskList {
	const water: CareTask[] = [];
	const measure: CareTask[] = [];
	const relaxed: Plant[] = [];
	const unconfigured: Plant[] = [];

	for (const plant of plants) {
		const plan = getPlan(plant);

		if (!plan) {
			// No server plan: fall back to the configured mode so the list still
			// shows something useful while offline/legacy clients are around.
			if (hasWateringConfig(plant)) {
				if (getWateringMode(plant) === 'Moisture') {
					measure.push({
						kind: 'measure',
						plant,
						urgency: 1,
						reasonKey: 'tasks.reasonFirstMeasurement'
					});
				}
				continue;
			}
			unconfigured.push(plant);
			continue;
		}

		const dueIn = plan.wateringDueAt
			? (new Date(plan.wateringDueAt).getTime() - now) / 86400000
			: null;
		const measureIn = plan.measurementDueAt
			? (new Date(plan.measurementDueAt).getTime() - now) / 86400000
			: null;

		const waterNow = plan.status === 'water_now';
		// A measurement is only due once the schedule actually says so. The
		// server keeps at least a day between readings, so this must not trigger
		// just because a plant is still collecting its first values.
		const measureNow = plan.status === 'measure' || (measureIn !== null && measureIn <= 0);

		if (waterNow) {
			water.push({
				kind: 'water',
				plant,
				urgency: dueIn ?? -1,
				reasonKey: plan.mode === 'Moisture' ? 'tasks.reasonMoistureLow' : 'tasks.reasonIntervalDue'
			});
		}

		if (measureNow) {
			measure.push({
				kind: 'measure',
				plant,
				urgency: plan.needsFirstMeasurement ? -100 : (measureIn ?? 0),
				reasonKey: plan.needsFirstMeasurement
					? 'tasks.reasonFirstMeasurement'
					: 'tasks.reasonMeasurementDue'
			});
		}

		if (!waterNow && !measureNow) {
			relaxed.push(plant);
		}
	}

	water.sort((a, b) => a.urgency - b.urgency);
	measure.sort((a, b) => a.urgency - b.urgency);

	return { water, measure, relaxed, unconfigured, total: water.length + measure.length };
}

/** Days until the next measurement, or null if none is scheduled. */
export function measurementDueInDays(plant: Plant): number | null {
	const plan = getPlan(plant);
	if (!plan?.measurementDueAt) return null;
	return daysUntil(plan.measurementDueAt);
}

/** Suggested default moisture value for a new reading (last measured value). */
export function suggestedMoistureValue(plant: Plant): number {
	const plan = getPlan(plant);
	const bounds = getMoistureBounds(plant);
	const ideal = plant.watering?.idealMoisture ?? (bounds.min + bounds.max) / 2;
	if (plan?.currentMoisture === undefined || plan?.currentMoisture === null) return ideal;
	return plan.currentMoisture;
}
