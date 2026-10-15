<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Scrollable from '$lib/components/layout/Scrollable.svelte';
	import PlantImage from '$lib/components/plants/PlantImage.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import HoldButton from '$lib/components/ui/HoldButton.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Chip from '$lib/components/ui/Chip.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import MoodFace from '$lib/components/ui/MoodFace.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import type { IconName } from '$lib/components/ui/icons';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { getWateringStore } from '$lib/stores/watering.svelte';
	import type { Plant } from '$lib/types/api';
	import {
		formatMoisture,
		getPlan,
		getPlantStatusText,
		getPlantWaterStatus,
		getWateringMode,
		getWateringProgress,
		hasWateringConfig as plantHasWateringConfig
	} from '$lib/utils/plant';
	import Timestamp from '$lib/components/ui/Timestamp.svelte';
	import { SvelteDate } from 'svelte/reactivity';
	import MoistureLogSheet from '$lib/components/plants/MoistureLogSheet.svelte';

	const store = getPlantsStore();
	const watering = getWateringStore();

	let measuring = $state(false);

	const plantId = $derived(page.params.plant ?? '');
	const plant = $derived(store.plants.find((p) => p.id === plantId || p.slug === plantId) ?? null);
	const photoUrls = $derived(
		(plant as unknown as { photoUrls?: string[] } | null)?.photoUrls ?? []
	);
	const hasWateringConfig = $derived(plant ? plantHasWateringConfig(plant) : false);
	const isMoisturePlan = $derived(plant ? getWateringMode(plant) === 'Moisture' : false);
	const plan = $derived(plant ? getPlan(plant) : null);

	// ---- Playful Tamagotchi-style mood derived from the watering schedule ----
	const mood = $derived.by(
		(): {
			face: 'happy' | 'soon' | 'thirsty';
			key: string;
			color: string;
			animate: boolean;
		} | null => {
			if (!plant || !hasWateringConfig) return null;
			const status = getPlantWaterStatus(plant);
			if (status === 'overdue') {
				return { face: 'thirsty', key: 'plants.moodThirsty', color: 'text-danger', animate: true };
			}
			if (status === 'due-soon') {
				return { face: 'soon', key: 'plants.moodDueSoon', color: 'text-warn', animate: false };
			}
			return { face: 'happy', key: 'plants.moodHappy', color: 'text-ok', animate: false };
		}
	);

	const statusText = $derived(plant ? getPlantStatusText(plant) : null);
	const progressPct = $derived(plant ? getWateringProgress(plant) * 100 : 0);

	const barColor = $derived(
		!plant || getPlantWaterStatus(plant) === 'overdue'
			? 'bg-danger'
			: getPlantWaterStatus(plant) === 'due-soon'
				? 'bg-warn'
				: 'bg-brand'
	);

	const nextWaterDate = $derived.by(() => {
		if (plan?.wateringDueAt) return new SvelteDate(plan.wateringDueAt);
		if (!plant?.watering?.lastWatered || !plant.watering.intervalDays) return null;
		const next = new SvelteDate(plant.watering.lastWatered);
		next.setDate(next.getDate() + plant.watering.intervalDays);
		return next;
	});

	// ---- Care overview: only fields that are actually set ----
	interface CareRow {
		icon: IconName;
		labelKey: string;
		value: string;
	}

	const careRows = $derived.by((): CareRow[] => {
		if (!plant) return [];
		const rows: CareRow[] = [];
		const every = (days?: number) =>
			days ? `${$tStore('plants.every')} ${days} ${$tStore('plants.days')}` : '';
		// Missing enum keys surface as the raw key (consistent with the rest of
		// the i18n layer) instead of silently falling back to the raw value.
		const option = (prefix: string, value?: string) => (value ? $tStore(`${prefix}.${value}`) : '');

		if (plant.sunlight) {
			rows.push({
				icon: 'sun',
				labelKey: 'plants.formSunlightRequirements',
				value: option('plants.sunlight', plant.sunlight)
			});
		}
		if (plant.watering?.intervalDays) {
			rows.push({
				icon: 'droplet',
				labelKey: 'plants.wateringTitle',
				value: [
					every(plant.watering.intervalDays),
					option('plants.wateringMethodOptions', plant.watering.method),
					option('plants.waterTypeOptions', plant.watering.waterType)
				]
					.filter(Boolean)
					.join(' · ')
			});
		}
		if (plant.fertilizing?.intervalDays) {
			rows.push({
				icon: 'flask-conical',
				labelKey: 'plants.fertilizingTitle',
				value: [
					every(plant.fertilizing.intervalDays),
					option('plants.fertilizerTypeOptions', plant.fertilizing.type),
					plant.fertilizing.npkRatio
				]
					.filter(Boolean)
					.join(' · ')
			});
		}
		if (plant.humidity?.targetHumidityPct) {
			rows.push({
				icon: 'cloud-drizzle',
				labelKey: 'plants.humidityTitle',
				value: [
					`${plant.humidity.targetHumidityPct}%`,
					plant.humidity.requiresMisting && plant.humidity.mistingIntervalDays
						? `${$tStore('plants.sprayEvery')} ${plant.humidity.mistingIntervalDays} ${$tStore('plants.days')}`
						: ''
				]
					.filter(Boolean)
					.join(' · ')
			});
		}
		if (plant.soil?.type) {
			rows.push({
				icon: 'layers',
				labelKey: 'plants.soilTitle',
				value: [
					plant.soil.type,
					plant.soil.repottingCycle
						? `${plant.soil.repottingCycle} ${$tStore('plants.years')}`
						: '',
					plant.soil.components?.length ? plant.soil.components.join(', ') : ''
				]
					.filter(Boolean)
					.join(' · ')
			});
		}
		if (plant.seasonality) {
			rows.push({
				icon: 'snowflake',
				labelKey: 'plants.seasonalityTitle',
				value: [
					plant.seasonality.winterRestPeriod ? $tStore('plants.winterRest') : '',
					plant.seasonality.minTempCelsius !== undefined
						? `≥ ${plant.seasonality.minTempCelsius}°C`
						: ''
				]
					.filter(Boolean)
					.join(' · ')
			});
		}
		return rows;
	});

	function editPlant(current: Plant) {
		goto(resolve(`/manage/${current.id}/edit`));
	}

	function formatDate(date?: string | null): string {
		return date ? new Date(date).toLocaleDateString() : '-';
	}
