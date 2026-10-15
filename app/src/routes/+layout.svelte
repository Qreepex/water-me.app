<script lang="ts">
	import { onMount } from 'svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Auth from '$lib/auth/Auth.svelte';
	import BottomNav from '$lib/components/layout/BottomNav.svelte';
	import { initializeI18n } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { initializeLanguage } from '$lib/stores/language';
	import { onAppResume } from '$lib/utils/lifecycle';
	import { Capacitor } from '@capacitor/core';
	import { SplashScreen } from '@capacitor/splash-screen';

	let { children } = $props();

	const plantsStore = getPlantsStore();

	onMount(() => {
		// Plants are loaded once per session, so changes made on another device
		// (e.g. a moisture measurement) would otherwise stay invisible until a
		// reload. Refresh whenever the user comes back to the app.
		const stopResumeListener = onAppResume(() => {
			if (plantsStore.hasLoaded) void plantsStore.syncFromServer();
		});

		void (async () => {
			// Hide splash screen once the app is ready
			if (Capacitor.isNativePlatform()) {
				await SplashScreen.hide();
			}

			// Initialize language from user profile or preferences
			await initializeLanguage();

			// Initialize i18n translations for the selected language
			await initializeI18n();

			// Register service worker for image caching (web only)
			if (typeof navigator !== 'undefined' && 'serviceWorker' in navigator) {
				try {
					await navigator.serviceWorker.register('/sw.js');
				} catch {
					// ignore
				}
			}

			// Do not auto-request notification permissions on startup.
			// Use $lib/notifications.requestNotificationPermissions() when user opts in.
		})();

		return stopResumeListener;
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="fixed inset-0 bg-canvas">
	<!-- Brand gradient wash behind the top of every page (DESIGN.md header gradient) -->
	<div
		aria-hidden="true"
		class="pointer-events-none absolute inset-x-0 top-0 h-80"
		style="background: var(--header-wash);"
	></div>
	<main class="pt-safe relative flex h-full flex-col overflow-hidden">
		<Auth>
			<div class="flex-1 overflow-hidden pt-2 md:px-10 md:pt-10 xl:px-32 xl:pt-14">
				<BottomNav />
				<div class="flex h-full flex-col pb-20">
					{@render children()}
				</div>
			</div>
		</Auth>
	</main>
</div>

<style>
	/* Safe area insets for mobile notches and status bars */
	.pt-safe {
		padding-top: env(safe-area-inset-top);
	}
</style>
