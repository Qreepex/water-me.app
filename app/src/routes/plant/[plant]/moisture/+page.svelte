<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { untrack } from 'svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Scrollable from '$lib/components/layout/Scrollable.svelte';
	import MoistureLogSheet from '$lib/components/plants/MoistureLogSheet.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Chip from '$lib/components/ui/Chip.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { getWateringStore } from '$lib/stores/watering.svelte';
	import type { MoistureLog } from '$lib/types/api';
	import { MoistureScale } from '$lib/types/api';
	import { formatMoisture, getMoistureBounds, getPlan } from '$lib/utils/plant';
	import Timestamp from '$lib/components/ui/Timestamp.svelte';
	import { SvelteDate } from 'svelte/reactivity';

	const store = getPlantsStore();
	const watering = getWateringStore();

	const plantId = $derived(page.params.plant ?? '');
	const plant = $derived(store.plants.find((p) => p.id === plantId || p.slug === plantId) ?? null);
	const plan = $derived(plant ? getPlan(plant) : null);

	let deletingId = $state<string | null>(null);

	const scale = $derived(plant?.watering?.moistureScale ?? MoistureScale.Scale1To10);
	const ideal = $derived(plant?.watering?.idealMoisture ?? null);
	const band = $derived(plant?.watering?.moistureBand ?? null);
	/** Lower edge of the target band: the server waters here. */
	const below = $derived(plan?.waterBelow ?? (ideal === null ? null : ideal - (band ?? 2)));
	const until = $derived(plan?.waterUntil ?? (ideal === null ? null : ideal + (band ?? 2)));
	const bounds = $derived(plant ? getMoistureBounds(plant) : { min: 1, max: 10, step: 0.5 });

	/** Newest first. */
	const logs = $derived(
		[...(plant?.moistureLogs ?? [])].sort(
			(a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
		)
	);

	function moisturePct(value: number): number {
		return Math.min(100, Math.max(0, ((value - bounds.min) / (bounds.max - bounds.min)) * 100));
	}

	const idealPct = $derived(ideal === null ? null : moisturePct(ideal));
	const belowPct = $derived(below === null ? null : moisturePct(below));
	const untilPct = $derived(until === null ? null : moisturePct(until));

	function barTone(value: number): string {
		if (below === null) return 'bg-brand';
		if (value < below) return 'bg-danger';
		if (ideal !== null && value <= ideal) return 'bg-warn';
		return 'bg-brand';
	}

	async function remove(log: MoistureLog): Promise<void> {
		if (!plant) return;
		if (!confirm($tStore('plants.deleteMeasurementConfirm'))) return;
		deletingId = log.id;
		await watering.deleteMoisture(plant.id, log.id);
		deletingId = null;
	}

	function trendIcon(): 'trending-down' | 'trending-up' | 'arrow-right' {
		if (!plan?.decayPerDay) return 'arrow-right';
		return plan.decayPerDay > 0.05
			? 'trending-down'
			: plan.decayPerDay < -0.05
				? 'trending-up'
				: 'arrow-right';
	}

	// Deep link: /plant/:id/moisture?measure=1 opens the entry sheet right away.
	let logging = $state(untrack(() => page.url.searchParams.get('measure') === '1'));
</script>

<PageHeader
	icon="line-chart"
	title="plants.moistureHistory"
	description={plant?.name}
	backHref={plant ? resolve(`/plant/${plant.id}`) : resolve('/')}
>
	{#if plant}
		<Button
			variant="primary"
			size="sm"
			text="plants.addMeasurement"
			icon="plus"
			onclick={() => (logging = true)}
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
				<!-- Current state + the algorithm's prediction -->
				<Card padded>
					<div class="flex items-start justify-between gap-3">
						<div>
							<p class="text-sm font-semibold text-ink-soft">{$tStore('plants.currentMoisture')}</p>
							<p class="text-4xl font-bold text-ink">
								{formatMoisture(plan?.currentMoisture ?? null, scale)}
							</p>
							{#if plan?.measuredAt}
								<p class="mt-1 text-xs text-ink-soft">
									{$tStore('plants.lastMeasured')}:
									<Timestamp date={new SvelteDate(plan.measuredAt)} />
								</p>
							{/if}
							{#if plan?.wateredSinceLastMeasurement}
								<p class="mt-1 text-xs text-warn">
									{$tStore('plants.measurementStaleHint')}
								</p>
							{/if}
						</div>
						{#if ideal !== null}
							<Chip
								tone="brand"
								icon="ruler"
								text={`${$tStore('plants.idealMoisture')}: ${formatMoisture(ideal, scale)}`}
							/>
						{/if}
					</div>

					<!-- Scale bar with the ideal marker -->
					{#if plan?.currentMoisture !== undefined && plan?.currentMoisture !== null}
						<div class="mt-4">
							<div class="relative h-3 overflow-hidden rounded-full bg-ink/10">
								<div
									class="h-full rounded-full transition-all duration-500 {barTone(
										plan.currentMoisture
									)}"
									style:width="{moisturePct(plan.currentMoisture)}%"
								></div>
							</div>
							{#if idealPct !== null}
								<div class="relative h-0" aria-hidden="true">
									{#if belowPct !== null && untilPct !== null}
										<div
											class="absolute -top-3 h-3 rounded-sm bg-brand/35 ring-1 ring-brand-dark/40"
											style:left="{belowPct}%; width: {untilPct - belowPct}%"
										></div>
									{/if}
									<div
										class="absolute -top-5 h-7 w-0.5 -translate-x-1/2 bg-ink/70"
										style:left="{idealPct}%"
									></div>
								</div>
							{/if}
						</div>
					{/if}

					<!-- Model summary -->
					{#if plan}
						<dl class="mt-5 grid grid-cols-2 gap-3 text-sm">
							<div class="rounded-xl bg-canvas p-3">
								<dt class="flex items-center gap-1.5 text-xs text-ink-soft">
									<Icon name={trendIcon()} size={13} />
									{$tStore('plants.decayPerDay')}
								</dt>
								<dd class="font-semibold text-ink">
									{plan.decayPerDay !== undefined && plan.decayPerDay !== null
										? `${Math.round(plan.decayPerDay * 10) / 10}/day`
										: '—'}
								</dd>
							</div>
							<div class="rounded-xl bg-canvas p-3">
								<dt class="text-xs text-ink-soft">{$tStore('plants.targetBand')}</dt>
								<dd class="font-semibold text-ink">
									{#if below !== null && until !== null}
										{formatMoisture(below, scale)} – {formatMoisture(until, scale)}
									{:else}
										—
									{/if}
								</dd>
							</div>
							<div class="rounded-xl bg-canvas p-3">
								<dt class="text-xs text-ink-soft">{$tStore('plants.measurementsUsed')}</dt>
								<dd class="font-semibold text-ink">{plan.sampleCount ?? 0}</dd>
							</div>
							<div class="rounded-xl bg-canvas p-3">
								<dt class="text-xs text-ink-soft">{$tStore('plants.nextWatering')}</dt>
								<dd class="font-semibold text-ink">
									{#if plan.wateringDueAt}
										<Timestamp date={new SvelteDate(plan.wateringDueAt)} mode="future" />
									{:else}
										{$tStore('plants.notEnoughDataYet')}
									{/if}
								</dd>
							</div>
							<div class="rounded-xl bg-canvas p-3">
								<dt class="text-xs text-ink-soft">{$tStore('plants.nextMeasurement')}</dt>
								<dd class="font-semibold text-ink">
									{#if plan.measurementDueAt}
										<Timestamp date={new SvelteDate(plan.measurementDueAt)} mode="future" />
									{:else}
										—
									{/if}
								</dd>
							</div>
						</dl>
					{/if}

					<Button
						variant="primary"
						class="mt-4 w-full"
						icon="ruler"
						text="plants.addMeasurement"
						onclick={() => (logging = true)}
					/>
				</Card>

				<!-- History -->
				{#if logs.length === 0}
					<Card padded class="text-center">
						<p class="text-sm text-ink-soft">{$tStore('plants.noMeasurementsYet')}</p>
						<p class="mt-1 text-xs text-ink-soft/80">{$tStore('plants.moistureHistoryHint')}</p>
					</Card>
				{:else}
					<Card padded>
						<h2 class="flex items-center gap-2 text-base font-bold text-ink">
							<Icon name="line-chart" size={18} class="text-brand-dark dark:text-brand" />
							{$tStore('plants.moistureHistory')}
						</h2>
						<ul class="mt-3 divide-y divide-ink/5">
							{#each logs as log (log.id)}
								<li class="flex items-center gap-3 py-3">
									<div class="min-w-0 flex-1">
										<div class="flex items-baseline gap-2">
											<span class="text-lg font-bold text-ink">
												{formatMoisture(log.value, log.scale ?? scale)}
											</span>
											{#if below !== null && log.value < below}
												<span
													class="rounded-full bg-danger/15 px-2 py-0.5 text-[11px] font-semibold text-danger"
												>
													{$tStore('plants.belowBand')}
												</span>
											{/if}
										</div>
										<p class="text-xs text-ink-soft">
											<Timestamp date={new SvelteDate(log.date)} />
											{#if log.note}· {log.note}{/if}
										</p>
									</div>
									<div class="h-2 w-16 shrink-0 overflow-hidden rounded-full bg-ink/10">
										<div
											class="h-full rounded-full {barTone(log.value)}"
											style:width="{moisturePct(log.value)}%"
										></div>
									</div>
									<button
										type="button"
										onclick={() => remove(log)}
										disabled={deletingId === log.id}
										aria-label={$tStore('common.delete')}
										class="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-full text-ink-soft transition hover:bg-danger/10 hover:text-danger disabled:opacity-40"
									>
										<Icon name="trash" size={16} />
									</button>
								</li>
							{/each}
						</ul>
					</Card>
				{/if}

				<!-- Why this plan? -->
				<Card padded class="bg-brand/5">
					<h2 class="flex items-center gap-2 text-base font-bold text-ink">
						<Icon name="info" size={18} class="text-brand-dark dark:text-brand" />
						{$tStore('plants.howItWorks')}
					</h2>
					<p class="mt-2 text-sm text-ink-soft">{$tStore('plants.moisturePlanHint')}</p>
					<Button
						variant="ghost"
						size="sm"
						class="mt-3"
						icon="settings"
						text="plants.editWateringPlan"
						onclick={() => goto(resolve(`/manage/${p.id}/edit#watering`))}
					/>
				</Card>
			</div>
		</Scrollable>
	{/if}
</PageContent>

{#if logging && plant}
	<MoistureLogSheet {plant} onclose={() => (logging = false)} />
{/if}
