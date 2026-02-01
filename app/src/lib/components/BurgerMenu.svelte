<script lang="ts">
	import { IMPRINT_URL, PRIVACY_POLICY_URL, WEBSITE_URL } from '$lib/constants';
	import { openExternalLink } from '$lib/os/browser';
	import { languageStore, setLanguage } from '$lib/stores/language';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import { imageCacheStore } from '$lib/stores/imageCache.svelte';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { invalidateApiCache } from '$lib/utils/cache';
	import PageHeader from './layout/PageHeader.svelte';
	import PageContent from './layout/PageContent.svelte';
	import Card from './ui/Card.svelte';
	import Button from './ui/Button.svelte';

	const plantsStore = getPlantsStore();
	let isRefreshing = $state(false);

	async function handleLanguageChange(lang: 'en' | 'de' | 'es') {
		await setLanguage(lang);
	}

	async function handleRefresh() {
		isRefreshing = true;
		try {
			plantsStore.setLoading(true);
			plantsStore.setError(null);

			// Clear image cache
			imageCacheStore.cleanup(0); // Force cleanup with 0 max age

			await invalidateApiCache(['/api/plants']);

			const result = await fetchData('/api/plants', {});
			if (!result.ok) {
				plantsStore.setError(result.error?.message || 'Failed to fetch plants');
				return;
			}

			plantsStore.setPlants(result.data || []);
		} finally {
			plantsStore.setLoading(false);
			isRefreshing = false;
		}
	}

	const buildDate = new Date().toLocaleDateString();
	const version = '1.0.0';
</script>

<div class="flex-shrink-0">
	<PageHeader icon="⚙️" title="menu.settings" />
</div>

<PageContent>
	<div class="overflow-y-auto">
		<div class="space-y-6 py-6">
			<!-- Refresh Button Section -->
			<Card rounded="2xl">
				<div class="p-6">
					<Button
						onclick={handleRefresh}
						disabled={isRefreshing}
						class="w-full"
						text="common.refresh"
					></Button>
				</div>
			</Card>

			<!-- Language Picker -->
			<Card rounded="2xl">
				<div class="p-6">
					<p class="mb-4 text-sm font-semibold text-[var(--text-light-main)]">
						{$tStore('common.language')}
					</p>
					<div class="flex gap-3">
						<Button
							onclick={() => handleLanguageChange('en')}
							class={`flex-1 ${
								$languageStore === 'en'
									? 'bg-[var(--p-emerald)] text-white'
									: 'bg-[var(--p-emerald)]/10 text-[var(--p-emerald-dark)] hover:bg-[var(--p-emerald)]/20'
							}`}
							text="🇬🇧 EN"
						></Button>
						<Button
							onclick={() => handleLanguageChange('de')}
							class={`flex-1 ${
								$languageStore === 'de'
									? 'bg-[var(--p-emerald)] text-white'
									: 'bg-[var(--p-emerald)]/10 text-[var(--p-emerald-dark)] hover:bg-[var(--p-emerald)]/20'
							}`}
							text="🇩🇪 DE"
						></Button>
						<Button
							onclick={() => handleLanguageChange('es')}
							class={`flex-1 ${
								$languageStore === 'es'
									? 'bg-[var(--p-emerald)] text-white'
									: 'bg-[var(--p-emerald)]/10 text-[var(--p-emerald-dark)] hover:bg-[var(--p-emerald)]/20'
							}`}
							text="🇪🇸 ES"
						></Button>
					</div>
				</div>
			</Card>

			<!-- Resources -->
			<Card rounded="2xl">
				<div class="p-6">
					<p class="mb-4 text-sm font-semibold text-[var(--text-light-main)]">
						{$tStore('menu.resources')}
					</p>
					<div class="space-y-2">
						<button
							class="w-full rounded-lg px-4 py-3 text-left text-sm text-[var(--p-emerald-dark)] transition-colors hover:bg-[var(--p-emerald)]/10"
							onclick={() => openExternalLink(WEBSITE_URL)}
						>
							🌐 {$tStore('menu.website')}
						</button>
						<button
							class="w-full rounded-lg px-4 py-3 text-left text-sm text-[var(--p-emerald-dark)] transition-colors hover:bg-[var(--p-emerald)]/10"
							onclick={() => openExternalLink(PRIVACY_POLICY_URL)}
						>
							🔒 {$tStore('menu.privacyPolicy')}
						</button>
						<button
							class="w-full rounded-lg px-4 py-3 text-left text-sm text-[var(--p-emerald-dark)] transition-colors hover:bg-[var(--p-emerald)]/10"
							onclick={() => openExternalLink(IMPRINT_URL)}
						>
							ℹ️ {$tStore('menu.imprint')}
						</button>
					</div>
				</div>
			</Card>

			<!-- Build Info -->
			<Card rounded="2xl">
				<div class="p-6">
					<p class="mb-3 text-sm font-semibold text-[var(--text-light-main)]">
						{$tStore('menu.about')}
					</p>
					<div class="space-y-2 text-sm text-[var(--text-light-main)]/70">
						<div class="flex items-center justify-between">
							<span>Version</span>
							<span class="font-medium text-[var(--text-light-main)]">{version}</span>
						</div>
						<div class="flex items-center justify-between">
							<span>Build Date</span>
							<span class="font-medium text-[var(--text-light-main)]">{buildDate}</span>
						</div>
					</div>
				</div>
			</Card>
		</div>
	</div>
</PageContent>
