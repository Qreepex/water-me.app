import Database from 'better-sqlite3';
import path from 'path';
import { fileURLToPath } from 'url';
import type { Plant } from './types/types';
import { PlantFlag, SunlightRequirement } from './types/types';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const dbPath = path.join(__dirname, '..', '..', 'plants.db');
const db = new Database(dbPath);

// Enable foreign keys
db.pragma('foreign_keys = ON');

// Initialize database schema
export function initializeDatabase() {
	db.exec(`
		CREATE TABLE IF NOT EXISTS plants (
			id TEXT PRIMARY KEY,
			species TEXT NOT NULL,
			name TEXT NOT NULL,
			sunLight TEXT NOT NULL,
			preferedTemperature REAL NOT NULL,
			wateringIntervalDays INTEGER NOT NULL,
			lastWatered TEXT NOT NULL,
			fertilizingIntervalDays INTEGER NOT NULL,
			lastFertilized TEXT NOT NULL,
			preferedHumidity REAL NOT NULL,
			sprayIntervalDays INTEGER,
			notes TEXT NOT NULL DEFAULT '[]',
			flags TEXT NOT NULL DEFAULT '[]',
			photoIds TEXT NOT NULL DEFAULT '[]',
			createdAt TEXT NOT NULL,
			updatedAt TEXT NOT NULL
		);
	`);
}

// Initialize on import
initializeDatabase();

export function getDb() {
	return db;
}

export function getAllPlants(): Plant[] {
	const stmt = db.prepare('SELECT * FROM plants ORDER BY name ASC');
	const plants = stmt.all() as PlantRow[];
	return plants.map(parsePlantRow);
}

export function getPlantById(id: string): Plant | null {
	const stmt = db.prepare('SELECT * FROM plants WHERE id = ?');
	const plant = stmt.get(id) as PlantRow | undefined;
	return plant ? parsePlantRow(plant) : null;
}

export function createPlant(
	plant: Omit<Plant, 'lastWatered' | 'lastFertilized' | 'id'> & {
		lastWatered?: string;
		lastFertilized?: string;
		id?: string;
	}
): Plant {
	const id = plant.id || generateId();
	const now = new Date().toISOString();
	const lastWatered = plant.lastWatered || now;
	const lastFertilized = plant.lastFertilized || now;

	const stmt = db.prepare(`
		INSERT INTO plants (
			id, species, name, sunLight, preferedTemperature,
			wateringIntervalDays, lastWatered, fertilizingIntervalDays,
			lastFertilized, preferedHumidity, sprayIntervalDays,
			notes, flags, photoIds, createdAt, updatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`);

	stmt.run(
		id,
		plant.species,
		plant.name,
		plant.sunLight,
		plant.preferedTemperature,
		plant.wateringIntervalDays,
		lastWatered,
		plant.fertilizingIntervalDays,
		lastFertilized,
		plant.preferedHumidity,
		plant.sprayIntervalDays || null,
		JSON.stringify(plant.notes || []),
		JSON.stringify(plant.flags || []),
		JSON.stringify(plant.photoIds || []),
		now,
		now
	);

	return getPlantById(id)!;
}

export function updatePlant(id: string, updates: Partial<Plant>): Plant | null {
	const existing = getPlantById(id);
	if (!existing) return null;

	const now = new Date().toISOString();

	const stmt = db.prepare(`
		UPDATE plants SET
			species = ?,
			name = ?,
			sunLight = ?,
			preferedTemperature = ?,
			wateringIntervalDays = ?,
			lastWatered = ?,
			fertilizingIntervalDays = ?,
			lastFertilized = ?,
			preferedHumidity = ?,
			sprayIntervalDays = ?,
			notes = ?,
			flags = ?,
			photoIds = ?,
			updatedAt = ?
		WHERE id = ?
	`);

	stmt.run(
		updates.species ?? existing.species,
		updates.name ?? existing.name,
		updates.sunLight ?? existing.sunLight,
		updates.preferedTemperature ?? existing.preferedTemperature,
		updates.wateringIntervalDays ?? existing.wateringIntervalDays,
		updates.lastWatered ?? existing.lastWatered,
		updates.fertilizingIntervalDays ?? existing.fertilizingIntervalDays,
		updates.lastFertilized ?? existing.lastFertilized,
		updates.preferedHumidity ?? existing.preferedHumidity,
		updates.sprayIntervalDays ?? existing.sprayIntervalDays,
		JSON.stringify(updates.notes ?? existing.notes),
		JSON.stringify(updates.flags ?? existing.flags),
		JSON.stringify(updates.photoIds ?? existing.photoIds),
		now,
		id
	);

	return getPlantById(id);
}

export function deletePlant(id: string): boolean {
	const stmt = db.prepare('DELETE FROM plants WHERE id = ?');
	const result = stmt.run(id);
	return (result.changes ?? 0) > 0;
}

function parsePlantRow(row: PlantRow): Plant {
	let parsedNotes: string[] = [];
	let parsedFlags: PlantFlag[] = [];
	let parsedPhotoIds: string[] = [];

	try {
		const notesData = JSON.parse(row.notes || '[]');
		if (Array.isArray(notesData)) {
			parsedNotes = notesData.filter((item): item is string => typeof item === 'string');
		}
	} catch {
		// Handle JSON parse errors, keep empty array
	}

	try {
		const flagsData = JSON.parse(row.flags || '[]');
		if (Array.isArray(flagsData)) {
			parsedFlags = flagsData.filter(
				(item): item is PlantFlag =>
					typeof item === 'string' && Object.values(PlantFlag).includes(item as PlantFlag)
			);
		}
	} catch {
		// Handle JSON parse errors, keep empty array
	}

	try {
		const photoIdsData = JSON.parse(row.photoIds || '[]');
		if (Array.isArray(photoIdsData)) {
			parsedPhotoIds = photoIdsData.filter((item): item is string => typeof item === 'string');
		}
	} catch {
		// Handle JSON parse errors, keep empty array
	}

	const sunLight = isSunlightRequirementValue(row.sunLight)
		? row.sunLight
		: SunlightRequirement.INDIRECT_SUN;

	return {
		id: row.id,
		species: row.species,
		name: row.name,
		sunLight,
		preferedTemperature: row.preferedTemperature,
		wateringIntervalDays: row.wateringIntervalDays,
		lastWatered: row.lastWatered,
		fertilizingIntervalDays: row.fertilizingIntervalDays,
		lastFertilized: row.lastFertilized,
		preferedHumidity: row.preferedHumidity,
		sprayIntervalDays: row.sprayIntervalDays ?? undefined,
		notes: parsedNotes,
		flags: parsedFlags,
		photoIds: parsedPhotoIds
	};
}

interface PlantRow {
	id: string;
	species: string;
	name: string;
	sunLight: string;
	preferedTemperature: number;
	wateringIntervalDays: number;
	lastWatered: string;
	fertilizingIntervalDays: number;
	lastFertilized: string;
	preferedHumidity: number;
	sprayIntervalDays: number | null;
	notes: string;
	flags: string;
	photoIds: string;
	createdAt: string;
	updatedAt: string;
}

function isSunlightRequirementValue(value: string): value is SunlightRequirement {
	return Object.values(SunlightRequirement).includes(value as SunlightRequirement);
}

function generateId(): string {
	return `plant_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`;
}
