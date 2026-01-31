<script lang="ts">
	import ButtonSpinner from '$lib/assets/ButtonSpinner.svelte';
	import { tStore } from '$lib/i18n';

	type Variant = 'primary' | 'secondary' | 'danger' | 'ghost';
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
		icon
	}: Props = $props();

	const variantClasses: Record<Variant, string> = {
		primary:
			'bg-[var(--p-emerald)] text-[var(--text-light-main)] hover:bg-[var(--p-emerald-dark)] font-semibold shadow-[0_4px_14px_rgba(0,238,87,0.4)]',
		secondary: 'bg-[var(--status-success)] text-white hover:opacity-90 font-medium',
		danger: 'bg-[var(--status-error)] text-white hover:opacity-90 font-medium',
		ghost:
			'bg-transparent border-2 border-[var(--p-emerald)] text-[var(--text-light-main)] hover:bg-[var(--bg-light)] font-medium'
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
	{:else if icon}
		<span class="inline-flex items-center gap-2">
			<span>{icon}</span>
			<span>{$tStore(text)}</span>
		</span>
	{:else}
		{$tStore(text)}
	{/if}
</button>
