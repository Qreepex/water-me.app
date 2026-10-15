<script lang="ts">
	import Field from '$lib/components/ui/form/Field.svelte';
	import Input from '$lib/components/ui/form/Input.svelte';
	import MoistureInput from '$lib/components/ui/form/MoistureInput.svelte';
	import SegmentedControl from '$lib/components/ui/form/SegmentedControl.svelte';
	import Select from '$lib/components/ui/form/Select.svelte';
	import { tStore } from '$lib/i18n';
	import { MoistureScale, WateringMethod, WateringPlanMode, WaterType } from '$lib/types/api';
	import type { FormData } from '$lib/types/forms';

	interface Props {
		formData: FormData;
	}

	let { formData = $bindable() }: Props = $props();

	const modeOptions = [
		{
			value: WateringPlanMode.Manual,
			label: 'plants.wateringModeManual',
			description: 'plants.wateringModeManualHint'
		},
		{
			value: WateringPlanMode.Moisture,
			label: 'plants.wateringModeMoisture',
			description: 'plants.wateringModeMoistureHint'
		}
	];

	const methodOptions = $derived.by(() =>
		Object.values(WateringMethod).map((method) => ({
			value: method,
			label: $tStore('plants.wateringMethodOptions.' + method)
		}))
	);

	const typeOptions = $derived.by(() =>
		Object.values(WaterType).map((type) => ({
			value: type,
			label: $tStore('plants.waterTypeOptions.' + type)
		}))
	);

	const scaleOptions = $derived.by(() => [
		{ value: MoistureScale.Scale1To10, label: $tStore('plants.moistureScaleTen') },
		{ value: MoistureScale.Percent, label: $tStore('plants.moistureScalePercent') }
	]);

	const isMoisture = $derived(formData.wateringMode === WateringPlanMode.Moisture);
</script>

<div class="space-y-5">
	<Field label="plants.wateringMode" hint="plants.wateringModeHint">
		<SegmentedControl
			bind:value={formData.wateringMode}
			label="plants.wateringMode"
			options={modeOptions}
		/>
	</Field>

	{#if isMoisture}
		<!-- Moisture-based plan: ideal value + scale, the scheduler does the rest -->
		<div class="rounded-xl bg-brand/5 p-3 text-sm text-ink-soft">
			{$tStore('plants.moisturePlanHint')}
		</div>

		<Field label="plants.moistureScale" forId="moisture-scale">
			<Select id="moisture-scale" bind:value={formData.moistureScale} options={scaleOptions} />
		</Field>

		<div>
			<p class="mb-1 block text-sm font-semibold text-ink-soft">
				{$tStore('plants.idealMoisture')}
			</p>
			<MoistureInput
				id="ideal-moisture"
				bind:value={formData.idealMoisture}
				scale={formData.moistureScale}
				ideal={formData.idealMoisture}
				band={formData.moistureBand}
			/>
			<p class="mt-1 text-xs text-ink-soft/70">{$tStore('plants.idealMoistureHint')}</p>
		</div>

		<Field
			label="plants.moistureBand"
			forId="moisture-band"
			labelSuffix="{formData.idealMoisture - formData.moistureBand}–{formData.idealMoisture +
				formData.moistureBand}"
			hint="plants.moistureBandHint"
		>
			<Input
				id="moisture-band"
				type="number"
				min={0.5}
				step={0.5}
				bind:value={formData.moistureBand}
			/>
		</Field>
	{:else}
		<Field label="plants.formIntervalDays" forId="water-interval">
			<Input id="water-interval" type="number" min={1} bind:value={formData.wateringIntervalDays} />
		</Field>
	{/if}

	<Field label="plants.wateringMethod" forId="water-method">
		<Select id="water-method" bind:value={formData.wateringMethod} options={methodOptions} />
	</Field>

	<Field label="plants.waterType" forId="water-type">
		<Select id="water-type" bind:value={formData.waterType} options={typeOptions} />
	</Field>
</div>
