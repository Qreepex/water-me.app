import { fetchData } from '$lib/auth/fetch.svelte';
import { invalidateApiCache } from '$lib/utils/cache';
import { t } from '$lib/i18n';
import type { MoistureScale, MoistureReading } from '$lib/types/api';
import { Haptics, NotificationType } from '@capacitor/haptics';
import { getPlantsStore } from './plants.svelte';

interface WaterOptions {
	/** Optional moisture value measured right before watering, keyed by plant ID. */
	moistures?: Record<string, MoistureReading>;
}

/**
 * Watering session state for the water view / plant detail:
 * tracks which plants are currently being watered (in-flight), and exposes
 * helpers for logging soil moisture measurements.
 */
export function createWateringStore() {
	const plantsStore = getPlantsStore();
	let inFlightIds = $state<string[]>([]);

	function haptic(type: NotificationType): void {
		Haptics.notification({ type }).catch(() => {
			/* haptics unavailable (web) */
		});
	}

	function setInFlight(ids: string[], active: boolean): void {
		inFlightIds = active
			? [...inFlightIds, ...ids.filter((id) => !inFlightIds.includes(id))]
			: inFlightIds.filter((id) => !ids.includes(id));
	}

	async function invalidatePlant(id: string): Promise<void> {
		await invalidateApiCache(['/api/plants', `/api/plants/${id}`], {
			waitForAck: true,
			timeoutMs: 500
		});
		await plantsStore.reloadSilently();
	}

	async function water(ids: string[], options: WaterOptions = {}): Promise<boolean> {
		const targets = ids.filter((id) => !inFlightIds.includes(id));
		if (targets.length === 0) return false;
		setInFlight(targets, true);
		plantsStore.setError(null);

		try {
			const body: Record<string, unknown> = { plantIds: targets };
			if (options.moistures && Object.keys(options.moistures).length > 0) {
				body.moistures = options.moistures;
			}

			const response = await fetchData('/api/plants/water', {
				method: 'post',
				body: body as never
			});

			if (!response.ok) {
				plantsStore.setError(response.error?.message || t('plants.failedToWaterPlant'));
				haptic(NotificationType.Error);
				return false;
			}

			await invalidateApiCache(['/api/plants'], { waitForAck: true, timeoutMs: 500 });
			await plantsStore.reloadSilently();
			haptic(NotificationType.Success);
			return true;
		} catch (err) {
			plantsStore.setError(err instanceof Error ? err.message : t('plants.failedToWaterPlant'));
			haptic(NotificationType.Error);
			return false;
		} finally {
			setInFlight(targets, false);
		}
	}

	/**
	 * Record a soil moisture measurement. Resolves to the updated plant, or null
	 * on failure.
	 */
	async function logMoisture(
		id: string,
		value: number,
		options: { scale?: MoistureScale; date?: string; note?: string } = {}
	): Promise<boolean> {
		if (inFlightIds.includes(id)) return false;
		setInFlight([id], true);
		plantsStore.setError(null);

		try {
			const response = await fetchData('/api/plants/{id}/moisture', {
				method: 'post',
				params: { id },
				body: {
					value,
					...(options.scale ? { scale: options.scale } : {}),
					...(options.date ? { date: options.date } : {}),
					...(options.note ? { note: options.note } : {})
				}
			});

			if (!response.ok) {
				plantsStore.setError(response.error?.message || t('plants.failedToSaveMoisture'));
				haptic(NotificationType.Error);
				return false;
			}

			if (response.data) {
				plantsStore.setPlants(plantsStore.plants.map((p) => (p.id === id ? response.data! : p)));
			}
			await invalidateApiCache(['/api/plants', `/api/plants/${id}`], {
				waitForAck: true,
				timeoutMs: 300
			});
			haptic(NotificationType.Success);
			return true;
		} catch (err) {
			plantsStore.setError(err instanceof Error ? err.message : t('plants.failedToSaveMoisture'));
			haptic(NotificationType.Error);
			return false;
		} finally {
			setInFlight([id], false);
		}
	}

	/** Remove a single measurement from the history. */
	async function deleteMoisture(id: string, logId: string): Promise<boolean> {
		setInFlight([id], true);
		plantsStore.setError(null);

		try {
			const response = await fetchData('/api/plants/{id}/moisture/{logId}', {
				method: 'delete',
				params: { id, logId }
			});

			if (!response.ok) {
				plantsStore.setError(response.error?.message || t('plants.failedToDeleteMoisture'));
				return false;
			}

			if (response.data) {
				plantsStore.setPlants(plantsStore.plants.map((p) => (p.id === id ? response.data! : p)));
			}
			await invalidateApiCache(['/api/plants', `/api/plants/${id}`], {
				waitForAck: true,
				timeoutMs: 300
			});
			return true;
		} catch (err) {
			plantsStore.setError(err instanceof Error ? err.message : t('plants.failedToDeleteMoisture'));
			return false;
		} finally {
			setInFlight([id], false);
		}
	}

	/** Invalidate the cached plant and reload it (used after editor saves). */
	async function refreshPlant(id: string): Promise<void> {
		await invalidatePlant(id);
	}

	return {
		get inFlightIds() {
			return inFlightIds;
		},
		isWatering(id: string): boolean {
			return inFlightIds.includes(id);
		},
		get isWateringAny(): boolean {
			return inFlightIds.length > 0;
		},
		/** Water a single plant, optionally recording the measured moisture. */
		waterPlant: (id: string, options?: WaterOptions) => water([id], options),
		/** Water multiple plants at once (e.g. "water all due"). */
		waterPlants: (ids: string[], options?: WaterOptions) => water(ids, options),
		logMoisture,
		deleteMoisture,
		refreshPlant
	};
}

let wateringStore: ReturnType<typeof createWateringStore> | null = null;

export function getWateringStore() {
	if (!wateringStore) {
		wateringStore = createWateringStore();
	}
	return wateringStore;
}
