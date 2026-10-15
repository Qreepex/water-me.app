<script lang="ts">
	import { tStore } from '$lib/i18n';
	import { MoistureScale } from '$lib/types/api';

	interface Props {
		value: number;
		scale?: MoistureScale;
		/** Value considered "too dry" (draws a marker when provided) */
		ideal?: number | null;
		disabled?: boolean;
		id?: string;
	}

	let {
		value = $bindable(),
		scale = MoistureScale.Scale1To10,
		ideal = null,
		disabled = false,
		id
	}: Props = $props();

	const bounds = $derived(
		scale === MoistureScale.Percent ? { min: 0, max: 100, step: 1 } : { min: 1, max: 10, step: 0.5 }
	);

	const pct = $derived(((value - bounds.min) / (bounds.max - bounds.min)) * 100);
	const idealPct = $derived(
		ideal === null || ideal === undefined
			? null
			: ((ideal - bounds.min) / (bounds.max - bounds.min)) * 100
	);

	/** Visual dryness feedback: matches the server's "water now" threshold. */
	const tone = $derived(
		ideal !== null && ideal !== undefined && value <= ideal
			? 'bg-danger'
			: ideal !== null && ideal !== undefined && value <= ideal + bounds.step * 2
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

	<!-- Scale visualization with the ideal marker -->
	<div class="px-1">
		<div class="relative h-3 overflow-hidden rounded-full bg-ink/10">
			<div
				class="h-full rounded-full transition-all duration-300 {tone}"
				style:width="{pct}%"
			></div>
		</div>
		{#if idealPct !== null}
			<div
				class="relative h-0"
				title={$tStore('plants.moistureIdealMarker')}
				style:left="{idealPct}%"
				aria-hidden="true"
			>
				<div class="absolute -top-4 h-4 w-0.5 -translate-x-1/2 bg-ink/70"></div>
			</div>
		{/if}
		<div class="mt-1.5 flex justify-between text-[10px] text-ink-soft">
			<span>{$tStore('plants.moistureDry')}</span>
			{#if ideal !== null && ideal !== undefined}
				<span class="font-semibold text-brand-dark dark:text-brand">
					{$tStore('plants.moistureIdealMarker')}: {ideal}
				</span>
			{/if}
			<span>{$tStore('plants.moistureWet')}</span>
		</div>
	</div>
</div>
