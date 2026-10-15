import type { IconName } from '$lib/components/ui/icons';
import type { Tone } from '$lib/components/ui/tone';
import type { Plant, WateringPlan } from '$lib/types/api';

export type SortOption =
	| 'name'
	| 'lastWatered'
	| 'lastFertilized'
	| 'wateringIntervalDays'
	| 'mistingIntervalDays'
	| 'nameAsc'
	| 'nameDesc'
	| 'lastWateredAsc'
	| 'lastWateredDesc'
	| 'speciesAsc'
	| 'speciesDesc';

/** Raw moisture values are 1 (dry) .. 10 (soaking wet). */
export const MOISTURE_MIN = 1;
export const MOISTURE_MAX = 10;
/** Percentage scale bounds. */
export const MOISTURE_PCT_MIN = 0;
export const MOISTURE_PCT_MAX = 100;

/** The plant's configured watering mode ("Manual" for legacy plants). */
export function getWateringMode(plant: Plant): 'Manual' | 'Moisture' {
	return plant.watering?.mode === 'Moisture' ? 'Moisture' : 'Manual';
}

export function hasWateringConfig(plant: Plant): boolean {
	if (getWateringMode(plant) === 'Moisture') return true;
	return !!plant.watering?.intervalDays;
}

/** Server-computed schedule, when present. */
export function getPlan(plant: Plant): WateringPlan | null {
	return plant.wateringPlan ?? null;
}

/** Is the plant due for watering (either mode)? */
export function isDueForWatering(plant: Plant): boolean {
	if (isReadingBelowBand(plant)) return true;
	const plan = getPlan(plant);
	if (plan) return plan.status === 'water_now' || plan.status === 'due_soon';
	return getPlantWaterStatus(plant) !== 'ok';
}

/** Does the moisture-based plan ask for a new measurement? */
export function needsMeasurement(plant: Plant): boolean {
	return getPlan(plant)?.status === 'measure';
}

/** Scale bounds for the plant's configured moisture scale. */
export function getMoistureBounds(plant: Plant): { min: number; max: number; step: number } {
	if (plant.watering?.moistureScale === 'Percent') {
		return { min: MOISTURE_PCT_MIN, max: MOISTURE_PCT_MAX, step: 1 };
	}
	return { min: MOISTURE_MIN, max: MOISTURE_MAX, step: 0.5 };
}

/**
 * Get watering status for a plant
 */
export function getWateringStatus(
	plant: Plant
): { text: string; tone: Tone; icon: IconName; args?: string[] } | null {
	if (!hasWateringConfig(plant)) return null;

	// A reading below the band is conclusive, even before a model exists.
	if (isReadingBelowBand(plant)) {
		return { text: 'plants.moistureWaterNowStatus', tone: 'danger', icon: 'alert-triangle' };
	}

	// Moisture plants prefer the server-computed plan when available.
	const plan = getPlan(plant);
	if (plan && getWateringMode(plant) === 'Moisture') {
		switch (plan.status) {
			case 'water_now':
				return { text: 'plants.moistureWaterNowStatus', tone: 'danger', icon: 'alert-triangle' };
			case 'due_soon':
				return { text: 'plants.moistureDueSoonStatus', tone: 'warn', icon: 'droplet' };
			case 'measure':
				return { text: 'plants.moistureMeasureStatus', tone: 'info', icon: 'ruler' };
			default:
				return { text: 'plants.moistureOkStatus', tone: 'ok', icon: 'check-circle' };
		}
	}

	const last = plant.watering?.lastWatered
		? new Date(plant.watering.lastWatered).getTime()
		: Date.now();
	const interval = plant.watering?.intervalDays ?? 0;
	const days = Math.floor((Date.now() - last) / (1000 * 60 * 60 * 24));
	const daysUntilWater = interval - days;

	if (daysUntilWater <= -1) {
		return {
			text: 'plants.waterDaysOverdue',
			tone: 'danger',
			icon: 'alert-triangle',
			args: [String(Math.abs(daysUntilWater))]
		};
	}
	if (daysUntilWater === 0) {
		return { text: 'plants.dueTodayStatus', tone: 'danger', icon: 'droplet' };
	}
	if (daysUntilWater === 1) {
		return { text: 'plants.dueTomorrowStatus', tone: 'warn', icon: 'droplet' };
	}
	if (daysUntilWater > 1) {
		return {
			text: 'plants.waterInDays',
			tone: 'ok',
			icon: 'check-circle',
			args: [String(daysUntilWater)]
		};
	}
	return { text: 'plants.needsWaterStatus', tone: 'danger', icon: 'alert-triangle' };
}

