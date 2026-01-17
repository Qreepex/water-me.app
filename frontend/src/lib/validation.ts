import { SunlightRequirement, PlantFlag } from './types/types';
import type { Plant } from './types/types';

export interface ValidationError {
	field: string;
	message: string;
}

// Validation constraints
const CONSTRAINTS = {
	species: { minLength: 1, maxLength: 100 },
	name: { minLength: 1, maxLength: 100 },
	temperature: { min: -50, max: 100 },
	humidity: { min: 0, max: 100 },
	wateringInterval: { min: 1, max: 365 },
	fertilizingInterval: { min: 1, max: 365 },
	sprayInterval: { min: 1, max: 365 },
	notesMaxLength: 1000,
	notes: { maxItems: 100, maxItemLength: 500 },
	photoIds: { maxItems: 100, maxIdLength: 255 } // non-data URIs capped; data: URIs allowed larger
};

export function validatePlantInput(data: unknown): { valid: boolean; errors: ValidationError[] } {
	const errors: ValidationError[] = [];

	if (!data || typeof data !== 'object') {
		return { valid: false, errors: [{ field: 'root', message: 'Request body is required' }] };
	}

	const obj = data as Record<string, unknown>;

	// Validate species
	if (typeof obj.species !== 'string' || obj.species.trim().length === 0) {
		errors.push({
			field: 'species',
			message: 'Species is required and must be a non-empty string'
		});
	} else if (obj.species.length > CONSTRAINTS.species.maxLength) {
		errors.push({
			field: 'species',
			message: `Species must be ${CONSTRAINTS.species.maxLength} characters or less`
		});
	}

	// Validate name
	if (typeof obj.name !== 'string' || obj.name.trim().length === 0) {
		errors.push({ field: 'name', message: 'Name is required and must be a non-empty string' });
	} else if (obj.name.length > CONSTRAINTS.name.maxLength) {
		errors.push({
			field: 'name',
			message: `Name must be ${CONSTRAINTS.name.maxLength} characters or less`
		});
	}

	// Validate sunLight
	if (typeof obj.sunLight !== 'string' || !isSunlightRequirement(obj.sunLight)) {
		errors.push({
			field: 'sunLight',
			message: `SunLight must be one of: ${Object.values(SunlightRequirement).join(', ')}`
		});
	}

	// Validate preferedTemperature
	if (typeof obj.preferedTemperature !== 'number') {
		errors.push({ field: 'preferedTemperature', message: 'PreferredTemperature must be a number' });
	} else if (
		obj.preferedTemperature < CONSTRAINTS.temperature.min ||
		obj.preferedTemperature > CONSTRAINTS.temperature.max
	) {
		errors.push({
			field: 'preferedTemperature',
			message: `PreferredTemperature must be between ${CONSTRAINTS.temperature.min} and ${CONSTRAINTS.temperature.max}`
		});
	}

	// Validate wateringIntervalDays
	if (typeof obj.wateringIntervalDays !== 'number' || obj.wateringIntervalDays < 1) {
		errors.push({
			field: 'wateringIntervalDays',
			message: 'WateringIntervalDays must be a positive number'
		});
	} else if (obj.wateringIntervalDays > CONSTRAINTS.wateringInterval.max) {
		errors.push({
			field: 'wateringIntervalDays',
			message: `WateringIntervalDays must be ${CONSTRAINTS.wateringInterval.max} or less`
		});
	}

	// Validate fertilizingIntervalDays
	if (typeof obj.fertilizingIntervalDays !== 'number' || obj.fertilizingIntervalDays < 1) {
		errors.push({
			field: 'fertilizingIntervalDays',
			message: 'FertilizingIntervalDays must be a positive number'
		});
	} else if (obj.fertilizingIntervalDays > CONSTRAINTS.fertilizingInterval.max) {
		errors.push({
			field: 'fertilizingIntervalDays',
			message: `FertilizingIntervalDays must be ${CONSTRAINTS.fertilizingInterval.max} or less`
		});
	}

	// Validate preferedHumidity
	if (typeof obj.preferedHumidity !== 'number') {
		errors.push({ field: 'preferedHumidity', message: 'PreferredHumidity must be a number' });
	} else if (
		obj.preferedHumidity < CONSTRAINTS.humidity.min ||
		obj.preferedHumidity > CONSTRAINTS.humidity.max
	) {
		errors.push({
			field: 'preferedHumidity',
			message: `PreferredHumidity must be between ${CONSTRAINTS.humidity.min} and ${CONSTRAINTS.humidity.max}`
		});
	}

	// Validate sprayIntervalDays (optional)
	if (obj.sprayIntervalDays !== undefined && obj.sprayIntervalDays !== null) {
		if (typeof obj.sprayIntervalDays !== 'number' || obj.sprayIntervalDays < 1) {
			errors.push({
				field: 'sprayIntervalDays',
				message: 'SprayIntervalDays must be a positive number'
			});
		} else if (obj.sprayIntervalDays > CONSTRAINTS.sprayInterval.max) {
			errors.push({
				field: 'sprayIntervalDays',
				message: `SprayIntervalDays must be ${CONSTRAINTS.sprayInterval.max} or less`
			});
		}
	}

	// Validate lastWatered (optional)
	if (obj.lastWatered !== undefined && obj.lastWatered !== null) {
		if (typeof obj.lastWatered !== 'string' || !isValidISODate(obj.lastWatered)) {
			errors.push({
				field: 'lastWatered',
				message: 'LastWatered must be a valid ISO 8601 date string'
			});
		}
	}

	// Validate lastFertilized (optional)
	if (obj.lastFertilized !== undefined && obj.lastFertilized !== null) {
		if (typeof obj.lastFertilized !== 'string' || !isValidISODate(obj.lastFertilized)) {
			errors.push({
				field: 'lastFertilized',
				message: 'LastFertilized must be a valid ISO 8601 date string'
			});
		}
	}

	// Validate notes (optional)
	if (obj.notes !== undefined && obj.notes !== null) {
		if (!Array.isArray(obj.notes)) {
			errors.push({ field: 'notes', message: 'Notes must be an array' });
		} else if (obj.notes.length > CONSTRAINTS.notes.maxItems) {
			errors.push({
				field: 'notes',
				message: `Notes array must contain ${CONSTRAINTS.notes.maxItems} items or less`
			});
		} else {
			const invalidNotes = obj.notes.filter(
				(note: unknown): boolean => typeof note !== 'string' || note.length === 0
			);
			if (invalidNotes.length > 0) {
				errors.push({ field: 'notes', message: 'All notes must be non-empty strings' });
			}
			const longNotes = obj.notes.filter(
				(note: unknown): boolean =>
					typeof note === 'string' && note.length > CONSTRAINTS.notes.maxItemLength
			);
			if (longNotes.length > 0) {
				errors.push({
					field: 'notes',
					message: `Each note must be ${CONSTRAINTS.notes.maxItemLength} characters or less`
				});
			}
		}
	}

	// Validate flags (optional)
	if (obj.flags !== undefined && obj.flags !== null) {
		if (!Array.isArray(obj.flags)) {
			errors.push({ field: 'flags', message: 'Flags must be an array' });
		} else {
			const invalidFlags = obj.flags.filter((flag: unknown) => !isPlantFlag(flag));
			if (invalidFlags.length > 0) {
				errors.push({
					field: 'flags',
					message: `Flags must be one of: ${Object.values(PlantFlag).join(', ')}`
				});
			}
		}
	}

	// Validate photoIds (optional)
	if (obj.photoIds !== undefined && obj.photoIds !== null) {
		if (!Array.isArray(obj.photoIds)) {
			errors.push({ field: 'photoIds', message: 'PhotoIds must be an array' });
		} else if (obj.photoIds.length > CONSTRAINTS.photoIds.maxItems) {
			errors.push({
				field: 'photoIds',
				message: `PhotoIds array must contain ${CONSTRAINTS.photoIds.maxItems} items or less`
			});
		} else {
			const invalidIds = obj.photoIds.filter((id: unknown) => {
				if (typeof id !== 'string') return true;
				const trimmed = id.trim();
				if (trimmed.length === 0) return true;
				// Allow long data URIs, enforce length only for non-data strings (e.g., filenames/URLs)
				const isDataUri = trimmed.startsWith('data:');
				if (!isDataUri && trimmed.length > CONSTRAINTS.photoIds.maxIdLength) return true;
				return false;
			});
			if (invalidIds.length > 0) {
				errors.push({
					field: 'photoIds',
					message: `Each photo ID must be a non-empty string; non-data IDs must be ${CONSTRAINTS.photoIds.maxIdLength} characters or less`
				});
			}
		}
	}

	return {
		valid: errors.length === 0,
		errors
	};
}

