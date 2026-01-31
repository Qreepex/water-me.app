<script lang="ts">
	import type { Plant } from '$lib/types/api';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import { onMount, onDestroy } from 'svelte';
	import WaterDrop from '$lib/assets/WaterDrop.svg.svelte';
	import Can from '$lib/assets/Can.svg.svelte';
	import { daysAgo } from '$lib/utils/plant';

	interface Props {
		plant: Plant;
		status: 'overdue' | 'due-soon' | 'ok';
		statusText: string;
		statusIcon: string;
		isWatering?: boolean;
		isSelected?: boolean;
		onWater: (id: string) => void;
		onSelect: (id: string) => void;
		onSkip: (id: string) => void;
	}

	const {
		plant,
		status,
		statusText,
		statusIcon,
		isWatering = false,
		isSelected = false,
		onWater,
		onSelect,
		onSkip
	}: Props = $props();

	let previewUrl = $state<string | null>(null);

	function getStatusColor(status: 'overdue' | 'due-soon' | 'ok'): string {
		switch (status) {
			case 'overdue':
				return 'bg-[var(--status-error)]/20 border-[var(--status-error)]';
			case 'due-soon':
				return 'bg-[var(--status-warn)]/20 border-[var(--status-warn)]';
			default:
				return 'bg-[var(--status-success)]/20 border-[var(--status-success)]';
		}
	}

	function getLastWateredText(): string {
		const lastWatered = plant.watering?.lastWatered;
		if (!lastWatered) return 'Never watered';
		return `Last watered: ${daysAgo(lastWatered)}`;
	}

	onMount(async () => {
		const firstId = plant.photoIds?.[0];
		const firstUrl = plant?.photoUrls?.[0] as string | undefined;
		if (firstId && firstUrl) {
			previewUrl = await getImageObjectURL(firstId, firstUrl);
		}
	});

	onDestroy(() => {
		if (previewUrl) revokeObjectURL(previewUrl);
	});
</script>

<div
	class="w-full rounded-2xl border-2 bg-[var(--card-light)] p-4 shadow-md backdrop-blur {`border-[var(--p-emerald)]/30 ${getStatusColor(
		status
	)}`}"
>
	<div class="flex items-center gap-4">
		<!-- Photo -->
		{#if previewUrl}
			<img
				src={previewUrl}
				alt={plant.name}
				class="h-16 w-16 flex-shrink-0 rounded-lg object-cover"
			/>
		{:else}
			<div
				class="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-lg bg-[var(--p-emerald)]/30 text-2xl"
			>
				🌿
			</div>
		{/if}

		<!-- Plant Info -->
		<div class="flex-1 text-left">
			<div class="mb-1 flex items-center gap-2">
				<span class="text-xl">{statusIcon}</span>
				<h3 class="text-lg font-bold text-[var(--text-light-main)]">{plant.name}</h3>
			</div>
			<p class="mb-1 text-sm text-[var(--text-light-main)]/60 italic">{plant.species}</p>
			<p class="mb-2 text-sm font-medium text-[var(--text-light-main)]">{statusText}</p>

			<!-- Additional Info -->
			<div class="space-y-1 text-xs text-[var(--text-light-main)]/70">
				{#if status !== 'overdue' && status !== 'due-soon'}
					<p>{getLastWateredText()}</p>
				{/if}
				{#if plant.location?.room}
					<p>📍 {plant.location.room}</p>
				{/if}
				{#if plant.notes}
					<p>📝 {plant.notes}</p>
				{/if}
			</div>
		</div>
	</div>

	<!-- Two-Step Action Area -->
	{#if !isSelected}
		<!-- First Step: Select to Confirm -->
		<button
			class="mt-4 w-full rounded-xl border-2 border-[var(--p-emerald)]/40 bg-[var(--p-emerald)]/5 px-4 py-3 text-center font-semibold text-[var(--text-light-main)] transition hover:bg-[var(--p-emerald)]/10"
			onclick={() => onSelect(plant.id)}
			aria-label="Select to confirm watering"
		>
			<div class="flex items-center justify-center gap-2">
				<Can class="h-5 w-5" />
				<span>Ready to water?</span>
			</div>
		</button>
	{:else}
		<!-- Second Step: Confirm Action Buttons -->
		<div class="mt-4 grid grid-cols-2 gap-3">
			<button
				class="flex items-center justify-center gap-2 rounded-xl border border-[var(--p-emerald)]/40 bg-white px-3 py-3 text-sm font-semibold text-[var(--text-light-main)] transition hover:bg-[var(--bg-light)]"
				onclick={() => onSelect(plant.id)}
				aria-label="Cancel"
			>
				<span>Cancel</span>
			</button>
			<button
				class="flex items-center justify-center gap-2 rounded-xl bg-[var(--p-emerald)] px-3 py-3 text-sm font-semibold text-[var(--text-light-main)] transition hover:bg-[var(--p-emerald-dark)] disabled:opacity-60"
				onclick={() => onWater(plant.id)}
				disabled={isWatering}
				aria-label="Confirm watering"
			>
				<Can class="h-5 w-5" />
				<span>{isWatering ? 'Watering…' : 'Confirm'}</span>
			</button>
		</div>
	{/if}
</div>
