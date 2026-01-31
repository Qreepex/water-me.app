<script lang="ts">
	import { onMount } from 'svelte';
	import { tStore } from '$lib/i18n';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { sortPlants, daysAgo, getWateringStatus } from '$lib/utils/plant';
	import SortControls from '$lib/components/SortControls.svelte';
	import PlantCard from '$lib/components/PlantCard.svelte';
	import PageContainer from '$lib/components/layout/PageContainer.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Alert from '$lib/components/ui/Message.svelte';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import type { SortOption } from '$lib/utils/plant';
	import type { Plant } from '$lib/types/api';

	const store = getPlantsStore();
	let sortBy = $state<SortOption>('name');

	async function loadPlants() {
		store.setLoading(true);
		try {
			const result = await fetchData('/api/plants', {});
			if (!result.ok) {
				const errorMsg = result.error?.message || 'Failed to fetch plants';
				store.setError(errorMsg);
				return;
			}
			store.setPlants(result.data);
			// Prefetch first photo into cache to avoid second load on Android
			await prefetchPreviews(result.data);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : 'Unknown error';
			store.setError(errorMsg);
		} finally {
			store.setLoading(false);
		}
	}

	async function prefetchPreviews(items: Plant[]): Promise<void> {
		const { getImageObjectURL } = await import('$lib/utils/imageCache');
		for (const p of items) {
			const firstId = p.photoIds?.[0];
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const firstUrl = (p as any)?.photoUrls?.[0] as string | undefined;
			if (firstId && firstUrl) {
				// Fire and forget; this populates IndexedDB cache
				try {
					await getImageObjectURL(firstId, firstUrl);
				} catch {
					// ignore
				}
			}
		}
	}

	function getSortedPlants(): Plant[] {
		return sortPlants(store.plants, sortBy);
	}

	onMount(() => {
		loadPlants();
	});
</script>

<PageContainer>
	<!-- Header -->
	<PageHeader icon="🌱" title="common.app" description="common.appDescription" />

	<!-- Controls -->
	<div class="mb-8 flex items-center justify-between">
		<SortControls {sortBy} onSortChange={(value: SortOption) => (sortBy = value)} />
		<div class="font-medium text-[var(--text-light-main)]">
			{store.plants.length}{store.plants.length === 1
				? ' ' + $tStore('common.plant')
				: ' ' + $tStore('common.plants')}
		</div>
	</div>

	<!-- Loading State -->
	{#if store.loading}
		<LoadingSpinner message="Loading your plants..." icon="🌿" />
	{:else if store.error}
		<Alert type="error" title="Error loading plants" description={store.error} />
	{:else if store.plants.length === 0}
		<EmptyState
			icon="🪴"
			title={$tStore('plants.noPlants')}
			description={$tStore('plants.startAddingPlants')}
		/>
	{:else}
		<!-- Plant Grid -->
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each getSortedPlants() as plant (plant.id)}
				<PlantCard {plant} {daysAgo} {getWateringStatus} />
			{/each}
		</div>
	{/if}
</PageContainer>
