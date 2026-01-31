<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		icon: string;
		title: string;
		description: string;
		children?: Snippet;
	}

	const { icon, title, description, children }: Props = $props();
</script>

<div
	class="rounded-2xl border-2 border-dashed border-[var(--p-emerald)] bg-[var(--bg-light)] p-12 text-center"
>
	<div class="mb-4 text-5xl">{icon}</div>
	<h2 class="mb-2 text-2xl font-bold text-[var(--text-light-main)]">{title}</h2>
	<p class="mb-6 text-[var(--text-light-main)]/60">{description}</p>
	{#if children}
		<div class="mt-6">
			{@render children?.()}
		</div>
	{/if}
</div>
