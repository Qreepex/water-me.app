<script lang="ts">
	import type { Plant } from '$lib/types/api';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import { onMount, onDestroy } from 'svelte';
	import Can from '$lib/assets/Can.svg.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { formatPastTimestamp, formatFutureTimestamp } from '$lib/utils/timestamp.svelte';

	interface Props {
		plant: Plant;
		status: 'overdue' | 'due-soon' | 'ok';
		statusText: string;
		statusIcon: string;
		isWatering?: boolean;
		isSelected?: boolean;
		showNextWater?: boolean;
		nextWaterDate?: Date;
		onWater: (id: string) => void;
		onSelect: (id: string) => void;
	}

	const {
		plant,
		status,
		statusText,
		statusIcon,
		isWatering = false,
		isSelected = false,
		showNextWater = false,
		nextWaterDate,
		onWater,
		onSelect
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
		return `Last watered: ${formatPastTimestamp(new Date(lastWatered))}`;
	}

	function getNextWaterText(): string {
		if (!nextWaterDate) return 'No watering schedule';
		return `Next water: ${formatFutureTimestamp(nextWaterDate)}`;
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
	class="w-full overflow-hidden rounded-2xl border-2 bg-[var(--card-light)] shadow-md backdrop-blur transition hover:shadow-lg {`border-[var(--p-emerald)]/30 ${getStatusColor(
		status
	)}`}"
>
	<div class="flex items-center gap-4 p-4">
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
				<h3 class="text-lg font-bold text-[var(--text-light-main)]">{plant.name}</h3>
			</div>
			<p class="mb-2 text-sm text-[var(--text-light-main)]/60 italic">{plant.species}</p>

			<!-- Status and Info -->
			<div class="flex items-center gap-2">
				<span class="text-lg">{statusIcon}</span>
				<span class="text-xs font-medium text-[var(--text-light-main)]">{statusText}</span>
			</div>

			<!-- Additional Info -->
			{#if status !== 'overdue' && status !== 'due-soon'}
				{#if showNextWater}
					<p class="mt-1 text-xs text-[var(--text-light-main)]/70">{getNextWaterText()}</p>
				{:else}
					<p class="mt-1 text-xs text-[var(--text-light-main)]/70">{getLastWateredText()}</p>
				{/if}
			{/if}

			{#if plant.location?.room}
				<p class="mt-1 text-xs text-[var(--text-light-main)]/70">📍 {plant.location.room}</p>
			{/if}
		</div>
	</div>

	<!-- Two-Step Action Area -->
	{#if !isSelected}
		<!-- First Step: Select to Confirm -->
		<div class="p-3">
			<Button
				variant="water"
				size="md"
				onclick={() => onSelect(plant.id)}
				iconComponent={Can}
				text="plants.readyToWater"
				class="w-full"
			/>
		</div>
	{:else}
		<!-- Second Step: Confirm Action Buttons -->
		<div class="flex gap-2 p-3">
			<Button
				variant="danger"
				size="md"
				onclick={() => onSelect(plant.id)}
				text="plants.cancel"
				class="flex-1"
			/>
			<Button
				variant="water"
				size="md"
				onclick={() => onWater(plant.id)}
				disabled={isWatering}
				loading={isWatering}
				loadingText="plants.watering"
				iconComponent={Can}
				text={isWatering ? 'plants.watering' : 'plants.confirm'}
				class="flex-1"
			/>
		</div>
	{/if}
</div>
