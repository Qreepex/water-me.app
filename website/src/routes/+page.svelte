<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';

	onMount(() => {
		// Detect user's preferred language
		let userLang = 'en';

		if (typeof window !== 'undefined') {
			const browserLang = navigator.language.split('-')[0];
			if (browserLang === 'de') {
				userLang = 'de';
			} else if (browserLang === 'es') {
				userLang = 'es';
			}
		}

		// Redirect to language-specific page
		goto(resolve(`/${userLang}`), { replaceState: true });
	});
</script>

<svelte:head>
	<title>Water Me - Redirecting...</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-[#f2faf5] via-white to-[#e8f9f0]"
>
	<div class="text-center">
		<div class="mb-4 text-6xl">🌿</div>
		<h1 class="text-2xl font-bold text-[#061f12]">Redirecting...</h1>
	</div>
</div>
