<script lang="ts">
	import { tStore } from '$lib/i18n';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Alert from '$lib/components/ui/Message.svelte';
	import PlantList from '$lib/components/PlantList.svelte';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import type { SortOption } from '$lib/utils/plant';
	import PageContent from '$lib/components/layout/PageContent.svelte';

	const store = getPlantsStore();
	let sortBy = $state<SortOption>('nameAsc');
</script>

<div class="flex h-full flex-col overflow-hidden">
	<!-- Header -->
	<div class="flex-shrink-0">
		<PageHeader icon="🌱" title="common.app" description="common.appDescription" />
	</div>

	<PageContent>
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
			<!-- Scrollable Plant List -->
			<PlantList plants={store.plants} {sortBy} onSortChange={(value) => (sortBy = value)} />
		{/if}
	</PageContent>
</div>
