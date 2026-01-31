<script lang="ts">
	import { IMPRINT_URL, PRIVACY_POLICY_URL, WEBSITE_URL } from '$lib/constants';
	import { openExternalLink } from '$lib/os/browser';
	import { languageStore, setLanguage } from '$lib/stores/language';
	import { tStore } from '$lib/i18n';
	import NotificationDebug from './NotificationDebug.svelte';
	import { resolve } from '$app/paths';
	import Burger from '$lib/assets/Burger.svg.svelte';
	import Bell from '$lib/assets/Bell.svg.svelte';

	const { onClose } = $props<{ onClose?: () => void }>();

	let showNotificationDebug = $state(false);

	async function handleLogout() {
		// await authStore.logout();
		onClose?.();
	}

	async function handleLanguageChange(lang: 'en' | 'de' | 'es') {
		await setLanguage(lang);
	}

	function toggleNotificationDebug() {
		showNotificationDebug = !showNotificationDebug;
	}
</script>

<!-- Settings Content -->
<div class="space-y-2 p-6">
	<!-- Profile Section -->
	<a
		href={resolve('/profile')}
		onclick={onClose}
		class="flex w-full items-center gap-3 rounded-xl px-6 py-4 text-left text-lg transition-colors hover:bg-emerald-50"
	>
		<Burger isActive={false} />
		<span class="font-medium text-emerald-900">User Profile</span>
	</a>
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
		<Bell />
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
<div class="border-t border-emerald-200 px-6 py-4 pb-24">
	<button
		onclick={handleLogout}
		class="w-full rounded-lg bg-red-50 px-6 py-4 text-center text-lg font-semibold text-red-600 transition-colors hover:bg-red-100"
	>
		Logout
	</button>
</div>
