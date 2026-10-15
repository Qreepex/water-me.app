import type { IconName } from '$lib/components/ui/icons';
import type { Plant } from '$lib/types/api';
import {
	FertilizerType,
	MoistureScale,
	SunlightRequirement,
	WateringMethod,
	WateringPlanMode,
	WaterType
} from '$lib/types/api';
import type { FormData } from '$lib/types/forms';
import { createEmptyFormData } from '$lib/types/forms';

export interface EditorSection {
	key: string;
	icon: IconName;
	titleKey: string;
}

/** Display order + identity of all editor sections (drives scroll-spy and section nav). */
export const EDITOR_SECTIONS: EditorSection[] = [
	{ key: 'basic', icon: 'leaf', titleKey: 'plants.basicInformation' },
	{ key: 'photos', icon: 'camera', titleKey: 'plants.photos' },
	{ key: 'location', icon: 'map-pin', titleKey: 'plants.location' },
	{ key: 'watering', icon: 'droplet', titleKey: 'plants.wateringTitle' },
	{ key: 'fertilizing', icon: 'flask-conical', titleKey: 'plants.fertilizingTitle' },
	{ key: 'humidity', icon: 'cloud-drizzle', titleKey: 'plants.humidityTitle' },
	{ key: 'soil', icon: 'layers', titleKey: 'plants.soilTitle' },
	{ key: 'seasonality', icon: 'snowflake', titleKey: 'plants.seasonalityTitle' },
	{ key: 'metadata', icon: 'tag', titleKey: 'plants.metadata' }
];

/** Create form data pre-filled from an existing plant (edit flows). */
export function initializeFormData(plant: Plant): FormData {
	return {
		name: plant.name,
		species: plant.species,
		isToxic: plant.isToxic,
		sunlight: plant.sunlight as SunlightRequirement,
		preferedTemperature: plant.preferedTemperature,
		room: plant.location?.room ?? '',
		position: plant.location?.position ?? '',
		isOutdoors: plant.location?.isOutdoors ?? false,
		wateringMode: plant.watering?.mode ?? WateringPlanMode.Manual,
		wateringIntervalDays: plant.watering?.intervalDays ?? 7,
		wateringMethod: plant.watering?.method ?? WateringMethod.Top,
		waterType: plant.watering?.waterType ?? WaterType.Tap,
		idealMoisture: plant.watering?.idealMoisture ?? 5,
		moistureScale: plant.watering?.moistureScale ?? MoistureScale.Scale1To10,
		fertilizingType: plant.fertilizing?.type ?? FertilizerType.Liquid,
		fertilizingIntervalDays: plant.fertilizing?.intervalDays ?? 30,
		npkRatio: plant.fertilizing?.npkRatio ?? '10:10:10',
		concentrationPercent: plant.fertilizing?.concentrationPercent ?? 50,
		activeInWinter: plant.fertilizing?.activeInWinter ?? false,
		targetHumidity: plant.humidity?.targetHumidityPct ?? 50,
		requiresMisting: plant.humidity?.requiresMisting ?? false,
		mistingIntervalDays: plant.humidity?.mistingIntervalDays ?? 3,
		requiresHumidifier: plant.humidity?.requiresHumidifier ?? false,
		soilType: plant.soil?.type ?? 'Generic',
		repottingCycle: plant.soil?.repottingCycle ?? 2,
		soilComponents: plant.soil?.components ?? [],
		winterRestPeriod: plant.seasonality?.winterRestPeriod ?? false,
		winterWaterFactor: plant.seasonality?.winterWaterFactor ?? 0.5,
		minTempCelsius: plant.seasonality?.minTempCelsius ?? 15,
		flags: plant.flags ?? [],
		notes: plant.notes ?? []
	};
}

/**
 * Full PATCH payload covering every section at once (single save for the
 * unified editor). Intentionally omits history timestamps (lastWatered etc.)
 * so saving settings never touches care history.
 */
export function buildFullPayload(formData: FormData): Record<string, unknown> {
	return {
		name: formData.name.trim(),
		species: formData.species.trim(),
		sunlight: formData.sunlight,
		isToxic: formData.isToxic,
		preferedTemperature: formData.preferedTemperature,
		location: {
			room: formData.room.trim(),
			position: formData.position.trim(),
			isOutdoors: formData.isOutdoors
		},
		watering: {
			mode: formData.wateringMode,
			intervalDays:
				formData.wateringMode === 'Moisture' ? undefined : formData.wateringIntervalDays,
			method: formData.wateringMethod,
			waterType: formData.waterType,
			idealMoisture: formData.wateringMode === 'Moisture' ? formData.idealMoisture : undefined,
			moistureScale: formData.wateringMode === 'Moisture' ? formData.moistureScale : undefined
		},
		fertilizing: {
			type: formData.fertilizingType,
			intervalDays: formData.fertilizingIntervalDays,
			npkRatio: formData.npkRatio.trim(),
			concentrationPercent: formData.concentrationPercent,
			activeInWinter: formData.activeInWinter
		},
		humidity: {
			targetHumidityPct: formData.targetHumidity,
			requiresMisting: formData.requiresMisting,
			mistingIntervalDays: formData.mistingIntervalDays,
			requiresHumidifier: formData.requiresHumidifier
		},
		soil: {
			type: formData.soilType.trim(),
			repottingCycle: formData.repottingCycle,
			components: formData.soilComponents
		},
		seasonality: {
			winterRestPeriod: formData.winterRestPeriod,
			winterWaterFactor: formData.winterWaterFactor,
			minTempCelsius: formData.minTempCelsius
		},
		flags: formData.flags,
		notes: formData.notes
	};
}