export function sanitizePlantInput(data: unknown): Partial<Plant> {
	if (!data || typeof data !== 'object') {
		return {};
	}

	const obj = data as Record<string, unknown>;

	return {
		...(typeof obj.species === 'string' && { species: obj.species.trim() }),
		...(typeof obj.name === 'string' && { name: obj.name.trim() }),
		...(typeof obj.sunLight === 'string' &&
			isSunlightRequirement(obj.sunLight) && { sunLight: obj.sunLight }),
		...(typeof obj.preferedTemperature === 'number' && {
			preferedTemperature: obj.preferedTemperature
		}),
		...(typeof obj.wateringIntervalDays === 'number' && {
			wateringIntervalDays: obj.wateringIntervalDays
		}),
		...(typeof obj.fertilizingIntervalDays === 'number' && {
			fertilizingIntervalDays: obj.fertilizingIntervalDays
		}),
		...(typeof obj.preferedHumidity === 'number' && { preferedHumidity: obj.preferedHumidity }),
		...(typeof obj.sprayIntervalDays === 'number' &&
			obj.sprayIntervalDays !== null && {
				sprayIntervalDays: obj.sprayIntervalDays
			}),
		...(typeof obj.lastWatered === 'string' && { lastWatered: obj.lastWatered }),
		...(typeof obj.lastFertilized === 'string' && { lastFertilized: obj.lastFertilized }),
		...(Array.isArray(obj.notes) && {
			notes: obj.notes
				.filter((note): note is string => typeof note === 'string')
				.map((note: string) => note.trim())
				.filter((note: string) => note.length > 0)
		}),
		...(Array.isArray(obj.flags) && {
			flags: obj.flags.filter((flag): flag is PlantFlag => isValidFlag(flag))
		}),
		...(Array.isArray(obj.photoIds) && {
			photoIds: obj.photoIds
				.filter((id): id is string => typeof id === 'string')
				.map((id: string) => id.trim())
				.filter((id: string) => id.length > 0)
		})
	};
}

function isValidFlag(flag: unknown): flag is PlantFlag {
	return Object.values(PlantFlag).includes(flag as PlantFlag);
}

function isPlantFlag(flag: unknown): flag is PlantFlag {
	return typeof flag === 'string' && Object.values(PlantFlag).includes(flag as PlantFlag);
}

function isSunlightRequirement(value: string): value is SunlightRequirement {
	return Object.values(SunlightRequirement).includes(value as SunlightRequirement);
}

function isValidISODate(dateString: string): boolean {
	if (typeof dateString !== 'string') return false;
	const date = new Date(dateString);
	return date instanceof Date && !isNaN(date.getTime()) && dateString === date.toISOString();
}
