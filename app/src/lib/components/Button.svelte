<script lang="ts">
	import { tStore } from '$lib/i18n';

	export let disabled: boolean = false;
	export let onclick: () => void = () => {};
	export let loading: boolean = false;
	export let text: string = 'Button';
	export let loadingText: string = 'Loading...';
</script>

<button
	{onclick}
	{disabled}
	class="w-full cursor-pointer rounded-lg bg-green-600 px-4 py-3 font-semibold text-white transition hover:bg-green-700 focus:ring-2 focus:ring-green-500 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
>
	{#if loading}
		<span class="inline-flex items-center gap-2">
			<svg
				class="h-4 w-4 animate-spin"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
			>
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path
					class="opacity-75"
					fill="currentColor"
					d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
				></path>
			</svg>
			{$tStore(loadingText)}
		</span>
	{:else}
		{$tStore(text)}
	{/if}
</button>
