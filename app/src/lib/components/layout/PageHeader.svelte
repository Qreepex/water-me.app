<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		icon?: string;
		children?: Snippet;
	}

	const { title, description, icon, children }: Props = $props();
</script>

<div class="mb-8">
	<div class="mb-4 flex items-center justify-between gap-3">
		<div class="flex items-center gap-3">
			{#if icon}
				<span class="text-5xl">{icon}</span>
			{/if}
			<div>
				<h1 class="text-4xl font-bold text-[var(--text-light-main)]">{title}</h1>
				{#if description}
					<p class="mt-1 text-[var(--text-light-main)]/60">{description}</p>
				{/if}
			</div>
		</div>

		{@render children?.()}
	</div>
</div>
