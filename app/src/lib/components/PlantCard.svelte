<script lang="ts">
	import { tStore } from '$lib/i18n';
	import type { Plant } from '$lib/types/api';
	import { imageCacheStore } from '$lib/stores/imageCache.svelte';
	import { onDestroy } from 'svelte';
	import { goto } from '$app/navigation';

	interface Props {
		plant: Plant;
		daysAgo: (dateString: string) => string;
		getWateringStatus: (plant: Plant) => { text: string; color: string };
	}

	const { plant, daysAgo, getWateringStatus }: Props = $props();

	const firstId = plant.photoIds?.[0];
	// Get the URL from cache once (already preloaded in Auth)
	const previewUrl = $state(firstId ? imageCacheStore.getImageURLSync(firstId) : null);

	function openPlant() {
		goto(`/plant/${plant.id}`);
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			openPlant();
		}
	}

	onDestroy(() => {
		if (firstId) {
			imageCacheStore.releaseImage(firstId);
		}
	});
</script>

<div
	role="button"
	tabindex="0"
	onclick={openPlant}
	onkeydown={onKeydown}
	class="group cursor-pointer overflow-hidden rounded-2xl border border-[var(--p-emerald)]/30 bg-[var(--card-light)] shadow-md transition-all duration-300 hover:border-[var(--p-emerald)]/60 hover:bg-[var(--card-light)]/80 hover:shadow-xl"
>
	<!-- Image -->
	<div
		class="relative flex h-48 items-center justify-center overflow-hidden rounded-t-2xl bg-gradient-to-br from-[var(--p-emerald)] to-[var(--p-emerald-dark)]"
	>
		{#if previewUrl}
			<img
				src={previewUrl}
				alt={plant.name}
				class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
			/>
		{:else}
			<div class="text-7xl transition-transform duration-300 group-hover:scale-110">🌱</div>
		{/if}
	</div>

	<!-- Content -->
	<div class="p-5">
		<!-- Name and Species -->
		<h3 class="mb-1 line-clamp-2 text-xl font-bold text-[var(--text-light-main)]">{plant.name}</h3>
		<p class="mb-4 line-clamp-1 text-sm text-[var(--status-success)]">{plant.species}</p>

		<!-- Watering Status -->
		<div class="mb-4">
			<div class={`mb-2 text-sm font-semibold ${getWateringStatus(plant).color}`}>
				{getWateringStatus(plant).text}
			</div>
			<p class="text-xs text-[var(--text-light-main)]/60">
				Watered {daysAgo(plant.watering?.lastWatered ?? '')}
			</p>
		</div>

		<!-- Metadata Grid -->
		<div class="grid grid-cols-2 gap-3 text-xs">
			<div class="rounded-lg bg-[var(--p-emerald)]/20 p-2">
				<div class="font-semibold text-[var(--p-emerald-dark)]">💧</div>
				<p class="mt-1 text-xs text-[var(--text-light-main)]/80">
					Every {plant.watering?.intervalDays}d
				</p>
			</div>
			{#if plant.fertilizing?.intervalDays}
				<div class="rounded-lg bg-[var(--status-warn)]/20 p-2">
					<div class="font-semibold text-[var(--status-warn)]">🥗</div>
					<p class="mt-1 text-xs text-[var(--text-light-main)]/80">
						Every {plant.fertilizing?.intervalDays}d
					</p>
				</div>
			{/if}
			{#if plant.sunlight}
				<div class="rounded-lg bg-[var(--status-info)]/20 p-2">
					<div class="font-semibold text-[var(--status-info)]">☀️</div>
					<p class="mt-1 text-xs text-[var(--text-light-main)]/80">
						{$tStore('plants.sunlight.' + plant.sunlight)}
					</p>
				</div>
			{/if}
			{#if plant.humidity}
				<div class="rounded-lg bg-[var(--p-emerald)]/20 p-2">
					<div class="font-semibold text-[var(--p-emerald-dark)]">💨</div>
					<p class="mt-1 text-xs text-[var(--text-light-main)]/80">
						{plant.humidity?.targetHumidityPct}%
					</p>
				</div>
			{/if}
		</div>

		<!-- Spray Info -->
		{#if plant.humidity?.requiresMisting}
			<div class="mb-3 rounded-lg bg-[var(--status-info)]/20 p-2">
				<p class="text-xs text-[var(--text-light-main)]/70">
					💦 Spray every <span class="font-semibold text-[var(--status-info)]"
						>{plant.humidity?.mistingIntervalDays}</span
					> days
				</p>
				<p class="mt-1 text-xs text-[var(--text-light-main)]/70">
					Last: <span class="font-semibold">{daysAgo(plant.fertilizing?.lastFertilized ?? '')}</span
					>
				</p>
			</div>
		{/if}

		<!-- Flags -->
		{#if plant.flags && plant.flags.length > 0}
			<div class="mb-3 flex flex-wrap gap-2">
				{#each plant.flags as flag (flag)}
					<span
						class="rounded-full bg-[var(--status-warn)]/30 px-2 py-1 text-xs font-medium text-[var(--status-warn)]"
						>⚡ {flag}</span
					>
				{/each}
			</div>
		{/if}

		<!-- Notes Preview -->
		{#if plant.notes && plant.notes.length > 0}
			<div class="border-t border-[var(--p-emerald)]/30 pt-3">
				<p class="line-clamp-2 text-xs text-[var(--text-light-main)]/60">📝 {plant.notes[0]}</p>
			</div>
		{/if}
	</div>
</div>
