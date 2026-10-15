import { fetchData } from '$lib/auth/fetch.svelte';
import { t } from '$lib/i18n';
import type { Plant } from '$lib/types/api';
import { invalidateApiCache } from '$lib/utils/cache';
import { imageCacheStore } from './imageCache.svelte';

interface PlantsStore {
	readonly plants: Plant[];
	readonly loading: boolean;
	readonly error: string | null;
	readonly hasLoaded: boolean;

	setPlants(plants: Plant[]): void;
	setLoading(loading: boolean): void;
	setError(error: string | null): void;
	/** Fetch plants from the API and prefetch their images. Sets loading state. */
	loadPlants(): Promise<void>;
	/** Like loadPlants but keeps the current list visible (no loading spinner). */
	reloadSilently(): Promise<void>;
	/**
	 * Drop cached API responses and reload silently. Used when the app resumes:
	 * the service worker serves GET /api/* cache-first for 15 minutes, so without
	 * invalidating first a refresh would just return the stale response and
	 * changes made on another device would stay invisible.
	 */
	syncFromServer(): Promise<void>;
	/** Invalidate API + image caches and reload everything. */
	refresh(): Promise<void>;
	reset(): void;
}

export function createPlantsStore(): PlantsStore {
	let plants = $state<Plant[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let hasLoaded = $state(false);

	async function fetchPlants(): Promise<void> {
		const result = await fetchData('/api/plants', {});
		if (!result.ok) {
			error = result.error?.message || t('plants.failedToFetchPlants');
			return;
		}
		plants = result.data || [];
		error = null;
		await prefetchPlantImages(plants);
	}

	return {
		get plants() {
			return plants;
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get hasLoaded() {
			return hasLoaded;
		},
		setPlants(newPlants: Plant[]) {
			plants = newPlants;
		},
		setLoading(newLoading: boolean) {
			loading = newLoading;
		},
		setError(newError: string | null) {
			error = newError;
		},
		async loadPlants() {
			loading = true;
			error = null;
			try {
				await fetchPlants();
			} catch (err) {
				error = err instanceof Error ? err.message : t('plants.failedToFetchPlants');
			} finally {
				loading = false;
				hasLoaded = true;
			}
		},
		async reloadSilently() {
			try {
				await fetchPlants();
			} catch {
				// Keep showing stale data on silent reload failures.
			}
		},
		async syncFromServer() {
			await invalidateApiCache(['/api/plants'], { waitForAck: true, timeoutMs: 500 });
			await this.reloadSilently();
		},
		async refresh() {
			imageCacheStore.cleanup(0);
			await invalidateApiCache(['/api/plants'], { waitForAck: true, timeoutMs: 500 });
			await this.loadPlants();
		},
		reset() {
			plants = [];
			loading = true;
			error = null;
			hasLoaded = false;
		}
	};
}

async function prefetchPlantImages(plants: Plant[]): Promise<void> {
	const prefetches: Promise<unknown>[] = [];

	for (const plant of plants) {
		const photoIds = plant.photoIds ?? [];
		const photoUrls = (plant as { photoUrls?: string[] }).photoUrls ?? [];

		for (let i = 0; i < photoIds.length; i++) {
			const photoId = photoIds[i];
			const photoUrl = photoUrls[i];
			if (!photoId || !photoUrl) continue;
			prefetches.push(imageCacheStore.getImageURL(photoId, photoUrl).catch(() => null));
		}
	}

	await Promise.allSettled(prefetches);
}

// Global plants store singleton
let plantsStore: PlantsStore | null = null;

export function getPlantsStore(): PlantsStore {
	if (!plantsStore) {
		plantsStore = createPlantsStore();
	}
	return plantsStore;
}
