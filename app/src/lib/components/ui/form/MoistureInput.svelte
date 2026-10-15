<script lang="ts">
	import { tStore } from '$lib/i18n';
	import { MoistureScale } from '$lib/types/api';

	interface Props {
		value: number;
		scale?: MoistureScale;
		/** Centre of the target band. Water at `ideal - band`, up to `ideal + band`. */
		ideal?: number | null;
		/** Half the width of the target band. Defaults to 2. */
		band?: number | null;
		disabled?: boolean;
		id?: string;
	}

	let {
		value = $bindable(),
		scale = MoistureScale.Scale1To10,
		ideal = null,
		band = null,
		disabled = false,
		id
	}: Props = $props();

	/** Default band half-width, mirroring services.DefaultMoistureBand. */
	const DEFAULT_BAND = 2;

	const bounds = $derived(
		scale === MoistureScale.Percent ? { min: 0, max: 100, step: 1 } : { min: 1, max: 10, step: 0.5 }
	);

	const toPct = (v: number) => ((v - bounds.min) / (bounds.max - bounds.min)) * 100;
	const pct = $derived(toPct(value));

	const effectiveBand = $derived(band && band > 0 ? band : DEFAULT_BAND);
	/** Lower edge: watering is recommended here. */
	const below = $derived(ideal === null || ideal === undefined ? null : ideal - effectiveBand);
	/** Upper edge: water up to here. */
	const until = $derived(ideal === null || ideal === undefined ? null : ideal + effectiveBand);

	const bandLeftPct = $derived(below === null ? null : Math.max(0, Math.min(100, toPct(below))));
	const bandWidthPct = $derived(
		below === null || until === null
			? null
			: Math.max(0, Math.min(100, toPct(until)) - toPct(below))
	);
	const idealPct = $derived(ideal === null || ideal === undefined ? null : toPct(ideal));

	/** Visual feedback: matches the server, which waters at the lower band edge. */
	const tone = $derived(
		below === null
			? 'bg-brand'
			: value < below
				? 'bg-danger'
				: value <= ideal!
					? 'bg-warn'
					: 'bg-brand'
	);

	function clamp(next: number): number {
		const stepped = Math.round(next / bounds.step) * bounds.step;
		return Math.min(bounds.max, Math.max(bounds.min, stepped));
	}
</script>

<div class="space-y-3">
	<div class="flex items-center gap-3">
		<button
			type="button"
			disabled={disabled || value <= bounds.min}
			onclick={() => (value = clamp(value - bounds.step))}
			class="flex h-11 w-11 shrink-0 cursor-pointer items-center justify-center rounded-xl border border-ink/10 bg-canvas text-xl font-bold text-ink-soft transition active:scale-95 disabled:opacity-40"
			aria-label={$tStore('plants.moistureDecrease')}
		>
			−
		</button>
		<div class="min-w-0 flex-1 text-center">
			<input
				{id}
				type="number"
				{disabled}
				min={bounds.min}
				max={bounds.max}
				step={bounds.step}
				{value}
				oninput={(e) => (value = clamp(Number(e.currentTarget.value)))}
				class="w-full rounded-xl border-2 border-brand/25 bg-canvas py-2 text-center text-2xl font-bold text-ink shadow-sm focus:border-brand focus:ring-2 focus:ring-brand/20 focus:outline-none"
			/>
			<p class="mt-1 text-xs text-ink-soft">
				{$tStore(scale === 'Percent' ? 'plants.moistureUnitPercent' : 'plants.moistureUnitTen')}
			</p>
		</div>
		<button
			type="button"
			disabled={disabled || value >= bounds.max}
			onclick={() => (value = clamp(value + bounds.step))}
			class="flex h-11 w-11 shrink-0 cursor-pointer items-center justify-center rounded-xl border border-ink/10 bg-canvas text-xl font-bold text-ink-soft transition active:scale-95 disabled:opacity-40"
			aria-label={$tStore('plants.moistureIncrease')}
		>
			+
		</button>
	</div>

	<!-- Scale visualization: current value bar plus the target band -->
	<div class="px-1">
		<div class="relative h-3 overflow-hidden rounded-full bg-ink/10">
			<div
				class="h-full rounded-full transition-all duration-300 {tone}"
				style:width="{pct}%"
			></div>
		</div>
		{#if bandLeftPct !== null && bandWidthPct !== null}
			<div class="relative h-0" aria-hidden="true">
				<!-- target band (water at the left edge, refill to the right edge) -->
				<div
					class="absolute -top-3 h-3 rounded-sm bg-brand/35 ring-1 ring-brand-dark/40"
					style:left="{bandLeftPct}%; width: {bandWidthPct}%"
				></div>
				<!-- ideal value (band centre) -->
				<div
					class="absolute -top-5 h-7 w-0.5 -translate-x-1/2 bg-ink/70"
					style:left="{idealPct}%"
				></div>
			</div>
		{/if}
		<div class="mt-2 flex justify-between text-[10px] text-ink-soft">
			<span>{$tStore('plants.moistureDry')}</span>
			{#if below !== null && until !== null}
				<span class="font-semibold text-brand-dark dark:text-brand">
					{$tStore('plants.moistureBandLabel')}: {below}–{until}
				</span>
			{/if}
			<span>{$tStore('plants.moistureWet')}</span>
		</div>
	</div>
</div>
