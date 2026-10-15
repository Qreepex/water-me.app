import type {
	PlantFlag,
	SunlightRequirement,
	WateringMethod,
	WateringPlanMode,
	WaterType,
	FertilizerType,
	MoistureScale
} from './api';
import {
	SunlightRequirement as SR,
	WateringMethod as WM,
	WaterType as WT,
	FertilizerType as FT
} from './api';
import { WateringPlanMode as WPM, MoistureScale as MS } from './api';

export interface FormData {
	// Basic info
	name: string;
	species: string;
	isToxic: boolean;
	sunlight: SunlightRequirement;
	preferedTemperature: number;

	// Location
	room: string;
	position: string;
	isOutdoors: boolean;

	// Watering
	wateringMode: WateringPlanMode;
	wateringIntervalDays: number;
	wateringMethod: WateringMethod;
	waterType: WaterType;
	idealMoisture: number;
	moistureScale: MoistureScale;

	// Fertilizing
	fertilizingType: FertilizerType;
	fertilizingIntervalDays: number;
	npkRatio: string;
	concentrationPercent: number;
	activeInWinter: boolean;

	// Humidity
	targetHumidity: number;
	requiresMisting: boolean;
	mistingIntervalDays: number;
	requiresHumidifier: boolean;

	// Soil
	soilType: string;
	repottingCycle: number;
	soilComponents: string[];

	// Seasonality
	winterRestPeriod: boolean;
	winterWaterFactor: number;
	minTempCelsius: number;

	// Metadata
	flags: PlantFlag[];
	notes: string[];
}

export function createEmptyFormData(sunlightDefault?: SunlightRequirement): FormData {
	return {
		name: '',
		species: '',
		isToxic: false,
		sunlight: sunlightDefault || SR.Indirect_Sun,
		preferedTemperature: 20,
		room: '',
		position: '',
		isOutdoors: false,
		wateringMode: WPM.Manual,
		wateringIntervalDays: 7,
		wateringMethod: WM.Top,
		waterType: WT.Tap,
		idealMoisture: 5,
		moistureScale: MS.Scale1To10,
		fertilizingType: FT.Liquid,
		fertilizingIntervalDays: 30,
		npkRatio: '10:10:10',
		concentrationPercent: 50,
		activeInWinter: false,
		targetHumidity: 50,
		requiresMisting: false,
		mistingIntervalDays: 3,
		requiresHumidifier: false,
		soilType: 'Generic',
		repottingCycle: 2,
		soilComponents: [],
		winterRestPeriod: false,
		winterWaterFactor: 0.5,
		minTempCelsius: 15,
		flags: [],
		notes: []
	};
}
