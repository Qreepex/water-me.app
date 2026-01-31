import type { Plant } from '$lib/types/api';

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

/**
 * Format days ago as human-readable text
 */
export function daysAgo(dateString: string): string {
	const days = Math.floor((Date.now() - new Date(dateString).getTime()) / (1000 * 60 * 60 * 24));
	if (days === 0) return 'Today';
	if (days === 1) return 'Yesterday';
	return `${days} days ago`;
}

/**
 * Get watering status for a plant
 */
export function getWateringStatus(plant: Plant): { text: string; color: string } {
	const last = plant.watering?.lastWatered
		? new Date(plant.watering.lastWatered).getTime()
		: Date.now();
	const interval = plant.watering?.intervalDays ?? 0;
	const days = Math.floor((Date.now() - last) / (1000 * 60 * 60 * 24));
	const daysUntilWater = interval - days;

	if (daysUntilWater <= 0) return { text: '🌵 Needs water!', color: 'text-red-600' };
	if (daysUntilWater <= 1) return { text: '⚠️ Water soon', color: 'text-yellow-600' };
	return { text: `✓ In ${daysUntilWater} days`, color: 'text-green-600' };
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
 * Get days until water is needed
 */
export function getDaysUntilWater(plant: Plant): number {
	const lastWatered = plant.watering?.lastWatered
		? new Date(plant.watering.lastWatered).getTime()
		: 0;
	const interval = plant.watering?.intervalDays ?? 0;
	const daysSinceWatered = Math.floor((Date.now() - lastWatered) / (1000 * 60 * 60 * 24));
	return interval - daysSinceWatered;
}

/**
 * Get watering status for water page
 */
export function getPlantWaterStatus(plant: Plant): 'overdue' | 'due-soon' | 'ok' {
	const daysUntil = getDaysUntilWater(plant);
	if (daysUntil <= 0) return 'overdue';
	if (daysUntil <= 1) return 'due-soon';
	return 'ok';
}

/**
 * Get status text for water page
 */
export function getPlantStatusText(plant: Plant): string {
	const daysUntil = getDaysUntilWater(plant);
	if (daysUntil < 0)
		return `${Math.abs(daysUntil)} ${Math.abs(daysUntil) === 1 ? 'day' : 'days'} overdue`;
	if (daysUntil === 0) return 'Due today';
	if (daysUntil === 1) return 'Due tomorrow';
	return `Due in ${daysUntil} days`;
}

/**
 * Get status icon for water page
 */
export function getStatusIcon(status: 'overdue' | 'due-soon' | 'ok'): string {
	switch (status) {
		case 'overdue':
			return '🚨';
		case 'due-soon':
			return '⚠️';
		default:
			return '✅';
	}
}
