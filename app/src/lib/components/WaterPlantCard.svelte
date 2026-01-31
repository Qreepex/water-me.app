<script lang="ts">
	import type { Plant } from '$lib/types/api';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import { onMount, onDestroy } from 'svelte';

	interface Props {
		plant: Plant;
		isSelected: boolean;
		status: 'overdue' | 'due-soon' | 'ok';
		statusText: string;
		statusIcon: string;
		onToggle: (id: string) => void;
	}

	const { plant, isSelected, status, statusText, statusIcon, onToggle }: Props = $props();

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

<button
	onclick={() => onToggle(plant.id)}
	class="w-full rounded-2xl border-2 bg-[var(--card-light)] p-4 shadow-md backdrop-blur transition {isSelected
		? 'border-[var(--p-emerald)] bg-[var(--p-emerald)]/10'
		: `border-[var(--p-emerald)]/30 ${getStatusColor(status)}`}"
>
	<div class="flex items-center gap-4">
		<!-- Checkbox -->
		<div
			class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg border-2 transition {isSelected
				? 'border-[var(--p-emerald)] bg-[var(--p-emerald)] text-white'
				: 'border-[var(--p-emerald)] bg-[var(--card-light)] text-[var(--p-emerald)]'}"
		>
			{#if isSelected}
				<span class="text-xl">✓</span>
			{:else}
				<span class="text-xl">○</span>
			{/if}
		</div>

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
			<p class="text-sm font-medium text-[var(--text-light-main)]">{statusText}</p>
		</div>
	</div>
</button>
