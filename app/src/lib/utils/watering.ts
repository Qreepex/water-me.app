import type { Plant } from '$lib/types/api';
import { getDaysUntilWater } from './plant';

/**
 * Sort plants by watering priority (most urgent first)
 */
export function sortByWateringPriority(plants: Plant[]): Plant[] {
	return [...plants].sort((a, b) => {
		const aDays = getDaysUntilWater(a);
		const bDays = getDaysUntilWater(b);
		// Overdue plants first (negative numbers), then due soon, then ok
		return aDays - bDays;
	});
}

/**
 * Filter plants that need watering today or are overdue
 */
export function getPlantsNeedingWaterToday(plants: Plant[]): Plant[] {
	return plants.filter((plant) => {
		const daysUntil = getDaysUntilWater(plant);
		return daysUntil <= 0;
	});
}

/**
 * Filter plants that need watering soon (within 1-2 days)
 */
export function getPlantsNeedingWaterSoon(plants: Plant[]): Plant[] {
	return plants.filter((plant) => {
		const daysUntil = getDaysUntilWater(plant);
		return daysUntil > 0 && daysUntil <= 2;
	});
}

/**
 * Filter plants that don't need watering soon
 */
export function getPlantsOk(plants: Plant[]): Plant[] {
	return plants.filter((plant) => {
		const daysUntil = getDaysUntilWater(plant);
		return daysUntil > 2;
	});
}

/**
 * Group plants by watering status
 */
export function groupPlantsByWateringStatus(plants: Plant[]): {
	overdue: Plant[];
	dueSoon: Plant[];
	ok: Plant[];
} {
	return {
		overdue: getPlantsNeedingWaterToday(plants),
		dueSoon: getPlantsNeedingWaterSoon(plants),
		ok: getPlantsOk(plants)
	};
}
