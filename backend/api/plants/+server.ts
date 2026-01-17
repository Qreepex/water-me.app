import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getAllPlants, createPlant } from '$lib/db';
import { validatePlantInput, sanitizePlantInput } from '$lib/validation';

export const GET: RequestHandler = async () => {
	try {
		const plants = getAllPlants();
		return json(plants);
	} catch (error) {
		console.error('Error fetching plants:', error);
		return json({ error: 'Failed to fetch plants' }, { status: 500 });
	}
};

export const POST: RequestHandler = async ({ request }) => {
	try {
		let body: unknown;
		try {
			body = await request.json();
		} catch {
			return json({ error: 'Invalid JSON in request body' }, { status: 400 });
		}

		// Validate input
		const validation = validatePlantInput(body);
		if (!validation.valid) {
			return json({ error: 'Validation failed', details: validation.errors }, { status: 400 });
		}

		// Sanitize input
		const sanitized = sanitizePlantInput(body);

		// Type guard - all required fields should be present after sanitization
		if (
			!sanitized.species ||
			!sanitized.name ||
			!sanitized.sunLight ||
			sanitized.preferedTemperature === undefined ||
			sanitized.wateringIntervalDays === undefined ||
			sanitized.fertilizingIntervalDays === undefined ||
			sanitized.preferedHumidity === undefined
		) {
			return json({ error: 'Missing required fields after sanitization' }, { status: 400 });
		}

		const bodyRecord = body as Record<string, unknown>;
		const plantId: string | undefined = isString(bodyRecord.id) ? bodyRecord.id : undefined;
		const plant = createPlant({
			id: plantId,
			species: sanitized.species,
			name: sanitized.name,
			sunLight: sanitized.sunLight,
			preferedTemperature: sanitized.preferedTemperature,
			wateringIntervalDays: sanitized.wateringIntervalDays,
			lastWatered: sanitized.lastWatered,
			fertilizingIntervalDays: sanitized.fertilizingIntervalDays,
			lastFertilized: sanitized.lastFertilized,
			preferedHumidity: sanitized.preferedHumidity,
			sprayIntervalDays: sanitized.sprayIntervalDays,
			notes: sanitized.notes || [],
			flags: sanitized.flags || [],
			photoIds: sanitized.photoIds || []
		});

		return json(plant, { status: 201 });
	} catch (error) {
		console.error('Error creating plant:', error);
		return json({ error: 'Failed to create plant' }, { status: 500 });
	}
};

function isString(value: unknown): value is string {
	return typeof value === 'string';
}
