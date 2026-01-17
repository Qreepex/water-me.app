import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getPlantById, updatePlant, deletePlant } from '$lib/db';
import { validatePlantInput, sanitizePlantInput } from '$lib/validation';

export const GET: RequestHandler = async ({ params }) => {
	try {
		if (!params.plantId || typeof params.plantId !== 'string') {
			return json({ error: 'Invalid plant ID' }, { status: 400 });
		}

		const plant = getPlantById(params.plantId);

		if (!plant) {
			return json({ error: 'Plant not found' }, { status: 404 });
		}

		return json(plant);
	} catch (error) {
		console.error('Error fetching plant:', error);
		return json({ error: 'Failed to fetch plant' }, { status: 500 });
	}
};

export const PUT: RequestHandler = async ({ params, request }) => {
	try {
		if (!params.plantId || typeof params.plantId !== 'string') {
			return json({ error: 'Invalid plant ID' }, { status: 400 });
		}

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

		const plant = updatePlant(params.plantId, sanitized);

		if (!plant) {
			return json({ error: 'Plant not found' }, { status: 404 });
		}

		return json(plant);
	} catch (error) {
		console.error('Error updating plant:', error);
		return json({ error: 'Failed to update plant' }, { status: 500 });
	}
};

export const PATCH: RequestHandler = async ({ params, request }) => {
	try {
		if (!params.plantId || typeof params.plantId !== 'string') {
			return json({ error: 'Invalid plant ID' }, { status: 400 });
		}

		let body: unknown;
		try {
			body = await request.json();
		} catch {
			return json({ error: 'Invalid JSON in request body' }, { status: 400 });
		}

		// For PATCH, we allow partial updates, so validate only provided fields
		const validation = validatePlantInput(body);
		if (!validation.valid) {
			return json({ error: 'Validation failed', details: validation.errors }, { status: 400 });
		}

		// Sanitize input
		const sanitized = sanitizePlantInput(body);

		const plant = updatePlant(params.plantId, sanitized);

		if (!plant) {
			return json({ error: 'Plant not found' }, { status: 404 });
		}

		return json(plant);
	} catch (error) {
		console.error('Error updating plant:', error);
		return json({ error: 'Failed to update plant' }, { status: 500 });
	}
};

export const DELETE: RequestHandler = async ({ params }) => {
	try {
		if (!params.plantId || typeof params.plantId !== 'string') {
			return json({ error: 'Invalid plant ID' }, { status: 400 });
		}

		const deleted = deletePlant(params.plantId);

		if (!deleted) {
			return json({ error: 'Plant not found' }, { status: 404 });
		}

		return json({ success: true });
	} catch (error) {
		console.error('Error deleting plant:', error);
		return json({ error: 'Failed to delete plant' }, { status: 500 });
	}
};
