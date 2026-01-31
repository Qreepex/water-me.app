<script lang="ts">
	interface Props {
		value?: string;
		placeholder?: string;
		ariaLabel?: string;
	}

	let { value = $bindable(''), placeholder = '', ariaLabel = 'Search' }: Props = $props();

	function clear() {
		value = '';
	}
</script>

<div class="relative flex-1">
	<input
		type="text"
		bind:value
		{placeholder}
		aria-label={ariaLabel}
		class="w-full rounded-lg border border-[var(--p-emerald)]/30 bg-white py-2 pr-10 pl-10 text-sm placeholder-[var(--text-light-main)]/40 transition-all focus:border-[var(--p-emerald)]/60 focus:ring-2 focus:ring-[var(--p-emerald)]/20 focus:outline-none"
	/>
	<svg
		class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-[var(--p-emerald)]/50"
		fill="none"
		stroke="currentColor"
		viewBox="0 0 24 24"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="2"
			d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
		/>
	</svg>
	{#if value}
		<button
			type="button"
			onclick={clear}
			class="absolute top-1/2 right-2 -translate-y-1/2 rounded-full p-1 text-[var(--text-light-main)]/60 transition hover:bg-[var(--p-emerald)]/10 hover:text-[var(--p-emerald-dark)]"
			aria-label="Clear search"
		>
			<svg viewBox="0 0 20 20" class="h-4 w-4" fill="currentColor">
				<path
					d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
				/>
			</svg>
		</button>
	{/if}
</div>
