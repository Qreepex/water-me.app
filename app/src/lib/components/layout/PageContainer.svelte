<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		maxWidth?: 'sm' | 'md' | 'lg' | 'xl';
		noPadding?: boolean;
		gradient?: boolean;
		children: Snippet;
	}

	const { maxWidth = 'lg', noPadding = false, gradient = true, children }: Props = $props();

	const widthClasses = {
		sm: 'max-w-2xl',
		md: 'max-w-4xl',
		lg: 'max-w-6xl',
		xl: 'max-w-7xl'
	};
</script>

<div
	class="min-h-screen {gradient ? 'bg-[var(--bg-light)]' : 'bg-[var(--card-light)]'} {!noPadding &&
		'p-6 md:p-10'}"
>
	<div class="mx-auto {widthClasses[maxWidth]} {!noPadding && 'pb-24'}">
		{@render children()}
	</div>
</div>
