<script lang="ts">
	import { IMPRINT_URL, PRIVACY_POLICY_URL, WEBSITE_URL } from '$lib/constants';
	import { openExternalLink } from '$lib/os/browser';
	import { languageStore, setLanguage } from '$lib/stores/language';
	import { tStore } from '$lib/i18n';
	const { onClose } = $props<{ onClose?: () => void }>();

	async function handleLogout() {}

	async function handleLanguageChange(lang: 'en' | 'de' | 'es') {
		await setLanguage(lang);
	}
</script>

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

<!-- Footer Links -->
<div class="space-y-2 border-t border-emerald-200 px-6 py-4">
	<button
		class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
		onclick={() => openExternalLink(WEBSITE_URL)}
	>
		{$tStore('menu.website')}
	</button>
	<button
		class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
		onclick={() => openExternalLink(PRIVACY_POLICY_URL)}
	>
		{$tStore('menu.privacyPolicy')}
	</button>
	<button
		class="block rounded-lg px-4 py-3 text-base text-emerald-700 transition-colors hover:bg-emerald-50"
		onclick={() => openExternalLink(IMPRINT_URL)}
	>
		{$tStore('menu.imprint')}
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
