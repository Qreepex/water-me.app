<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { languageStore } from '$lib/i18n.svelte';
	import { page } from '$app/stores';

	let { children } = $props();

	$effect(() => {
		// Update language based on URL
		const path = $page.url.pathname;
		if (path.startsWith('/de')) {
			languageStore.setLanguage('de');
		} else if (path.startsWith('/es')) {
			languageStore.setLanguage('es');
		} else {
			languageStore.setLanguage('en');
		}
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

{@render children()}
