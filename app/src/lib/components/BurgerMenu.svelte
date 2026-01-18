<script lang="ts">
	import { goto } from '$app/navigation';
	import { IMPRINT_URL, PRIVACY_POLICY_URL, WEBSITE_URL, API_BASE_URL } from '$lib/constants';
	import { openExternalLink } from '$lib/os/browser';
	import { languageStore, setLanguage } from '$lib/stores/language';
	import { tStore } from '$lib/i18n';
	import NotificationDebug from './NotificationDebug.svelte';
	import { resolve } from '$app/paths';

	let isOpen = $state(false);
	let showNotificationDebug = $state(false);

	function toggleMenu() {
		isOpen = !isOpen;
	}

	function closeMenu() {
		isOpen = false;
		showNotificationDebug = false;
	}

	async function handleLogout() {
		// await authStore.logout();
		closeMenu();
	}

	async function handleLanguageChange(lang: 'en' | 'de' | 'es') {
		await setLanguage(lang);
		// Update user language in profile if logged in
		// if ($authStore.user && $authStore.token) {
		// 	try {
		// 		const response = await fetch(API_BASE_URL + '/api/user', {
		// 			method: 'PUT',
		// 			headers: {
		// 				'Content-Type': 'application/json',
		// 				Authorization: `Bearer ${$authStore.token}`
		// 			},
		// 			body: JSON.stringify({
		// 				language: lang
		// 			})
		// 		});

		// 		if (response.ok) {
		// 			const updatedUser = await response.json();
		// 			await authStore.setUser(updatedUser);
		// 		}
		// 	} catch (err) {
		// 		console.error('Failed to update language preference:', err);
		// 	}
		// }
	}

	function navigateTo(path: '/app' | '/app/profile' | '/app/manage') {
		goto(resolve(path));
		closeMenu();
	}

	function toggleNotificationDebug() {
		showNotificationDebug = !showNotificationDebug;
	}
</script>

