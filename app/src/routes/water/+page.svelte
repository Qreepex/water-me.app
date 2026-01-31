<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import {
		getDaysUntilWater,
		getPlantWaterStatus,
		getPlantStatusText,
		getStatusIcon
	} from '$lib/utils/plant';
	import { sortByWateringPriority } from '$lib/utils/watering';
	import WaterPlantCard from '$lib/components/WaterPlantCard.svelte';
	import PageContainer from '$lib/components/layout/PageContainer.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { onMount } from 'svelte';

	const store = getPlantsStore();
	let watering = $state(false);
	let selectedIds = $state<string[]>([]);

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

	function togglePlant(id: string): void {
		if (selectedIds.includes(id)) {
			selectedIds = selectedIds.filter((sid) => sid !== id);
		} else {
			selectedIds = [...selectedIds, id];
		}
	}

	function selectAll(): void {
		selectedIds = store.plants.map((p) => p.id);
	}

	function selectDueToday(): void {
		const due = store.plants.filter((p) => getDaysUntilWater(p) <= 0);
		selectedIds = due.map((p) => p.id);
	}

	function clearSelection(): void {
		selectedIds = [];
	}

	async function waterSelectedPlants(): Promise<void> {
		if (selectedIds.length === 0) return;

		watering = true;
		store.setError(null);

		try {
			const response = await fetchData('/api/plants/water', {
				method: 'post',
				body: {
					plantIds: selectedIds
				}
			});

			if (!response.ok) {
				const errorMsg = response.error?.message || 'Failed to water plants';
				store.setError(errorMsg);
				return;
			}

			// Update the plants locally
			const now = new Date().toISOString();
			const updated = store.plants.map((p) => {
				if (selectedIds.includes(p.id) && p.watering) {
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

			// Clear selection
			selectedIds = [];
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Failed to water plants';
			store.setError(errorMsg);
		} finally {
			watering = false;
		}
	}

	onMount(() => {
		loadPlants();
	});
</script>

<PageContainer gradient>
	<!-- Header -->
	<PageHeader icon="💧" title="Water Plants" description="Quick watering view" />

	<!-- Error Message -->
	{#if store.error}
		<Alert type="error" title="Error" description={store.error} />
	{/if}

	<!-- Loading State -->
	{#if store.loading}
		<LoadingSpinner message="Loading your plants..." icon="🌱" />
	{:else if store.plants.length === 0}
		<!-- Empty State -->
		<EmptyState
			icon="🪴"
			title="No plants yet"
			description="Add plants to start tracking watering schedules"
		>
			<Button variant="primary" onclick={() => goto(resolve('/manage/create'))} text="addPlant" />
		</EmptyState>
	{:else}
		<!-- Quick Actions Bar -->
		<div
			class="mb-6 flex flex-wrap gap-2 rounded-lg border border-[var(--p-emerald)]/30 bg-[var(--card-light)] p-4 shadow-md backdrop-blur"
		>
			<Button
				variant="secondary"
				size="sm"
				onclick={selectDueToday}
				class="flex-1"
				text="selectOverdue"
			/>
			<Button variant="primary" size="sm" onclick={selectAll} class="flex-1" text="selectAll" />
			<Button variant="ghost" size="sm" onclick={clearSelection} class="flex-1" text="clear" />
		</div>

		<!-- Plant List -->
		<div class="space-y-4">
			{#each sortByWateringPriority(store.plants) as plant (plant.id)}
				<WaterPlantCard
					{plant}
					isSelected={selectedIds.includes(plant.id)}
					status={getPlantWaterStatus(plant)}
					statusText={getPlantStatusText(plant)}
					statusIcon={getStatusIcon(getPlantWaterStatus(plant))}
					onToggle={togglePlant}
				/>
			{/each}
		</div>
	{/if}
</PageContainer>

<!-- Bottom Action Bar (Fixed) -->
{#if selectedIds.length > 0}
	<div
		class="fixed right-0 bottom-0 left-0 border-t border-[var(--p-emerald)] bg-[var(--card-light)]/95 px-4 py-4 shadow-xl backdrop-blur"
	>
		<div class="mx-auto flex max-w-6xl items-center justify-between gap-4">
			<div class="text-sm font-medium text-[var(--text-light-main)]">
				{selectedIds.length}
				{selectedIds.length === 1 ? 'plant' : 'plants'} selected
			</div>
			<Button
				variant="primary"
				onclick={waterSelectedPlants}
				disabled={watering}
				text={watering ? '💧 Watering...' : '💧 Water Selected'}
			/>
		</div>
	</div>
{/if}
