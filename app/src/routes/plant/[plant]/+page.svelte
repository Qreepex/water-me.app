<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { imageCacheStore } from '$lib/stores/imageCache.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import Alert from '$lib/components/ui/Message.svelte';
	import { resolve } from '$app/paths';
	import { SvelteMap, SvelteSet } from 'svelte/reactivity';
	import Scrollable from '$lib/components/layout/Scrollable.svelte';
	import { formatPastTimestamp } from '$lib/utils/timestamp.svelte';

	const store = getPlantsStore();
	const plantId = $derived($page.params.plant ?? '');
	const plant = $derived.by(
		() => store.plants.find((p) => p.id === plantId || p.slug === plantId) ?? null
	);

	let imageUrls = $state<Record<string, string | null>>({});
	let loadingImages = $state(false);

	const primaryPhotoId = $derived.by(() => plant?.photoIds?.[0] ?? null);
	const primaryPhotoUrl = $derived.by(() =>
		primaryPhotoId ? (imageUrls[primaryPhotoId] ?? null) : null
	);

	$effect(() => {
		if (!plant) return;

		const ids = new SvelteSet<string>();
		(plant.photoIds ?? []).forEach((id) => id && ids.add(id));
		(plant.growthHistory ?? []).forEach((entry) => entry.photoId && ids.add(entry.photoId));

		const idList = Array.from(ids);
		const remoteUrls = ((plant as unknown as { photoUrls?: string[] })?.photoUrls ??
			[]) as string[];
		const idToRemote = new SvelteMap<string, string>();
		(plant.photoIds ?? []).forEach((id, index) => {
			const url = remoteUrls[index];
			if (id && url) idToRemote.set(id, url);
		});

		loadingImages = idList.length > 0;
		let cancelled = false;
		const nextUrls: Record<string, string | null> = {};

		const loadAll = async () => {
			for (const id of idList) {
				let url = imageCacheStore.getImageURLSync(id);
				if (!url) {
					url = await imageCacheStore.getImageURL(id, idToRemote.get(id));
				}
				if (cancelled) return;
				nextUrls[id] = url;
			}
			if (!cancelled) imageUrls = nextUrls;
		};

		loadAll().finally(() => {
			if (!cancelled) loadingImages = false;
		});

		return () => {
			cancelled = true;
			idList.forEach((id) => imageCacheStore.releaseImage(id));
		};
	});

	function formatDate(date?: string | null): string {
		if (!date) return '-';
		return new Date(date).toLocaleDateString();
	}

	function editPlant() {
		if (!plant) return;
		goto(resolve(`/manage/${plant.id}`));
	}

	function formatBool(value?: boolean): string {
		if (value === undefined || value === null) return '-';
		return value ? ($tStore('common.yes') ?? 'Yes') : ($tStore('common.no') ?? 'No');
	}
</script>

<PageHeader icon="🪴" title={plant?.name ?? 'plants.myPlants'} description={plant?.species}>
		<Button variant="primary" size="sm" text="common.edit" icon="✏️" onclick={editPlant} />
</PageHeader>