function hasChanged<K extends keyof FormData>(
	formData: FormData,
	initial: FormData,
	key: K
): boolean {
	if (Array.isArray(formData[key]) && Array.isArray(initial[key])) {
		return JSON.stringify(formData[key]) !== JSON.stringify(initial[key]);
	}
	return formData[key] !== initial[key];
}

/**
 * Create payload containing only fields that differ from the empty defaults
 * (POST /api/plants). `initial` defaults to a fresh empty form.
 */
export function buildCreatePayload(
	formData: FormData,
	initial: FormData = createEmptyFormData()
): Record<string, unknown> {
	const now = new Date().toISOString();
	const payload: Record<string, unknown> = {
		name: formData.name.trim(),
		species: formData.species.trim()
	};

	if (hasChanged(formData, initial, 'isToxic')) payload.isToxic = formData.isToxic;
	if (hasChanged(formData, initial, 'sunlight')) payload.sunlight = formData.sunlight;
	if (hasChanged(formData, initial, 'preferedTemperature')) {
		payload.preferedTemperature = formData.preferedTemperature;
	}

	if (
		hasChanged(formData, initial, 'room') ||
		hasChanged(formData, initial, 'position') ||
		hasChanged(formData, initial, 'isOutdoors')
	) {
		payload.location = {
			room: formData.room.trim(),
			position: formData.position.trim(),
			isOutdoors: formData.isOutdoors
		};
	}

	if (
		hasChanged(formData, initial, 'wateringMode') ||
		hasChanged(formData, initial, 'wateringIntervalDays') ||
		hasChanged(formData, initial, 'wateringMethod') ||
		hasChanged(formData, initial, 'waterType') ||
		hasChanged(formData, initial, 'idealMoisture') ||
		hasChanged(formData, initial, 'moistureScale')
	) {
		payload.watering = {
			mode: formData.wateringMode,
			intervalDays:
				formData.wateringMode === 'Moisture' ? undefined : formData.wateringIntervalDays,
			method: formData.wateringMethod,
			waterType: formData.waterType,
			idealMoisture: formData.wateringMode === 'Moisture' ? formData.idealMoisture : undefined,
			moistureScale: formData.wateringMode === 'Moisture' ? formData.moistureScale : undefined,
			lastWatered: now
		};
	}

	if (
		hasChanged(formData, initial, 'fertilizingType') ||
		hasChanged(formData, initial, 'fertilizingIntervalDays') ||
		hasChanged(formData, initial, 'npkRatio') ||
		hasChanged(formData, initial, 'concentrationPercent') ||
		hasChanged(formData, initial, 'activeInWinter')
	) {
		payload.fertilizing = {
			type: formData.fertilizingType,
			intervalDays: formData.fertilizingIntervalDays,
			npkRatio: formData.npkRatio.trim(),
			concentrationPercent: formData.concentrationPercent,
			activeInWinter: formData.activeInWinter,
			lastFertilized: now
		};
	}

	if (
		hasChanged(formData, initial, 'targetHumidity') ||
		hasChanged(formData, initial, 'requiresMisting') ||
		hasChanged(formData, initial, 'mistingIntervalDays') ||
		hasChanged(formData, initial, 'requiresHumidifier')
	) {
		payload.humidity = {
			targetHumidityPct: formData.targetHumidity,
			requiresMisting: formData.requiresMisting,
			mistingIntervalDays: formData.mistingIntervalDays,
			requiresHumidifier: formData.requiresHumidifier
		};
	}

	if (
		hasChanged(formData, initial, 'soilType') ||
		hasChanged(formData, initial, 'repottingCycle') ||
		hasChanged(formData, initial, 'soilComponents')
	) {
		payload.soil = {
			type: formData.soilType.trim(),
			repottingCycle: formData.repottingCycle,
			components: formData.soilComponents
		};
	}

	if (
		hasChanged(formData, initial, 'winterRestPeriod') ||
		hasChanged(formData, initial, 'winterWaterFactor') ||
		hasChanged(formData, initial, 'minTempCelsius')
	) {
		payload.seasonality = {
			winterRestPeriod: formData.winterRestPeriod,
			winterWaterFactor: formData.winterWaterFactor,
			minTempCelsius: formData.minTempCelsius
		};
	}

	if (hasChanged(formData, initial, 'flags')) payload.flags = formData.flags;
	if (hasChanged(formData, initial, 'notes')) payload.notes = formData.notes;

	return payload;
}

export function isFormDirty(current: FormData, original: FormData): boolean {
	return JSON.stringify(current) !== JSON.stringify(original);
}

/** Name + species are mandatory when creating or renaming a plant. */
export function isBasicInfoValid(formData: FormData): boolean {
	return formData.name.trim().length > 0 && formData.species.trim().length > 0;
}