/**
 * Sort plants by the given option
 */
export function sortPlants(plants: Plant[], sortBy: SortOption): Plant[] {
	const sorted = [...plants];
	switch (sortBy) {
		case 'name':
		case 'nameAsc':
			return sorted.sort((a, b) => (a.name ?? '').localeCompare(b.name ?? ''));
		case 'nameDesc':
			return sorted.sort((a, b) => (b.name ?? '').localeCompare(a.name ?? ''));
		case 'lastWatered':
		case 'lastWateredDesc':
			return sorted.sort((a, b) => {
				const aw = a.watering?.lastWatered ? new Date(a.watering.lastWatered).getTime() : 0;
				const bw = b.watering?.lastWatered ? new Date(b.watering.lastWatered).getTime() : 0;
				return bw - aw;
			});
		case 'lastWateredAsc':
			return sorted.sort((a, b) => {
				const aw = a.watering?.lastWatered ? new Date(a.watering.lastWatered).getTime() : 0;
				const bw = b.watering?.lastWatered ? new Date(b.watering.lastWatered).getTime() : 0;
				return aw - bw;
			});
		case 'speciesAsc':
			return sorted.sort((a, b) => (a.species ?? '').localeCompare(b.species ?? ''));
		case 'speciesDesc':
			return sorted.sort((a, b) => (b.species ?? '').localeCompare(a.species ?? ''));
		case 'lastFertilized':
			return sorted.sort((a, b) => {
				const af = a.fertilizing?.lastFertilized
					? new Date(a.fertilizing.lastFertilized).getTime()
					: 0;
				const bf = b.fertilizing?.lastFertilized
					? new Date(b.fertilizing.lastFertilized).getTime()
					: 0;
				return bf - af;
			});
		case 'wateringIntervalDays':
			return sorted.sort(
				(a, b) => (a.watering?.intervalDays ?? 999) - (b.watering?.intervalDays ?? 999)
			);
		case 'mistingIntervalDays':
			return sorted.sort(
				(a, b) =>
					(a.humidity?.mistingIntervalDays ?? 999) - (b.humidity?.mistingIntervalDays ?? 999)
			);
		default:
			return sorted;
	}
}

/**
 * Get days until water is needed. Prefers the server-computed plan and falls
 * back to the fixed interval for legacy/manual plants.
 */
export function getDaysUntilWater(plant: Plant): number {
	const plan = getPlan(plant);
	if (plan?.wateringDueAt) {
		const due = new Date(plan.wateringDueAt).getTime();
		return Math.floor((due - Date.now()) / (1000 * 60 * 60 * 24));
	}

	const lastWatered = plant.watering?.lastWatered
		? new Date(plant.watering.lastWatered).getTime()
		: 0;
	const interval = plant.watering?.intervalDays ?? 0;
	const daysSinceWatered = Math.floor((Date.now() - lastWatered) / (1000 * 60 * 60 * 24));
	return interval - daysSinceWatered;
}

/**
 * Is the newest reading itself already below the target band?
 *
 * The server reports this too, but the app checks independently: a single
 * reading below the band is conclusive on its own and must never surface as
 * "moisture is fine", regardless of whether a model has been fitted yet.
 * A reading taken right after watering describes saturated soil and is ignored.
 */
