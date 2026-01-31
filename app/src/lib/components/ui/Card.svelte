<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		rounded?: 'md' | 'lg' | 'xl' | '2xl';
		shadow?: 'sm' | 'md' | 'lg';
		children: Snippet;
	}

	const { rounded = 'lg', shadow = 'md', children }: Props = $props();

	const roundedClasses = {
		md: 'rounded-md',
		lg: 'rounded-lg',
		xl: 'rounded-xl',
		'2xl': 'rounded-2xl'
	};
	const shadowClasses = { sm: 'shadow-sm', md: 'shadow-md', lg: 'shadow-lg' };
</script>

<div
	class="overflow-hidden border border-[var(--p-emerald)]/30 {roundedClasses[
		rounded
	]} bg-[var(--card-light)] backdrop-blur transition hover:shadow-lg {shadowClasses[shadow]}"
>
	{@render children()}
</div>
