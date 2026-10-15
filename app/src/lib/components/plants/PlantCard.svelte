<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { tStore } from '$lib/i18n';
	import type { Plant } from '$lib/types/api';
	import { formatMoisture, getPlan, getWateringMode, getWateringStatus } from '$lib/utils/plant';
	import Timestamp from '../ui/Timestamp.svelte';
	import Chip from '../ui/Chip.svelte';
	import Icon from '../ui/Icon.svelte';
	import { toneBg, toneText } from '../ui/tone';
	import type { IconName } from '../ui/icons';
	import PlantImage from './PlantImage.svelte';

	interface Props {
		plant: Plant;
		/** Stagger index for the entrance animation */
		index?: number;
	}

	const { plant, index = 0 }: Props = $props();

	const photoUrl = $derived((plant as { photoUrls?: string[] }).photoUrls?.[0]);
	const wateringStatus = $derived(getWateringStatus(plant));

	interface MetaTile {
		icon: IconName;
		text: string;
		tone: keyof typeof toneBg;
	}

	const metaTiles = $derived.by((): MetaTile[] => {
		const tiles: MetaTile[] = [];
		if (getWateringMode(plant) === 'Moisture') {
			const ideal = plant.watering?.idealMoisture;
			if (ideal !== undefined && ideal !== null) {
				tiles.push({
					icon: 'ruler',
					text: `${$tStore('plants.idealMoisture')}: ${formatMoisture(ideal, plant.watering?.moistureScale)}`,
					tone: 'brand'
				});
			}
		} else if (plant.watering?.intervalDays) {
			tiles.push({
				icon: 'droplet',
				text: `${$tStore('plants.every')} ${plant.watering.intervalDays} ${$tStore('plants.days')}`,
				tone: 'info'
			});
		}
		if (plant.fertilizing?.intervalDays) {
			tiles.push({
				icon: 'flask-conical',
				text: `${$tStore('plants.every')} ${plant.fertilizing.intervalDays} ${$tStore('plants.days')}`,
				tone: 'warn'
			});
		}
		if (plant.sunlight) {
			tiles.push({
				icon: 'sun',
				text: $tStore('plants.sunlight.' + plant.sunlight),
				tone: 'brand'
			});
		}
		if (plant.humidity?.targetHumidityPct) {
			tiles.push({
				icon: 'cloud-drizzle',
				text: `${plant.humidity.targetHumidityPct}%`,
				tone: 'info'
			});
		}
		return tiles;
	});

	function openPlant() {
		goto(resolve(`/plant/${plant.id}`));
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			openPlant();
		}
	}
</script>

<div
	role="button"
	tabindex="0"
	onclick={openPlant}
	onkeydown={onKeydown}
	class="group pop-in cursor-pointer touch-manipulation overflow-hidden rounded-2xl bg-surface shadow transition-all duration-200 active:scale-[0.97] active:shadow-sm"
	style:animation-delay="{Math.min(index * 30, 240)}ms"
>
	<!-- Image -->
	<div class="relative h-40 overflow-hidden">
		<PlantImage
			photoId={plant.photoIds?.[0]}
			remoteUrl={photoUrl}
			alt={plant.name}
			class="h-full w-full object-cover transition-transform duration-300 group-active:scale-105"
		/>
	</div>

	<!-- Content -->
	<div class="p-3.5">
		<h3 class="line-clamp-1 text-base font-bold text-ink">{plant.name}</h3>
		<p class="mb-2 line-clamp-1 text-xs text-ok dark:text-[#2fd077]">{plant.species}</p>

		<!-- Watering Status -->
		{#if wateringStatus}
			<div class="mb-3">
				<div
					class="mb-1 flex items-center gap-1.5 text-xs font-semibold {toneText[
						wateringStatus.tone
					]}"
				>
					<Icon name={wateringStatus.icon} size={14} />
					<span>{$tStore(wateringStatus.text, wateringStatus.args)}</span>
				</div>
				{#if plant.watering?.lastWatered}
					<p class="text-xs text-ink-soft">
						{$tStore('plants.lastWatered')}:
						<Timestamp date={new Date(plant.watering.lastWatered)} />
					</p>
				{/if}
			</div>
		{/if}

		<!-- Moisture readout for moisture-based plants -->
		{#if getWateringMode(plant) === 'Moisture'}
			{@const plan = getPlan(plant)}
			{#if plan?.currentMoisture !== undefined && plan?.currentMoisture !== null}
				<p class="mb-3 text-xs text-ink-soft">
					{$tStore('plants.currentMoisture')}:
					<span class="font-semibold text-ink">
						{formatMoisture(plan.currentMoisture, plant.watering?.moistureScale)}
					</span>
				</p>
			{/if}
		{/if}

		<!-- Metadata Grid -->
		{#if metaTiles.length > 0}
			<div class="grid grid-cols-2 gap-2">
				{#each metaTiles as tile (tile.icon + tile.text)}
					<div class="rounded-lg p-1.5 text-xs {toneBg[tile.tone]}">
						<Icon name={tile.icon} size={13} />
						<span class="ml-1 text-ink/80">{tile.text}</span>
					</div>
				{/each}
			</div>
		{/if}

		<!-- Misting Info -->
		{#if plant.humidity?.requiresMisting && plant.humidity?.mistingIntervalDays}
			<div class="mt-2 rounded-lg bg-info/15 p-2 text-xs text-ink/80">
				<span class="inline-flex items-center gap-1"
					><Icon name="cloud-drizzle" size={13} /> {$tStore('plants.sprayEvery')}</span
				>
				<span class="font-semibold text-info">{plant.humidity.mistingIntervalDays}</span>
				{$tStore('plants.days')}
				{#if plant.humidity.lastMisted}
					· {$tStore('plants.lastSprayedStatus')}:
					<span class="font-semibold">
						<Timestamp date={new Date(plant.humidity.lastMisted)} />
					</span>
				{/if}
			</div>
		{/if}

		<!-- Flags -->
		{#if plant.flags && plant.flags.length > 0}
			<div class="mt-2 flex flex-wrap gap-1.5">
				{#each plant.flags as flag (flag)}
					<Chip tone="warn" icon="zap" text={flag} />
				{/each}
			</div>
		{/if}

		<!-- Notes Preview -->
		{#if plant.notes && plant.notes.length > 0}
			<div class="mt-2 border-t border-brand/15 pt-2">
				<p class="line-clamp-2 text-xs text-ink-soft">
					<span class="inline-flex items-center gap-1.5"
						><Icon name="file-text" size={13} /> {plant.notes[0]}</span
					>
				</p>
			</div>
		{/if}
	</div>
</div>
