<script lang="ts">
	import ButtonSpinner from '$lib/assets/ButtonSpinner.svelte';
	import { tStore } from '$lib/i18n';
	import type { Snippet } from 'svelte';

	type Variant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'water';
	type Size = 'sm' | 'md' | 'lg';

	interface Props {
		variant?: Variant;
		size?: Size;
		disabled?: boolean;
		loading?: boolean;
		onclick?: () => void;
		loadingText?: string;
		text: string;
		class?: string;
		icon?: string;
		iconComponent?: Snippet;
	}

	const {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		onclick,
		loadingText = 'loading',
		text,
		class: className = '',
		icon,
		iconComponent
	}: Props = $props();

	const Icon = $derived(iconComponent);

	const variantClasses: Record<Variant, string> = {
		primary:
			'bg-[var(--p-emerald)] text-[var(--text-light-main)] hover:bg-[var(--p-emerald-dark)] font-semibold shadow-[0_4px_14px_rgba(0,238,87,0.4)]',
		secondary: 'bg-[var(--status-success)] text-white hover:opacity-90 font-medium',
		danger: 'bg-[var(--status-error)] text-white hover:opacity-90 font-medium',
		ghost:
			'bg-transparent border-2 border-[var(--p-emerald)] text-[var(--text-light-main)] hover:bg-[var(--bg-light)] font-medium',
		water:
			'bg-[var(--status-info)] text-white hover:opacity-90 font-semibold shadow-[0_4px_14px_rgba(33,158,188,0.4)]'
	};

	const sizeClasses: Record<Size, string> = {
		sm: 'px-4 py-2 text-sm rounded-lg',
		md: 'px-6 py-3 text-base rounded-lg',
		lg: 'px-8 py-4 text-lg rounded-lg'
	};
</script>

<button
	{disabled}
	{onclick}
	class="cursor-pointer transition disabled:opacity-50 {variantClasses[variant]} {sizeClasses[
		size
	]} {className}"
>
	{#if loading}
		<span class="inline-flex items-center gap-2">
			<ButtonSpinner />
			{$tStore(loadingText)}
		</span>
	{:else if icon || iconComponent}
		<span class="inline-flex items-center gap-2">
			{#if iconComponent}
				<Icon />
			{:else if icon}
				<span>{icon}</span>
			{/if}
			<span>{$tStore(text)}</span>
		</span>
	{:else}
		{$tStore(text)}
	{/if}
</button>
