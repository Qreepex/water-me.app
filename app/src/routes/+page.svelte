<script lang="ts">
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Alert from '$lib/components/ui/Message.svelte';
	import PlantList from '$lib/components/PlantList.svelte';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import type { SortOption } from '$lib/utils/plant';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	const store = getPlantsStore();
	let sortBy = $state<SortOption>('nameAsc');
</script>

<!-- Header -->
<div class="flex-shrink-0">
	<PageHeader icon="🌱" title="common.app" description="common.appDescription" />
</div>

<PageContent>
	<!-- Loading State -->
	{#if store.loading}
		<LoadingSpinner message="common.loadingPlants" icon="🌿" />
	{:else if store.error}
		<Alert type="error" title="common.errorLoadingPlants" description={store.error} />
	{:else if store.plants.length === 0}
			<EmptyState icon="🪴" title="plants.noPlants" description="plants.startAddingPlants">
			<Button variant="primary" onclick={() => goto(resolve('/manage/new'))} text="plants.addPlant" />
		</EmptyState>
	{:else}
		<!-- Scrollable Plant List -->
		<PlantList plants={store.plants} {sortBy} onSortChange={(value) => (sortBy = value)} />
	{/if}
</PageContent>
