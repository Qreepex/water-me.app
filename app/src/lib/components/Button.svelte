<script lang="ts">
	import ButtonSpinner from '$lib/assets/ButtonSpinner.svelte';
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
			<ButtonSpinner />
			{$tStore(loadingText)}
		</span>
	{:else}
		{$tStore(text)}
	{/if}
</button>
