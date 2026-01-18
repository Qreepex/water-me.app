<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { initializePushNotifications, cleanupPushNotifications } from '$lib/notifications';
	import { initializeI18n } from '$lib/i18n';
	import { initializeLanguage } from '$lib/stores/language';
	import { browser } from '$app/environment';
	import Auth from '$lib/auth/Auth.svelte';
	import { Capacitor } from '@capacitor/core';
	import { SplashScreen } from '@capacitor/splash-screen';

	let { children } = $props();
	let fcmToken = $state<string | null>(null);

	onMount(async () => {
		// hide splash screen once the app is ready
		if (Capacitor.isNativePlatform()) {
			await SplashScreen.hide();
		}

		// Initialize language from user profile or preferences
		await initializeLanguage();

		// Initialize i18n translations for the selected language
		await initializeI18n();

		// Initialize push notifications
		const result = await initializePushNotifications();
		fcmToken = result.token;

		if (fcmToken) {
			console.log('✅ FCM Token registered:', fcmToken);
		}
	});

	onDestroy(() => {
		if (browser) {
			cleanupPushNotifications();
		}
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="relative min-h-screen bg-gradient-to-br from-emerald-50 to-green-50">
	<main class="pt-safe pb-safe px-4">
		<Auth>
			{@render children()}
		</Auth>
	</main>
</div>

<style>
	/* Safe area insets for mobile notches and status bars */
	.pt-safe {
		padding-top: env(safe-area-inset-top);
	}
	.pb-safe {
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