export function isReadingBelowBand(plant: Plant): boolean {
	const plan = getPlan(plant);
	if (!plan) return false;
	if (plan.wateredSinceLastMeasurement) return false;
	const current = plan.currentMoisture;
	const below = plan.waterBelow;
	if (current === undefined || current === null || below === undefined || below === null) {
		return false;
	}
	return current <= below;
}

/**
 * Thirst level as 0..1 progress (0 = just watered, 1 = overdue), for progress bars.
 */
export function getWateringProgress(plant: Plant): number {
	// Moisture plans: progress from "watered" to the lower edge of the target
	// band, which is where the server recommends watering.
	const plan = getPlan(plant);
	if (getWateringMode(plant) === 'Moisture' && plan) {
		const current = plan.currentMoisture;
		const below = plan.waterBelow ?? null;
		const full = plan.estimatedFullMoisture;
		if (current === undefined || current === null) return 0;
		if (below === null || full === undefined || full === null) {
			return plan.status === 'water_now' ? 1 : 0;
		}
		if (full <= below) return current <= below ? 1 : 0;
		return Math.min(1, Math.max(0, (full - current) / (full - below)));
	}

	const interval = plant.watering?.intervalDays ?? 0;
	if (!interval) return 0;
	const daysSince = Math.max(0, interval - getDaysUntilWater(plant));
	return Math.min(1, daysSince / interval);
}

/**
 * Get watering status for water page
 */
export function getPlantWaterStatus(plant: Plant): 'overdue' | 'due-soon' | 'ok' {
	if (isReadingBelowBand(plant)) return 'overdue';
	const plan = getPlan(plant);
	if (plan) {
		if (plan.status === 'water_now') return 'overdue';
		if (plan.status === 'due_soon') return 'due-soon';
		return 'ok';
	}
	const daysUntil = getDaysUntilWater(plant);
	if (daysUntil <= 0) return 'overdue';
	if (daysUntil <= 1) return 'due-soon';
	return 'ok';
}

/**
 * Get status text for water page
 */
export function getPlantStatusText(plant: Plant): { key: string; args?: string[] } {
	const plan = getPlan(plant);
	if (plan && getWateringMode(plant) === 'Moisture') {
		switch (plan.status) {
			case 'water_now':
				return { key: 'plants.statusMoistureWaterNow' };
			case 'due_soon':
				return { key: 'plants.statusMoistureDueSoon' };
			case 'measure':
				return { key: 'plants.statusMeasure' };
			default:
				return { key: 'plants.statusMoistureOk' };
		}
	}

	const daysUntil = getDaysUntilWater(plant);
	if (daysUntil < 0)
		return { key: 'plants.statusDaysOverdue', args: [String(Math.abs(daysUntil))] };
	if (daysUntil === 0) return { key: 'plants.statusDueToday' };
	if (daysUntil === 1) return { key: 'plants.statusDueTomorrow' };
	return { key: 'plants.statusDueInDays', args: [String(daysUntil)] };
}

/**
 * Get status icon for water page
 */
export function getStatusVisual(status: 'overdue' | 'due-soon' | 'ok'): {
	icon: IconName;
	tone: Tone;
} {
	switch (status) {
		case 'overdue':
			return { icon: 'alert-triangle', tone: 'danger' };
		case 'due-soon':
			return { icon: 'alert-circle', tone: 'warn' };
		default:
			return { icon: 'check-circle', tone: 'ok' };
	}
}

/**
 * Human-readable phase of the moisture curve, used for the plant detail page.
 */
export function getMoistureTrend(plant: Plant): 'rising' | 'falling' | 'flat' | null {
	const plan = getPlan(plant);
	if (!plan?.decayPerDay) return null;
	if (plan.decayPerDay > 0.05) return 'falling';
	if (plan.decayPerDay < -0.05) return 'rising';
	return 'flat';
}

/** Format a moisture value with its scale for display. */
export function formatMoisture(value: number | null | undefined, scale?: string | null): string {
	if (value === undefined || value === null) return '—';
	const rounded = Math.round(value * 10) / 10;
	return scale === 'Percent' ? `${rounded}%` : `${rounded}/10`;
}