</script>

<PageHeader
	icon="leaf"
	title={plant?.name ?? 'plants.myPlants'}
	description={plant?.species}
	backHref={resolve('/')}
>
	{#if plant}
		<Button
			variant="primary"
			size="sm"
			text="common.edit"
			icon="pencil"
			onclick={() => editPlant(plant)}
		/>
	{/if}
</PageHeader>

<PageContent>
	{#if store.loading}
		<Spinner message="common.loadingPlantDetails" />
	{:else if !plant}
		<Alert type="error" title="common.error" description="plants.notFound" />
	{:else}
		{@const p = plant}
		<Scrollable noPadding>
			<div class="space-y-4 px-2 pb-8">
				<!-- Hero -->
				<Card>
					<div class="relative h-52">
						<PlantImage
							photoId={p.photoIds?.[0]}
							remoteUrl={photoUrls[0]}
							alt={p.name}
							class="h-full w-full object-cover"
							fallbackClass="text-7xl"
						/>
						<div
							class="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-black/65 to-transparent"
						></div>
						<div class="absolute inset-x-0 bottom-0 p-4">
							<h2 class="truncate text-2xl font-bold text-white drop-shadow">{p.name}</h2>
							{#if p.species}
								<p class="truncate text-sm text-white/85 italic">{p.species}</p>
							{/if}
						</div>
					</div>
				</Card>

				<!-- Chips: only what is actually set, no noise -->
				{#if p.location?.room || p.isToxic || p.flags?.length}
					<div class="flex flex-wrap gap-2 px-1">
						{#if p.location?.room}
							<Chip
								tone="brand"
								icon="map-pin"
								text={p.location.position
									? `${p.location.room} · ${p.location.position}`
									: p.location.room}
							/>
						{/if}
						{#if p.isToxic}
							<Chip tone="danger" icon="skull" text={$tStore('plants.toxic')} />
						{/if}
						{#each p.flags ?? [] as flag (flag)}
							<Chip tone="warn" icon="zap" text={flag} />
						{/each}
					</div>
				{/if}

				<!-- Mood & watering status -->
				{#if mood}
					<Card padded>
						<div class="flex items-center gap-4">
							<div
								class="shrink-0 {mood.color} {mood.animate ? 'animate-bounce' : ''}"
								role="img"
								aria-label="Plant mood"
							>
								<MoodFace mood={mood.face} size={64} />
							</div>
							<div class="min-w-0 flex-1">
								<p class="text-lg leading-tight font-bold text-ink">{$tStore(mood.key)}</p>
								{#if statusText}
									<p class="text-sm text-ink-soft">
										{$tStore(statusText.key, statusText.args)}
									</p>
								{/if}
							</div>
						</div>

						<!-- Thirst-o-meter -->
						<div
							class="mt-4 h-2.5 overflow-hidden rounded-full bg-ink/10"
							role="progressbar"
							aria-valuemin={0}
							aria-valuemax={100}
							aria-valuenow={Math.round(progressPct)}
						>
							<div
								class="h-full rounded-full transition-all duration-500 {barColor}"
								style:width="{progressPct}%"
							></div>
						</div>
						<div class="mt-2 flex items-center justify-between gap-2 text-xs text-ink-soft">
							<span class="inline-flex items-center gap-1.5">
								<Icon name="clock" size={13} />
								{#if p.watering?.lastWatered}
									<Timestamp date={new SvelteDate(p.watering.lastWatered)} />
								{:else}
									—
								{/if}
							</span>
							{#if nextWaterDate}
								<span class="inline-flex items-center gap-1.5">
									<Icon name="arrow-right" size={13} />
									<Timestamp date={nextWaterDate} mode="future" />
								</span>
							{/if}
						</div>

						<div class="mt-4">
							<HoldButton
								variant="water"
								size="lg"
								icon="droplet"
								label={watering.isWatering(p.id)
									? $tStore('plants.watering')
									: $tStore('plants.holdToWater')}
								disabled={watering.isWatering(p.id)}
								onconfirm={() => watering.waterPlant(p.id)}
							/>
						</div>
					</Card>
				{:else}
					<!-- No schedule yet: helpful nudge instead of empty stats -->
					<Card padded>
						<div class="flex items-center gap-4">
							<div class="shrink-0 text-brand-dark dark:text-brand">
								<MoodFace mood="soon" size={56} />
							</div>
							<div class="min-w-0 flex-1">
								<p class="text-sm text-ink-soft">{$tStore('plants.setupScheduleHint')}</p>
							</div>
						</div>
						<Button
							variant="primary"
							class="mt-3 w-full"
							icon="droplet"
							text="plants.configureWateringSchedule"
							onclick={() => goto(resolve(`/manage/${p.id}/edit#watering`))}
						/>
					</Card>
				{/if}

				<!-- Moisture plan: current value, prediction, quick actions -->
				{#if isMoisturePlan}
					<Card padded>
						<div class="flex items-start justify-between gap-3">
							<div>
								<p class="text-sm font-semibold text-ink-soft">
									{$tStore('plants.currentMoisture')}
								</p>
								<p class="text-3xl font-bold text-ink">
									{formatMoisture(plan?.currentMoisture ?? null, p.watering?.moistureScale)}
								</p>
								{#if plan?.measuredAt}
									<p class="mt-0.5 text-xs text-ink-soft">
										{$tStore('plants.lastMeasured')}:
										<Timestamp date={new SvelteDate(plan.measuredAt)} />
									</p>
								{/if}
							</div>
							{#if p.watering?.idealMoisture !== undefined && p.watering?.idealMoisture !== null}
								<Chip
									tone="brand"
									icon="ruler"
									text={`${$tStore('plants.idealMoisture')}: ${formatMoisture(
										p.watering.idealMoisture,
										p.watering.moistureScale
									)}`}
								/>
							{/if}
						</div>

						{#if plan?.waterBelow !== undefined && plan?.waterBelow !== null && plan?.waterUntil !== undefined && plan?.waterUntil !== null}
							<p class="mt-3 text-xs text-ink-soft">
								{$tStore('plants.targetBand')}:
								<span class="font-semibold text-ink">
									{formatMoisture(plan.waterBelow, p.watering?.moistureScale)} –
									{formatMoisture(plan.waterUntil, p.watering?.moistureScale)}
								</span>
							</p>
						{/if}

						<div class="mt-4 flex flex-col gap-2">
							<Button
								variant="primary"
								class="w-full"
								icon="ruler"
								text="plants.addMeasurement"
								onclick={() => (measuring = true)}
							/>
							<Button
								variant="ghost"
								class="w-full"
								icon="line-chart"
								text="plants.moistureHistory"
								onclick={() => goto(resolve(`/plant/${p.id}/moisture`))}
							/>
						</div>
					</Card>
				{/if}

				<!-- Care overview (only set fields) -->
				{#if careRows.length > 0}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="leaf" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.careOverview')}
						</h2>
						<ul class="mt-3 divide-y divide-ink/5">
							{#each careRows as row (row.icon + row.labelKey)}
								<li class="flex items-baseline gap-3 py-2.5">
									<span class="shrink-0 self-center text-ink-soft" aria-hidden="true">
										<Icon name={row.icon} size={17} />
									</span>
									<div class="min-w-0 text-sm">
										<span class="font-semibold text-ink">{$tStore(row.labelKey)}</span>
										<span class="block text-ink-soft">{row.value}</span>
									</div>
								</li>
							{/each}
						</ul>
					</Card>
				{/if}

				<!-- Photos strip (only when photos exist) -->
				{#if p.photoIds?.length}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="camera" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.gallery')}
						</h2>
						<div class="-mx-1 mt-3 flex snap-x snap-mandatory gap-3 overflow-x-auto px-1 pb-1">
							{#each p.photoIds as photoId, i (photoId)}
								<PlantImage
									{photoId}
									remoteUrl={photoUrls[i]}
									alt={p.name}
									class="h-36 w-36 shrink-0 snap-start rounded-xl object-cover"
								/>
							{/each}
						</div>
					</Card>
				{/if}

				<!-- Notes (only when present) -->
				{#if p.notes?.length}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="file-text" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.notesTitle')}
						</h2>
						<ul class="mt-3 space-y-2 text-sm text-ink">
							{#each p.notes as note, idx (idx)}
								<li class="rounded-xl bg-canvas p-3">{note}</li>
							{/each}
						</ul>
					</Card>
				{/if}

				<!-- Growth history (only when present) -->
				{#if p.growthHistory?.length}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="sprout" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.growthHistory')}
						</h2>
						<ul class="mt-3 space-y-2">
							{#each p.growthHistory as entry (entry.id)}
								<li class="flex items-center gap-3 rounded-xl bg-canvas p-3">
									{#if entry.photoId}
										<PlantImage
											photoId={entry.photoId}
											alt={p.name}
											class="h-12 w-12 shrink-0 rounded-lg object-cover"
										/>
									{/if}
									<div class="min-w-0 flex-1 text-sm">
										<span class="font-semibold text-ink">{formatDate(entry.date)}</span>
										<span class="block text-ink-soft">
											{entry.heightCm} cm · {entry.leafCount}
											{$tStore('plants.leafCount').toLowerCase()} · {entry.health}
										</span>
									</div>
								</li>
							{/each}
						</ul>
					</Card>
				{/if}

				<!-- Pest history (only when present) -->
				{#if p.pestHistory?.length}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="bug" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.pestHistory')}
						</h2>
						<ul class="mt-3 space-y-2">
							{#each p.pestHistory as pest (pest.id)}
								<li class="rounded-xl bg-canvas p-3 text-sm">
									<div class="flex items-center justify-between gap-2">
										<span class="font-semibold text-ink">{pest.pest}</span>
										<span class="text-xs text-ink-soft">{formatDate(pest.detectedAt)}</span>
									</div>
									{#if pest.notes}
										<p class="mt-1 text-ink-soft">{pest.notes}</p>
									{/if}
								</li>
							{/each}
						</ul>
					</Card>
				{/if}
			</div>
		</Scrollable>
	{/if}
</PageContent>

{#if measuring && plant}
	<MoistureLogSheet {plant} onclose={() => (measuring = false)} />
{/if}
