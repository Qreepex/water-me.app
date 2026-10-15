<script lang="ts">
	import HoldButton from '$lib/components/ui/HoldButton.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import MoodFace from '$lib/components/ui/MoodFace.svelte';
	import { toneBg, toneText } from '$lib/components/ui/tone';
	import { tStore } from '$lib/i18n';
	import { getWateringStore } from '$lib/stores/watering.svelte';
	import type { Plant } from '$lib/types/api';
	import {
		getPlantStatusText,
		getPlantWaterStatus,
		getStatusVisual,
		getWateringProgress
	} from '$lib/utils/plant';
	import Timestamp from '../ui/Timestamp.svelte';
	import PlantImage from './PlantImage.svelte';

	interface Props {
		plant: Plant;
		/** Compact styling for already-watered plants */
		compact?: boolean;
		/** Stagger index for the entrance animation */
		index?: number;
	}

	const { plant, compact = false, index = 0 }: Props = $props();
	const watering = getWateringStore();

	const status = $derived(getPlantWaterStatus(plant));
	const visual = $derived(getStatusVisual(status));
	const statusText = $derived(getPlantStatusText(plant));
	const progress = $derived(getWateringProgress(plant));
	const photoUrls = $derived((plant as { photoUrls?: string[] }).photoUrls ?? []);

	const moodFace = $derived(
		status === 'overdue' ? 'thirsty' : status === 'due-soon' ? 'soon' : 'happy'
	);

	const statusRing = $derived(
		status === 'overdue' ? 'ring-danger/60' : status === 'due-soon' ? 'ring-warn/70' : 'ring-ok/40'
	);

	const barClass = $derived(
		status === 'overdue' ? 'bg-danger' : status === 'due-soon' ? 'bg-warn' : 'bg-brand'
	);

	const isWatering = $derived(watering.isWatering(plant.id));
</script>

<div
	class="pop-in w-full overflow-hidden rounded-3xl border border-ink/5 bg-surface shadow-md transition-shadow hover:shadow-lg"
	style:animation-delay="{Math.min(index * 40, 280)}ms"
>
	<!-- Hero: full-bleed photo with name overlay + mood badge -->
	<div class="relative {compact ? 'h-20' : 'h-28'}">
		<PlantImage
			photoId={plant.photoIds?.[0]}
			remoteUrl={photoUrls[0]}
			alt={plant.name}
			class="h-full w-full object-cover"
			fallbackClass={compact ? '[&_svg]:h-7 [&_svg]:w-7' : ''}
		/>
		<div
			class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent"
		></div>

		<!-- Mood badge -->
		<div
			class="absolute top-2 right-2 rounded-full bg-surface/95 p-1 shadow-md ring-2 {statusRing}"
		>
			<MoodFace mood={moodFace} size={compact ? 34 : 42} class={toneText[visual.tone]} />
		</div>

		<div class="absolute inset-x-3 bottom-2 flex min-w-0 items-baseline gap-2">
			<h3 class="truncate text-lg font-bold text-white drop-shadow">{plant.name}</h3>
			{#if plant.species}
				<p class="truncate text-xs text-white/80 italic">{plant.species}</p>
			{/if}
		</div>
	</div>

	<!-- Status + thirst meter -->
	<div class="p-3.5">
		<div class="flex items-center justify-between gap-2">
			<span class="inline-flex items-center gap-1.5 text-sm font-semibold {toneText[visual.tone]}">
				<Icon name={visual.icon} size={16} />
				{$tStore(statusText.key, statusText.args)}
			</span>
			{#if plant.location?.room}
				<span
					class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium {toneBg[
						'neutral'
					]} text-ink-soft"
				>
					<Icon name="map-pin" size={12} />
					{plant.location.room}
				</span>
			{/if}
		</div>

		<div class="mt-2.5 flex items-center gap-2">
			<div class="h-2 flex-1 overflow-hidden rounded-full bg-ink/10" role="presentation">
				<div
					class="h-full rounded-full transition-all duration-500 {barClass}"
					style:width="{progress * 100}%"
				></div>
			</div>
			<span class="shrink-0 text-[11px] text-ink-soft">
				{#if plant.watering?.lastWatered}
					<Timestamp date={new Date(plant.watering.lastWatered)} />
				{:else}
					—
				{/if}
			</span>
		</div>
	</div>

	<!-- Big, thumb-friendly hold action -->
	<div class="px-3.5 pb-3.5">
		<HoldButton
			variant={compact ? 'secondary' : 'water'}
			size={compact ? 'md' : 'lg'}
			icon="droplet"
			label={isWatering ? $tStore('plants.watering') : $tStore('plants.holdToWater')}
			disabled={isWatering}
			onconfirm={() => watering.waterPlant(plant.id)}
		/>
	</div>
</div>