<div class="relative">
	<!-- Burger Icon Button -->
	<button
		onclick={toggleMenu}
		class="rounded-full border border-emerald-200 bg-white p-3 shadow-lg transition-all hover:scale-105 hover:shadow-xl"
		aria-label="Menu"
	>
		<svg class="h-6 w-6 text-emerald-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M4 6h16M4 12h16M4 18h16"
			></path>
		</svg>
	</button>

	{#if isOpen}
		<!-- Full Screen Menu Overlay -->
		<div class="safe-area fixed inset-0 z-[100] overflow-y-auto bg-white">
			<!-- Close Button -->
			<div class="flex items-center justify-between border-b border-emerald-200 p-6">
				<h2 class="text-2xl font-bold text-emerald-700">Menu</h2>
				<button
					onclick={closeMenu}
					class="rounded-full p-2 transition-colors hover:bg-emerald-100"
					aria-label="Close menu"
				>
					<svg
						class="h-6 w-6 text-emerald-700"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						></path>
					</svg>
				</button>
			</div>

			<!-- User Info -->
			<!-- {#if $authStore.user}
				<div class="border-b border-emerald-200 bg-emerald-50 p-6">
					<div class="flex items-center gap-3">
						<div class="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-600">
							<svg class="h-7 w-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
								></path>
							</svg>
						</div>
						<div>
							<p class="font-semibold text-emerald-900">
								{$authStore.user.username || 'User'}
							</p>
							<p class="text-sm text-emerald-700">{$authStore.user.email}</p>
						</div>
					</div>
				</div>
			{/if} -->

			<!-- Menu Items -->
			<div class="space-y-2 p-6">
				<button
					onclick={() => navigateTo('/app')}
					class="flex w-full items-center gap-3 rounded-xl px-6 py-4 text-left text-lg transition-colors hover:bg-emerald-50"
				>
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
						></path>
					</svg>
					<span class="font-medium text-emerald-900">My Plants</span>
				</button>

				<button
					onclick={() => navigateTo('/app/profile')}
					class="flex w-full items-center gap-3 rounded-xl px-6 py-4 text-left text-lg transition-colors hover:bg-emerald-50"
				>
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
						></path>
					</svg>
					<span class="font-medium text-emerald-900">User Profile</span>
				</button>

				<button
					onclick={() => navigateTo('/app/manage')}
					class="flex w-full items-center gap-3 rounded-xl px-6 py-4 text-left text-lg transition-colors hover:bg-emerald-50"
				>
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 6v6m0 0v6m0-6h6m-6 0H6"
						></path>
					</svg>
					<span class="font-medium text-emerald-900">{$tStore('menu.managePlants')}</span>
				</button>
			</div>

			<!-- Language Picker -->
			<div class="border-t border-emerald-200 px-6 py-4">
				<p class="mb-3 text-sm font-semibold text-emerald-700">
					{$tStore('common.language')}
				</p>
				<div class="flex gap-3">
					<button
						onclick={() => handleLanguageChange('en')}
						class={`rounded-lg px-6 py-3 text-base font-medium transition-colors ${
							$languageStore === 'en'
								? 'bg-emerald-600 text-white'
								: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
						}`}
					>
						EN
					</button>
					<button
						onclick={() => handleLanguageChange('de')}
						class={`rounded-lg px-6 py-3 text-base font-medium transition-colors ${
							$languageStore === 'de'
								? 'bg-emerald-600 text-white'
								: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
						}`}
					>
						DE
					</button>
					<button
						onclick={() => handleLanguageChange('es')}
						class={`rounded-lg px-6 py-3 text-base font-medium transition-colors ${
							$languageStore === 'es'
								? 'bg-emerald-600 text-white'
								: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
						}`}
					>
						ES
					</button>
				</div>
			</div>

			<!-- Notification Debug -->
			<div class="border-t border-emerald-200 px-6 py-4">
				<button
					onclick={toggleNotificationDebug}
					class="flex w-full items-center justify-between rounded-lg px-4 py-3 text-left text-base text-emerald-700 transition-colors hover:bg-emerald-50"
				>
					<span>🔔 Push Notifications (Debug)</span>
					<svg
						class="h-4 w-4 transition-transform {showNotificationDebug ? 'rotate-180' : ''}"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"
						></path>
					</svg>
				</button>
				{#if showNotificationDebug}
					<div class="mt-3">
						<NotificationDebug />
					</div>
				{/if}
			</div>

			<!-- Footer Links -->
			<div class="space-y-2 border-t border-emerald-200 px-6 py-4">
				<button
					class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
					onclick={() => openExternalLink(WEBSITE_URL)}
				>
					Website
				</button>
				<button
					class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
					onclick={() => openExternalLink(PRIVACY_POLICY_URL)}
				>
					Privacy Policy
				</button>
				<button
					class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
					onclick={() => openExternalLink(IMPRINT_URL)}
				>
					Impressum
				</button>
			</div>

			<!-- Build Info -->
			<div class="border-t border-emerald-200 px-6 py-4 text-sm text-emerald-600">
				<p>Build: {new Date().toISOString().split('T')[0]}</p>
				<p>Version: 1.0.0</p>
			</div>

			<!-- Logout Button -->
			<div class="border-t border-emerald-200 px-6 py-4">
				<button
					onclick={handleLogout}
					class="w-full rounded-lg bg-red-50 px-6 py-4 text-center text-lg font-semibold text-red-600 transition-colors hover:bg-red-100"
				>
					Logout
				</button>
			</div>
		</div>
	{/if}
</div>

<style>
	.safe-area {
		padding-top: env(safe-area-inset-top);
		padding-bottom: env(safe-area-inset-bottom);
		padding-left: env(safe-area-inset-left);
		padding-right: env(safe-area-inset-right);
	}
</style>
