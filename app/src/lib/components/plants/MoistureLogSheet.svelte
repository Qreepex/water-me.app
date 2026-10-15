<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { untrack } from 'svelte';
	import { tStore } from '$lib/i18n';
	import { getWateringStore } from '$lib/stores/watering.svelte';
	import type { Plant } from '$lib/types/api';
	import { MoistureScale } from '$lib/types/api';
	import { formatMoisture, getPlan } from '$lib/utils/plant';
	import { suggestedMoistureValue } from '$lib/utils/tasks';
	import Button from '$lib/components/ui/Button.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import MoistureInput from '$lib/components/ui/form/MoistureInput.svelte';

	interface Props {
		plant: Plant;
		/** Also offer a plain "water now" action (used in the water view). */
		allowWatering?: boolean;
		onclose: () => void;
	}

	const { plant, allowWatering = false, onclose }: Props = $props();
	const watering = getWateringStore();

	// The sheet is mounted fresh per plant, so capturing the initial value is intentional.
	let value = $state(untrack(() => suggestedMoistureValue(plant)));
	let saving = $state(false);
	let note = $state('');

	const plan = $derived(getPlan(plant));
	const ideal = $derived(plant.watering?.idealMoisture ?? null);
	const scale = $derived(plant.watering?.moistureScale ?? MoistureScale.Scale1To10);
	const tooDry = $derived(ideal !== null && value <= ideal);

	async function save(): Promise<void> {
		saving = true;
		const ok = await watering.logMoisture(plant.id, value, {
			scale,
			...(note.trim() ? { note: note.trim() } : {})
		});
		saving = false;
		if (ok) onclose();
	}

	async function saveAndWater(): Promise<void> {
		saving = true;
		const ok = await watering.waterPlant(plant.id, {
			moistures: { [plant.id]: value },
			scale
		});
		saving = false;
		if (ok) onclose();
	}
</script>

<!-- Backdrop -->
<div
	class="fixed inset-0 z-50 flex items-end justify-center bg-black/40 p-0 backdrop-blur-sm md:items-center md:p-6"
	role="presentation"
	onclick={(e) => e.target === e.currentTarget && onclose()}
>
	<Card padded class="max-h-[88vh] w-full max-w-lg overflow-y-auto rounded-b-none md:rounded-2xl">
		<div class="flex items-start justify-between gap-3">
			<div class="min-w-0">
				<p class="text-xs font-semibold tracking-wide text-brand-dark uppercase dark:text-brand">
					{$tStore('plants.moistureMeasureTitle')}
				</p>
				<h2 class="truncate text-xl font-bold text-ink">{plant.name}</h2>
				{#if plant.species}
					<p class="truncate text-sm text-ink-soft italic">{plant.species}</p>
				{/if}
			</div>
			<button
				type="button"
				onclick={onclose}
				aria-label={$tStore('common.close')}
				class="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-full bg-canvas text-ink-soft transition hover:bg-ink/10"
			>
				<Icon name="x" size={18} />
			</button>
		</div>

		<div class="mt-4">
			<MoistureInput bind:value {scale} {ideal} />
		</div>

		<!-- Live interpretation of the entered value -->
		<div
			class="mt-4 flex items-start gap-2 rounded-xl px-3 py-2.5 text-sm {tooDry
				? 'bg-danger/15 text-danger'
				: 'bg-brand/10 text-ink-soft'}"
		>
			<Icon name={tooDry ? 'alert-triangle' : 'info'} size={16} />
			<span>
				{#if ideal === null}
					{$tStore('plants.moistureNoIdealHint')}
				{:else if tooDry}
					{$tStore('plants.moistureTooDryHint', [formatMoisture(ideal, scale)])}
				{:else}
					{$tStore('plants.moistureFineHint', [formatMoisture(ideal, scale)])}
				{/if}
			</span>
		</div>

		{#if plan?.currentMoisture !== undefined && plan?.currentMoisture !== null}
			<p class="mt-2 text-xs text-ink-soft">
				{$tStore('plants.lastMeasured')}:
				<span class="font-semibold text-ink">{formatMoisture(plan.currentMoisture, scale)}</span>
			</p>
		{/if}

		<div class="mt-4">
			<input
				type="text"
				bind:value={note}
				maxlength={200}
				placeholder={$tStore('plants.moistureNotePlaceholder')}
				class="w-full rounded-xl border-2 border-brand/25 bg-canvas px-4 py-2.5 text-base text-ink shadow-sm transition placeholder:text-ink/30 focus:border-brand focus:ring-2 focus:ring-brand/20 focus:outline-none"
			/>
		</div>

		<div class="mt-5 flex flex-col gap-2">
			<Button
				variant="primary"
				size="lg"
				icon="check"
				text="plants.saveMeasurement"
				loadingText="common.save"
				loading={saving}
				disabled={saving}
				onclick={save}
			/>
			{#if tooDry || allowWatering}
				<Button
					variant="water"
					size="lg"
					icon="droplet"
					text="plants.measureAndWater"
					loadingText="plants.watering"
					loading={saving}
					onclick={saveAndWater}
				/>
			{/if}
			<Button
				variant="ghost"
				text="plants.viewMoistureHistory"
				icon="line-chart"
				onclick={() => goto(resolve(`/plant/${plant.id}/moisture`))}
			/>
		</div>
	</Card>
</div>
