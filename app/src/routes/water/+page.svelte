<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { getPlantWaterStatus, getPlantStatusText, getStatusIcon } from '$lib/utils/plant';
	import { sortByWateringPriority } from '$lib/utils/watering';
	import WaterPlantCard from '$lib/components/WaterPlantCard.svelte';
	import PageContainer from '$lib/components/layout/PageContainer.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Alert from '$lib/components/ui/Message.svelte';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { onMount } from 'svelte';
	import { Haptics, NotificationType } from '@capacitor/haptics';
	import type { Plant } from '$lib/types/api';
	import { SvelteDate } from 'svelte/reactivity';
	import { tStore } from '$lib/i18n';

	const store = getPlantsStore();
	let selectedForWateringId = $state<string | null>(null);
	let wateringIds = $state<string[]>([]);
	let dismissedIds = $state<string[]>([]);

	async function loadPlants(): Promise<void> {
		store.setLoading(true);
		try {
			const response = await fetchData('/api/plants', {});
			if (!response.ok) {
				const errorMsg = response.error?.message || 'Failed to load plants';
				store.setError(errorMsg);
				return;
			}
			store.setPlants(response.data || []);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to load plants';
			store.setError(errorMsg);
		} finally {
			store.setLoading(false);
		}
	}

	function invalidateCache(): void {
		if (typeof navigator !== 'undefined' && 'serviceWorker' in navigator) {
			navigator.serviceWorker.controller?.postMessage({
				type: 'INVALIDATE_CACHE',
				urls: ['/api/plants']
			});
		}
	}

	function dismissPlant(id: string): void {
		if (!dismissedIds.includes(id)) {
			dismissedIds = [...dismissedIds, id];
		}
		selectedForWateringId = null;
	}

	function isWatering(id: string): boolean {
		return wateringIds.includes(id);
	}

	function toggleWateringSelection(id: string): void {
		if (selectedForWateringId === id) {
			selectedForWateringId = null;
		} else {
			selectedForWateringId = id;
		}
	}

	async function waterPlant(id: string): Promise<void> {
		if (wateringIds.includes(id)) return;
		wateringIds = [...wateringIds, id];
		store.setError(null);

		try {
			const response = await fetchData('/api/plants/water', {
				method: 'post',
				body: {
					plantIds: [id]
				}
			});

			if (!response.ok) {
				const errorMsg = response.error?.message || 'Failed to water plant';
				store.setError(errorMsg);
				try {
					await Haptics.notification({ type: NotificationType.Error });
				} catch {
					console.error('Haptics notification error');
				}
				return;
			}

			const now = new Date().toISOString();
			const updated = store.plants.map((p) => {
				if (p.id === id && p.watering) {
					return {
						...p,
						watering: {
							...p.watering,
							lastWatered: now
						}
					};
				}
				return p;
			});
			store.setPlants(updated);

			// Clear selection and invalidate cache
			selectedForWateringId = null;
			invalidateCache();

			try {
				await Haptics.notification({ type: NotificationType.Success });
			} catch {
				console.error('Haptics notification error');
			}
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to water plant';
			store.setError(errorMsg);
			try {
				await Haptics.notification({ type: NotificationType.Error });
			} catch {
				console.error('Haptics notification error');
			}
		} finally {
			wateringIds = wateringIds.filter((pid) => pid !== id);
		}
	}

	function getVisiblePlants() {
		const visible = sortByWateringPriority(store.plants).filter(
			(p) => !dismissedIds.includes(p.id)
		);
		// Sort by due status - due plants first, then others
		return visible.sort((a, b) => {
			const statusA = getPlantWaterStatus(a);
			const statusB = getPlantWaterStatus(b);
			const priorityMap = { overdue: 0, 'due-soon': 1, ok: 2 };
			return priorityMap[statusA] - priorityMap[statusB];
		});
	}

	function getDuePlants() {
		return getVisiblePlants().filter((p) => {
			const status = getPlantWaterStatus(p);
			return status === 'overdue' || status === 'due-soon';
		});
	}

	function getNotDuePlants() {
		return getVisiblePlants().filter((p) => {
			const status = getPlantWaterStatus(p);
			return status === 'ok';
		});
	}

	function getNextWaterDate(plant: Plant): Date {
		const lastWatered = plant.watering?.lastWatered
			? new Date(plant.watering.lastWatered)
			: new Date();
		const intervalDays = plant.watering?.intervalDays ?? 0;
		const nextWaterDate = new SvelteDate(lastWatered);
		nextWaterDate.setDate(nextWaterDate.getDate() + intervalDays);
		return nextWaterDate;
	}

	onMount(() => {
		loadPlants();
	});
</script>

<PageContainer gradient>
	<!-- Header -->
	<PageHeader icon="💧" title="menu.waterPlants" description="menu.wateringDescription" />

	<!-- Error Message -->
	{#if store.error}
		<Alert type="error" title="Error" description={store.error} />
	{/if}

	<!-- Loading State -->
	{#if store.loading}
		<LoadingSpinner message="Loading your plants..." icon="🌱" />
	{:else if store.plants.length === 0}
		<!-- Empty State -->
		<EmptyState icon="🪴" title="plants.noPlants" description="plants.startAddingPlants">
			<Button variant="primary" onclick={() => goto(resolve('/manage/create'))} text="addPlant" />
		</EmptyState>
	{:else if getVisiblePlants().length === 0}
		<EmptyState icon="✓" title="plants.allWatered" description="plants.allPlantsWatered" />
	{:else}
		<!-- Due Plants Section -->
		{#if getDuePlants().length > 0}
			<div class="mb-8">
				<div class="mb-4 flex items-center gap-2">
					<h2 class="text-xl font-bold text-[var(--text-light-main)]">
						🌵 {$tStore('plants.needsWater')}
					</h2>
					<span
						class="ml-auto rounded-full bg-[var(--status-error)] px-3 py-1 text-sm font-semibold text-white"
					>
						{getDuePlants().length}
					</span>
				</div>
				<div class="space-y-3">
					{#each getDuePlants() as plant (plant.id)}
						<WaterPlantCard
							{plant}
							status={getPlantWaterStatus(plant)}
							statusText={getPlantStatusText(plant)}
							statusIcon={getStatusIcon(getPlantWaterStatus(plant))}
							isWatering={isWatering(plant.id)}
							isSelected={selectedForWateringId === plant.id}
							onWater={waterPlant}
							onSelect={toggleWateringSelection}
							showNextWater={false}
						/>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Not Due Plants Section -->
		{#if getNotDuePlants().length > 0}
			<div>
				<div class="mb-4 flex items-center gap-2">
					<h2 class="text-xl font-bold text-[var(--text-light-main)]">
						✅ {$tStore('plants.watered')}
					</h2>
					<span
						class="ml-auto rounded-full bg-[var(--status-success)] px-3 py-1 text-sm font-semibold text-white"
					>
						{getNotDuePlants().length}
					</span>
				</div>
				<div class="space-y-3">
					{#each getNotDuePlants() as plant (plant.id)}
						<WaterPlantCard
							{plant}
							status={getPlantWaterStatus(plant)}
							statusText={getPlantStatusText(plant)}
							statusIcon={getStatusIcon(getPlantWaterStatus(plant))}
							isWatering={isWatering(plant.id)}
							isSelected={selectedForWateringId === plant.id}
							onWater={waterPlant}
							onSelect={toggleWateringSelection}
							onSkip={dismissPlant}
							showNextWater={true}
							nextWaterDate={getNextWaterDate(plant)}
						/>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</PageContainer>
