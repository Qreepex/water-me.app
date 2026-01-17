<script lang="ts">
	import { onMount } from 'svelte';
	import type { Plant } from '$lib/types/types';

	type SortOption =
		| 'name'
		| 'lastWatered'
		| 'lastFertilized'
		| 'sprayInterval'
		| 'wateringInterval';

	let plants: Plant[] = [];
	let loading = true;
	let error: string | null = null;
	let sortBy: SortOption = 'name';

	onMount(async () => {
		try {
			const response = await fetch('/api/plants');
			if (!response.ok) throw new Error('Failed to fetch plants');
			plants = await response.json();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			loading = false;
		}
	});

	function getSortedPlants(): Plant[] {
		const sorted = [...plants];

		switch (sortBy) {
			case 'name':
				return sorted.sort((a, b) => a.name.localeCompare(b.name));
			case 'lastWatered':
				return sorted.sort(
					(a, b) => new Date(b.lastWatered).getTime() - new Date(a.lastWatered).getTime()
				);
			case 'lastFertilized':
				return sorted.sort(
					(a, b) => new Date(b.lastFertilized).getTime() - new Date(a.lastFertilized).getTime()
				);
			case 'wateringInterval':
				return sorted.sort((a, b) => a.wateringIntervalDays - b.wateringIntervalDays);
			case 'sprayInterval':
				return sorted.sort((a, b) => (a.sprayIntervalDays || 999) - (b.sprayIntervalDays || 999));
			default:
				return sorted;
		}
	}

	function daysAgo(dateString: string): string {
		const days = Math.floor((Date.now() - new Date(dateString).getTime()) / (1000 * 60 * 60 * 24));
		if (days === 0) return 'Today';
		if (days === 1) return 'Yesterday';
		return `${days} days ago`;
	}

	function getWateringStatus(plant: Plant): { text: string; color: string } {
		const days = Math.floor(
			(Date.now() - new Date(plant.lastWatered).getTime()) / (1000 * 60 * 60 * 24)
		);
		const daysUntilWater = plant.wateringIntervalDays - days;

		if (daysUntilWater <= 0) return { text: '🌵 Needs water!', color: 'text-red-600' };
		if (daysUntilWater <= 1) return { text: '⚠️ Water soon', color: 'text-yellow-600' };
		return { text: `✓ In ${daysUntilWater} days`, color: 'text-green-600' };
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-green-50 via-emerald-50 to-teal-50 p-8">
	<div class="mx-auto max-w-7xl">
		<!-- Header -->
		<div class="mb-12">
			<h1 class="mb-2 flex items-center gap-3 text-5xl font-bold text-green-800">🌱 My Plants</h1>
			<p class="text-lg text-green-700">Take care of your green friends</p>
		</div>

		<!-- Controls -->
		<div class="mb-8 flex items-center justify-between">
			<div class="flex items-center gap-3">
				<label for="sort" class="font-semibold text-green-800">Sort by:</label>
				<select
					id="sort"
					bind:value={sortBy}
					class="rounded-lg border-2 border-green-300 bg-white px-4 py-2 font-medium text-green-800 transition hover:border-green-400 focus:border-green-500 focus:outline-none"
				>
					<option value="name">Plant Name</option>
					<option value="lastWatered">Last Watered</option>
					<option value="lastFertilized">Last Fertilized</option>
					<option value="wateringInterval">Watering Frequency</option>
					<option value="sprayInterval">Spray Frequency</option>
				</select>
			</div>
			<div class="font-medium text-green-800">
				{plants.length}
				{plants.length === 1 ? 'plant' : 'plants'}
			</div>
		</div>

		<!-- Loading & Error States -->
		{#if loading}
			<div class="flex min-h-96 items-center justify-center">
				<div class="text-center">
					<div class="mb-4 animate-bounce text-6xl">🌿</div>
					<p class="text-lg font-medium text-green-700">Loading your plants...</p>
				</div>
			</div>
		{:else if error}
			<div class="rounded-lg border-2 border-red-400 bg-red-100 px-6 py-4 text-red-800">
				<p class="font-bold">Error loading plants</p>
				<p>{error}</p>
			</div>
		{:else if plants.length === 0}
			<div class="py-16 text-center">
				<div class="mb-4 text-8xl">🪴</div>
				<p class="text-xl font-medium text-green-800">No plants yet!</p>
				<p class="mt-2 text-green-700">Start adding your plants to track their care.</p>
			</div>
		{:else}
			<!-- Plant Grid -->
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
				{#each getSortedPlants() as plant (plant.id)}
					<div
						class="group overflow-hidden rounded-2xl bg-white shadow-md transition-all duration-300 hover:scale-105 hover:shadow-xl"
					>
						<!-- Image -->
						<div
							class="relative flex h-48 items-center justify-center overflow-hidden bg-gradient-to-br from-green-200 to-emerald-300"
						>
							{#if plant.photoIds.length > 0}
								<img
									src={plant.photoIds[0]}
									alt={plant.name}
									class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
								/>
							{:else}
								<div class="text-7xl transition-transform duration-300 group-hover:scale-110">
									🌱
								</div>
							{/if}
						</div>

						<!-- Content -->
						<div class="p-5">
							<!-- Name and Species -->
							<h3 class="mb-1 line-clamp-2 text-xl font-bold text-green-800">{plant.name}</h3>
							<p class="mb-4 line-clamp-1 text-sm text-green-600">{plant.species}</p>

							<!-- Watering Status -->
							<div class="mb-4">
								<div class={`mb-2 text-sm font-semibold ${getWateringStatus(plant).color}`}>
									{getWateringStatus(plant).text}
								</div>
								<p class="text-xs text-gray-600">
									Watered {daysAgo(plant.lastWatered)}
								</p>
							</div>

							<!-- Metadata Grid -->
							<div class="mb-4 grid grid-cols-2 gap-3 text-xs">
								<div class="rounded-lg bg-blue-50 p-2">
									<div class="font-semibold text-blue-600">💧</div>
									<p class="mt-1 text-xs text-gray-700">Every {plant.wateringIntervalDays}d</p>
								</div>
								<div class="rounded-lg bg-yellow-50 p-2">
									<div class="font-semibold text-yellow-600">🥗</div>
									<p class="mt-1 text-xs text-gray-700">Every {plant.fertilizingIntervalDays}d</p>
								</div>
								<div class="rounded-lg bg-purple-50 p-2">
									<div class="font-semibold text-purple-600">☀️</div>
									<p class="mt-1 text-xs text-gray-700">
										{plant.sunLight.split(' ').slice(0, 1).join('')}
									</p>
								</div>
								<div class="rounded-lg bg-teal-50 p-2">
									<div class="font-semibold text-teal-600">💨</div>
									<p class="mt-1 text-xs text-gray-700">{plant.preferedHumidity}%</p>
								</div>
							</div>

							<!-- Spray Info -->
							{#if plant.sprayIntervalDays}
								<div class="mb-3 rounded-lg bg-cyan-50 p-2">
									<p class="text-xs text-gray-600">
										💦 Spray every <span class="font-semibold text-cyan-700"
											>{plant.sprayIntervalDays}</span
										> days
									</p>
									<p class="mt-1 text-xs text-gray-600">
										Last: <span class="font-semibold">{daysAgo(plant.lastFertilized)}</span>
									</p>
								</div>
							{/if}

							<!-- Flags -->
							{#if plant.flags.length > 0}
								<div class="mb-3 flex flex-wrap gap-2">
									{#each plant.flags as flag}
										<span
											class="rounded-full bg-orange-100 px-2 py-1 text-xs font-medium text-orange-800"
										>
											⚡ {flag}
										</span>
									{/each}
								</div>
							{/if}

							<!-- Notes Preview -->
							{#if plant.notes.length > 0}
								<div class="border-t border-gray-200 pt-3">
									<p class="line-clamp-2 text-xs text-gray-600">
										📝 {plant.notes[0]}
									</p>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
</style>
