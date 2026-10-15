<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import List from '$lib/components/layout/List.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Scrollable from '$lib/components/layout/Scrollable.svelte';
	import PlantImage from '$lib/components/plants/PlantImage.svelte';
	import WaterPlantCard from '$lib/components/plants/WaterPlantCard.svelte';
	import MoistureLogSheet from '$lib/components/plants/MoistureLogSheet.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import HoldButton from '$lib/components/ui/HoldButton.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import MoodFace from '$lib/components/ui/MoodFace.svelte';
	import SectionHeader from '$lib/components/ui/SectionHeader.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { getWateringStore } from '$lib/stores/watering.svelte';
	import type { Plant } from '$lib/types/api';
	import { formatMoisture, getPlan } from '$lib/utils/plant';
	import { buildTaskList } from '$lib/utils/tasks';

	const store = getPlantsStore();
	const watering = getWateringStore();

	const tasks = $derived(buildTaskList(store.plants));
	const waterTargets = $derived(tasks.water.map((t) => t.plant.id));
	let measuringPlant = $state<Plant | null>(null);

	/**
	 * Water everything that is due.
	 *
	 * Deliberately no moisture values are sent: a watering is recorded as its
	 * own event, and re-submitting the last reading would add a duplicate data
	 * point at the watering moment and distort the model.
	 */
	async function waterAllDue(): Promise<void> {
		await watering.waterPlants(waterTargets);
	}
</script>

<PageHeader icon="list-checks" title="tasks.title" description="tasks.description" />

<PageContent>
	{#if store.error}
		<div class="mx-2 mb-3 shrink-0">
			<Alert type="error" title="common.error" description={store.error} />
		</div>
	{/if}

	{#if store.loading}
		<Spinner message="common.loadingPlants" />
	{:else if store.plants.length === 0}
		<EmptyState icon="leaf" title="plants.noPlants" description="plants.startAddingPlants">
			<Button
				variant="primary"
				onclick={() => goto(resolve('/manage/new'))}
				text="plants.addPlant"
			/>
		</EmptyState>
	{:else}
		<Scrollable noPadding>
			<!-- extra bottom padding reserves space for the floating water-all bar -->
			<div class="space-y-6 px-2 {tasks.water.length > 0 ? 'pb-44' : 'pb-8'}">
				{#if tasks.total === 0 && tasks.unconfigured.length === 0}
					<!-- Everything is taken care of -->
					<Card padded class="pop-in border-none bg-brand/15 text-center">
						<div class="flex justify-center">
							<MoodFace mood="happy" size={72} class="animate-pulse text-ok" />
						</div>
						<p class="mt-3 text-xl font-bold text-ink">{$tStore('tasks.allDone')}</p>
						<p class="mt-1 text-sm text-ink-soft">{$tStore('tasks.allDoneHint')}</p>
					</Card>
				{/if}

				<!-- 1. Watering tasks -->
				{#if tasks.water.length > 0}
					<section>
						<SectionHeader
							icon="droplet"
							title="tasks.waterSection"
							count={tasks.water.length}
							tone="danger"
						/>
						<List noPadding class="md:grid-cols-2 2xl:grid-cols-3">
							{#each tasks.water as task, i (task.plant.id)}
								<WaterPlantCard plant={task.plant} index={i} />
							{/each}
						</List>
					</section>
				{/if}

				<!-- 2. Measurement tasks -->
				{#if tasks.measure.length > 0}
					<section>
						<SectionHeader
							icon="ruler"
							title="tasks.measureSection"
							count={tasks.measure.length}
							tone="info"
						/>
						<List noPadding class="md:grid-cols-2 2xl:grid-cols-3">
							{#each tasks.measure as task, i (task.plant.id)}
								{@const plan = getPlan(task.plant)}
								{@const scale = task.plant.watering?.moistureScale}
								<button
									type="button"
									onclick={() => (measuringPlant = task.plant)}
									class="pop-in flex w-full cursor-pointer items-center gap-3 rounded-2xl bg-surface p-3.5 text-left shadow-md transition hover:shadow-lg active:scale-[0.99]"
									style:animation-delay="{Math.min(i * 40, 280)}ms"
								>
									<div class="relative h-14 w-14 shrink-0 overflow-hidden rounded-xl">
										<PlantImage
											photoId={task.plant.photoIds?.[0]}
											remoteUrl={(task.plant as { photoUrls?: string[] }).photoUrls?.[0]}
											alt={task.plant.name}
											class="h-full w-full object-cover"
											fallbackClass="[&_svg]:h-6 [&_svg]:w-6"
										/>
									</div>
									<div class="min-w-0 flex-1">
										<p class="truncate font-semibold text-ink">{task.plant.name}</p>
										<p class="truncate text-xs text-info">{$tStore(task.reasonKey)}</p>
										{#if plan?.currentMoisture !== undefined && plan?.currentMoisture !== null}
											<p class="truncate text-xs text-ink-soft">
												{$tStore('plants.lastMeasured')}:
												{formatMoisture(plan.currentMoisture, scale)}
											</p>
										{/if}
									</div>
									<span
										class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-info/15 text-info"
									>
										<Icon name="ruler" size={20} />
									</span>
								</button>
							{/each}
						</List>
					</section>
				{/if}

				<!-- 3. Nothing to do today -->
				{#if tasks.relaxed.length > 0}
					<section>
						<SectionHeader
							icon="check-circle"
							title="tasks.relaxedSection"
							count={tasks.relaxed.length}
							tone="ok"
						/>
						<List noPadding class="md:grid-cols-2 2xl:grid-cols-3">
							{#each tasks.relaxed as plant, i (plant.id)}
								<WaterPlantCard {plant} compact index={i} />
							{/each}
						</List>
					</section>
				{/if}

				<!-- 4. No plan configured yet -->
				{#if tasks.unconfigured.length > 0}
					<section>
						<SectionHeader
							icon="settings"
							title="plants.noWateringConfig"
							count={tasks.unconfigured.length}
							tone="neutral"
						/>
						<List noPadding>
							{#each tasks.unconfigured as plant (plant.id)}
								<div
									class="pop-in flex items-center justify-between gap-3 rounded-2xl bg-surface p-4 shadow"
								>
									<div class="min-w-0 flex-1">
										<p class="truncate font-semibold text-ink">{plant.name}</p>
										<p class="truncate text-xs text-ink-soft">{plant.species}</p>
									</div>
									<Button
										variant="secondary"
										size="sm"
										text="plants.configureWateringSchedule"
										icon="droplet"
										onclick={() => goto(resolve(`/manage/${plant.id}/edit#watering`))}
									/>
								</div>
							{/each}
						</List>
					</section>
				{/if}
			</div>
		</Scrollable>

		<!-- Floating water-all action, always within thumb reach -->
		{#if tasks.water.length > 0}
			<div
				class="fixed right-3 left-3 z-50 md:right-10 md:left-10 xl:right-32 xl:left-32"
				style="bottom: calc(env(safe-area-inset-bottom) + 5.5rem);"
			>
				<HoldButton
					variant="water"
					size="lg"
					icon="droplets"
					label={$tStore('plants.waterAllPlants', [String(tasks.water.length)])}
					disabled={watering.isWateringAny}
					onconfirm={waterAllDue}
					class="shadow-xl"
				/>
			</div>
		{/if}
	{/if}
</PageContent>

{#if measuringPlant}
	<MoistureLogSheet plant={measuringPlant} allowWatering onclose={() => (measuringPlant = null)} />
{/if}