<PageContent>
	{#if store.loading}
		<LoadingSpinner message="plants.loadingPlants" icon="🌿" />
	{:else if !plant}
		<Alert type="error" title="common.error" description="plants.notFound" />
	{:else}
		<Scrollable>
			<Card rounded="2xl">
				<div class="grid gap-6 md:grid-cols-[minmax(0,320px)_1fr]">
					<div
						class="flex h-64 items-center justify-center overflow-hidden rounded-t-2xl bg-gradient-to-br from-[var(--p-emerald)] to-[var(--p-emerald-dark)] md:h-full md:rounded-l-2xl md:rounded-tr-none"
					>
						{#if primaryPhotoUrl}
							<img src={primaryPhotoUrl} alt={plant.name} class="h-full w-full object-cover" />
						{:else}
							<div class="text-7xl">🌿</div>
						{/if}
					</div>
					<div class="p-6">
						<div class="flex flex-wrap items-center gap-2">
							{#if plant.sunlight}
								<span
									class="rounded-full bg-[var(--status-info)]/20 px-3 py-1 text-xs font-medium text-[var(--status-info)]"
								>
									☀️ {$tStore('plants.sunlight.' + plant.sunlight)}
								</span>
							{/if}
							{#if plant.isToxic !== undefined}
								<span
									class="rounded-full bg-[var(--status-warn)]/20 px-3 py-1 text-xs font-medium text-[var(--status-warn)]"
								>
									☠️ {plant.isToxic ? $tStore('plants.toxic') : $tStore('plants.notToxic')}
								</span>
							{/if}
							{#if plant.location?.room}
								<span
									class="rounded-full bg-[var(--p-emerald)]/20 px-3 py-1 text-xs font-medium text-[var(--p-emerald-dark)]"
								>
									📍 {plant.location.room}
									{#if plant.location?.position}
										· {plant.location.position}
									{/if}
								</span>
							{/if}
							{#if plant.location?.isOutdoors !== undefined}
								<span
									class="rounded-full bg-[var(--status-success)]/20 px-3 py-1 text-xs font-medium text-[var(--status-success)]"
								>
									{plant.location.isOutdoors
										? $tStore('plants.outdoors')
										: $tStore('plants.indoors')}
								</span>
							{/if}
						</div>

						<div class="mt-6 grid grid-cols-2 gap-4">
							<div class="rounded-xl bg-[var(--p-emerald)]/10 p-4">
								<div class="text-xs font-semibold text-[var(--p-emerald-dark)]">
									{$tStore('plants.lastWatered')}
								</div>
								<div class="mt-1 text-sm text-[var(--text-light-main)]">
									{plant.watering?.lastWatered
										? formatPastTimestamp(new Date(plant.watering.lastWatered))
										: '-'}
								</div>
							</div>
							<div class="rounded-xl bg-[var(--status-warn)]/10 p-4">
								<div class="text-xs font-semibold text-[var(--status-warn)]">
									{$tStore('plants.lastFertilized')}
								</div>
								<div class="mt-1 text-sm text-[var(--text-light-main)]">
									{plant.fertilizing?.lastFertilized
										? formatPastTimestamp(new Date(plant.fertilizing.lastFertilized))
										: '-'}
								</div>
							</div>
							<div class="rounded-xl bg-[var(--status-info)]/10 p-4">
								<div class="text-xs font-semibold text-[var(--status-info)]">
									{$tStore('plants.wateringFrequency')}
								</div>
								<div class="mt-1 text-sm text-[var(--text-light-main)]">
									{plant.watering?.intervalDays
										? `${$tStore('plants.every')} ${plant.watering.intervalDays} ${$tStore('plants.days')}`
										: '-'}
								</div>
							</div>
							<div class="rounded-xl bg-[var(--p-emerald)]/10 p-4">
								<div class="text-xs font-semibold text-[var(--p-emerald-dark)]">
									{$tStore('plants.sprayFrequency')}
								</div>
								<div class="mt-1 text-sm text-[var(--text-light-main)]">
									{plant.humidity?.mistingIntervalDays
										? `${$tStore('plants.every')} ${plant.humidity.mistingIntervalDays} ${$tStore('plants.days')}`
										: '-'}
								</div>
							</div>
						</div>
					</div>
				</div>
			</Card>

			<Card rounded="2xl">
				<div class="p-6">
					<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
						{$tStore('plants.careDetails')}
					</h2>
					<div class="mt-4 grid gap-4 md:grid-cols-2">
						<div class="rounded-xl border border-[var(--p-emerald)]/20 bg-white/60 p-4">
							<div class="text-sm font-semibold text-[var(--p-emerald-dark)]">
								{$tStore('plants.wateringTitle')}
							</div>
							<div class="mt-2 space-y-1 text-sm text-[var(--text-light-main)]">
								<div>
									{$tStore('plants.wateringMethod')}: {plant.watering?.method ?? '-'}
								</div>
								<div>
									{$tStore('plants.waterType')}: {plant.watering?.waterType ?? '-'}
								</div>
								<div>
									{$tStore('plants.lastWatered')}: {formatDate(plant.watering?.lastWatered)}
								</div>
							</div>
						</div>

						{#if plant.fertilizing}
							<div class="rounded-xl border border-[var(--status-warn)]/20 bg-white/60 p-4">
								<div class="text-sm font-semibold text-[var(--status-warn)]">
									{$tStore('plants.fertilizingTitle')}
								</div>
								<div class="mt-2 space-y-1 text-sm text-[var(--text-light-main)]">
									<div>
										{$tStore('plants.fertilizerType')}: {plant.fertilizing?.type ?? '-'}
									</div>
									<div>
										{$tStore('plants.npkRatio')}: {plant.fertilizing?.npkRatio ?? '-'}
									</div>
									<div>
										{$tStore('plants.concentration')}: {plant.fertilizing?.concentrationPercent ??
											'-'}%
									</div>
									<div>
										{$tStore('plants.lastFertilized')}: {formatDate(
											plant.fertilizing?.lastFertilized
										)}
									</div>
								</div>
							</div>
						{/if}

						{#if plant.humidity}
							<div class="rounded-xl border border-[var(--status-info)]/20 bg-white/60 p-4">
								<div class="text-sm font-semibold text-[var(--status-info)]">
									{$tStore('plants.humidityTitle')}
								</div>
								<div class="mt-2 space-y-1 text-sm text-[var(--text-light-main)]">
									<div>
										{$tStore('plants.targetHumidity')}: {plant.humidity?.targetHumidityPct ?? '-'}%
									</div>
									<div>
										{$tStore('plants.mistingInterval')}: {plant.humidity?.mistingIntervalDays ??
											'-'}
									</div>
									<div>
										{$tStore('plants.lastMisted')}: {formatDate(plant.humidity?.lastMisted)}
									</div>
									<div>
										{$tStore('plants.humidifier')}: {formatBool(plant.humidity?.requiresHumidifier)}
									</div>
								</div>
							</div>
						{/if}

						{#if plant.soil}
							<div class="rounded-xl border border-[var(--p-emerald)]/20 bg-white/60 p-4">
								<div class="text-sm font-semibold text-[var(--p-emerald-dark)]">
									{$tStore('plants.soilTitle')}
								</div>
								<div class="mt-2 space-y-1 text-sm text-[var(--text-light-main)]">
									<div>{$tStore('plants.soilType')}: {plant.soil?.type ?? '-'}</div>
									<div>
										{$tStore('plants.repottingCycle')}: {plant.soil?.repottingCycle ?? '-'}
									</div>
									<div>
										{$tStore('plants.lastRepotted')}: {formatDate(plant.soil?.lastRepotted)}
									</div>
									{#if plant.soil?.components?.length}
										<div>
											{$tStore('plants.soilComponents')}: {plant.soil.components.join(', ')}
										</div>
									{/if}
								</div>
							</div>
						{/if}

						{#if plant.seasonality}
							<div class="rounded-xl border border-[var(--status-warn)]/20 bg-white/60 p-4">
								<div class="text-sm font-semibold text-[var(--status-warn)]">
									{$tStore('plants.seasonalityTitle')}
								</div>
								<div class="mt-2 space-y-1 text-sm text-[var(--text-light-main)]">
									<div>
										{$tStore('plants.winterRest')}: {formatBool(
											plant.seasonality?.winterRestPeriod
										)}
									</div>
									<div>
										{$tStore('plants.winterWaterFactor')}: {plant.seasonality?.winterWaterFactor ??
											'-'}
									</div>
									<div>
										{$tStore('plants.minTemp')}: {plant.seasonality?.minTempCelsius ?? '-'}°C
									</div>
								</div>
							</div>
						{/if}
					</div>
				</div>
			</Card>

			<Card rounded="2xl">
				<div class="p-6">
					<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
						{$tStore('plants.gallery')}
					</h2>
					{#if loadingImages}
						<p class="mt-2 text-sm text-[var(--text-light-main)]/60">
							{$tStore('plants.loadingImages')}
						</p>
					{:else if plant.photoIds?.length}
						<div class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
							{#each plant.photoIds as photoId (photoId)}
								<div class="aspect-square overflow-hidden rounded-xl bg-[var(--p-emerald)]/10">
									{#if imageUrls[photoId]}
										<img
											src={imageUrls[photoId]}
											alt={plant.name}
											class="h-full w-full object-cover"
										/>
									{:else}
										<div class="flex h-full w-full items-center justify-center text-2xl">🌿</div>
									{/if}
								</div>
							{/each}
						</div>
					{:else}
						<p class="mt-2 text-sm text-[var(--text-light-main)]/60">
							{$tStore('plants.noImages')}
						</p>
					{/if}
				</div>
			</Card>

			{#if plant.notes?.length}
				<Card rounded="2xl">
					<div class="p-6">
						<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
							{$tStore('plants.notesTitle')}
						</h2>
						<ul class="mt-4 space-y-2 text-sm text-[var(--text-light-main)]">
							{#each plant.notes as note, idx (idx)}
								<li class="rounded-lg bg-white/70 p-3">📝 {note}</li>
							{/each}
						</ul>
					</div>
				</Card>
			{/if}

			{#if plant.flags?.length}
				<Card rounded="2xl">
					<div class="p-6">
						<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
							{$tStore('plants.flags')}
						</h2>
						<div class="mt-4 flex flex-wrap gap-2">
							{#each plant.flags as flag (flag)}
								<span
									class="rounded-full bg-[var(--status-warn)]/20 px-3 py-1 text-xs font-medium text-[var(--status-warn)]"
								>
									⚡ {flag}
								</span>
							{/each}
						</div>
					</div>
				</Card>
			{/if}

			{#if plant.pestHistory?.length}
				<Card rounded="2xl">
					<div class="p-6">
						<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
							{$tStore('plants.pestHistory')}
						</h2>
						<div class="mt-4 space-y-3">
							{#each plant.pestHistory as pest (pest.id)}
								<div class="rounded-xl border border-[var(--status-warn)]/20 bg-white/70 p-4">
									<div class="flex items-center justify-between">
										<div class="font-semibold text-[var(--text-light-main)]">{pest.pest}</div>
										<div class="text-xs text-[var(--text-light-main)]/60">
											{formatDate(pest.detectedAt)}
										</div>
									</div>
									<p class="mt-2 text-sm text-[var(--text-light-main)]/80">
										{pest.notes}
									</p>
								</div>
							{/each}
						</div>
					</div>
				</Card>
			{/if}

			<Card rounded="2xl">
				<div class="p-6">
					<h2 class="text-lg font-semibold text-[var(--text-light-main)]">
						{$tStore('plants.growthHistory')}
					</h2>
					{#if plant.growthHistory?.length}
						<div class="mt-4 space-y-3">
							{#each plant.growthHistory as entry (entry.id)}
								<div class="rounded-xl border border-[var(--p-emerald)]/20 bg-white/70 p-4">
									<div class="flex flex-wrap items-center justify-between gap-2">
										<div class="text-sm font-semibold text-[var(--text-light-main)]">
											{formatDate(entry.date)}
										</div>
										<div class="text-xs text-[var(--text-light-main)]/60">
											{$tStore('plants.health')}: {entry.health}
										</div>
									</div>
									<div class="mt-3 grid gap-3 md:grid-cols-[96px_1fr]">
										{#if entry.photoId}
											<div class="h-24 w-24 overflow-hidden rounded-lg bg-[var(--p-emerald)]/10">
												{#if imageUrls[entry.photoId]}
													<img
														src={imageUrls[entry.photoId]}
														alt={plant.name}
														class="h-full w-full object-cover"
													/>
												{:else}
													<div class="flex h-full w-full items-center justify-center">🌱</div>
												{/if}
											</div>
										{/if}
										<div class="text-sm text-[var(--text-light-main)]">
											<div>{$tStore('plants.height')}: {entry.heightCm} cm</div>
											<div>{$tStore('plants.leafCount')}: {entry.leafCount}</div>
											{#if entry.condition}
												<div>{$tStore('plants.condition')}: {entry.condition}</div>
											{/if}
										</div>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="mt-2 text-sm text-[var(--text-light-main)]/60">
							{$tStore('plants.noGrowthHistory')}
						</p>
					{/if}
				</div>
			</Card>
		</Scrollable>
	{/if}
</PageContent>
