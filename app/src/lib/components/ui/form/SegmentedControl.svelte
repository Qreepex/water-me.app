<script lang="ts" generics="T extends string">
	import { tStore } from '$lib/i18n';

	interface Props {
		value: T;
		/** i18n key for the accessible group label */
		label: string;
		options: { value: T; label: string; description?: string; icon?: string }[];
	}

	let { value = $bindable(), label, options }: Props = $props();
</script>

<div role="radiogroup" aria-label={$tStore(label)} class="grid gap-2 sm:grid-cols-2">
	{#each options as option (option.value)}
		{@const active = value === option.value}
		<button
			type="button"
			role="radio"
			aria-checked={active}
			onclick={() => (value = option.value)}
			class="flex min-h-11 cursor-pointer flex-col items-start gap-0.5 rounded-xl border-2 px-4 py-3 text-left transition active:scale-[0.98] {active
				? 'border-brand bg-brand/15 shadow-[0_4px_14px_rgba(0,238,87,0.25)]'
				: 'border-ink/10 bg-canvas hover:border-brand/40'}"
		>
			<span class="text-sm font-semibold {active ? 'text-ink' : 'text-ink-soft'}">
				{$tStore(option.label)}
			</span>
			{#if option.description}
				<span class="text-xs text-ink-soft">{$tStore(option.description)}</span>
			{/if}
		</button>
	{/each}
</div>
