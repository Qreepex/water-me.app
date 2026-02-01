<script lang="ts">
	import { resolve } from '$app/paths';
	import { GITHUB_URL, INSTAGRAM_URL, TIKTOK_URL, WEB_APP_URL } from '$lib';
	import { tStore, languageStore } from '$lib/i18n.svelte';
	import { derived } from 'svelte/store';

	// Create reactive translations
	const t = derived(languageStore, () => {
		return (key: string) => tStore(key);
	});

	let displayStats = $state({
		users: 0,
		plants: 0,
		reminders: 0
	});
	let hasAnimated = $state(false);

	$effect(() => {
		if (hasAnimated) return;
		hasAnimated = true;

		const targetStats = {
			users: 2850,
			plants: 45320,
			reminders: 128950
		};

		const duration = 2000;
		const startTime = Date.now();

		const animate = () => {
			const elapsed = Date.now() - startTime;
			const progress = Math.min(elapsed / duration, 1);

			displayStats.users = Math.floor(targetStats.users * progress);
			displayStats.plants = Math.floor(targetStats.plants * progress);
			displayStats.reminders = Math.floor(targetStats.reminders * progress);

			if (progress < 1) {
				requestAnimationFrame(animate);
			}
		};

		animate();
	});
</script>

<svelte:head>
	<title>Water Me - Never Miss a Plant Care Reminder</title>
	<meta
		name="description"
		content="Smart plant care companion app. Get reminders for watering, spraying, and fertilizing. Track your plant collection with ease."
	/>
	<meta name="keywords" content="plant care, app, reminders, watering, gardening, plants" />
	<meta property="og:title" content="Water Me - Smart Plant Care Companion" />
	<meta
		property="og:description"
		content="Never forget to water your plants again. Track watering schedules, get smart reminders, and keep your plants thriving."
	/>
	<meta property="og:type" content="website" />
	<meta property="og:image" content="/og-image.jpg" />
	<link rel="canonical" href="https://water-me.app/de" />
	<link rel="alternate" hreflang="en" href="https://water-me.app/en" />
	<link rel="alternate" hreflang="de" href="https://water-me.app/de" />
	<link rel="alternate" hreflang="es" href="https://water-me.app/es" />
	<link rel="alternate" hreflang="x-default" href="https://water-me.app/en" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-[#f2faf5] via-white to-[#e8f9f0]">
	<!-- Navigation -->
	<nav class="sticky top-0 z-50 border-b border-[#00ee57]/10 bg-white/80 backdrop-blur-md">
		<div class="mx-auto max-w-6xl px-4 py-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div class="text-3xl">🌿</div>
					<span class="text-2xl font-bold text-[#061f12]">Water Me</span>
				</div>
				<div class="flex gap-4">
					<div class="hidden gap-8 sm:flex">
						<a href="#features" class="text-[#4b6658] transition hover:text-[#00ee57]">
							{$t('nav.features')}
						</a>
						<a href="#stats" class="text-[#4b6658] transition hover:text-[#00ee57]">
							{$t('nav.stats')}
						</a>
						<a href="#download" class="text-[#4b6658] transition hover:text-[#00ee57]">
							{$t('nav.download')}
						</a>
						<a
							href={WEB_APP_URL}
							target="_blank"
							rel="noopener noreferrer external"
							class="inline-flex items-center gap-2 rounded-full bg-[#00ee57] px-6 py-2 font-semibold text-[#061f12] transition hover:bg-[#00a343] hover:text-white"
						>
							{$t('nav.open_app')}
						</a>
					</div>
					<select
						bind:value={$languageStore}
						onchange={() => {
							const lang = $languageStore;
							window.location.href = resolve(`/${lang}`);
						}}
						class="rounded border border-[#00ee57]/20 bg-white px-2 py-1 text-[#061f12]"
					>
						<option value="en">English</option>
						<option value="de">Deutsch</option>
						<option value="es">Español</option>
					</select>
				</div>
			</div>
		</div>
	</nav>
	<!-- Hero Section -->
	<section class="relative overflow-hidden px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<div class="grid items-center gap-12 lg:grid-cols-2">
				<!-- Left Content -->
				<div>
					<h1 class="text-5xl leading-tight font-bold text-[#061f12] sm:text-6xl">
						{$t('hero.title')}
					</h1>
					<p class="mt-6 text-xl text-[#4b6658]">
						{$t('hero.subtitle')}
					</p>

					<!-- CTA Buttons -->
					<div class="mt-8 flex gap-4">
						<a
							href="#download"
							class="inline-flex items-center gap-2 rounded-full bg-[#00ee57] px-8 py-4 font-semibold text-[#061f12] shadow-lg shadow-[#00ee57]/40 transition hover:shadow-xl hover:shadow-[#00ee57]/50"
						>
							📥 {$t('hero.cta_primary')}
						</a>
						<a
							href={WEB_APP_URL}
							target="_blank"
							rel="noopener noreferrer external"
							class="inline-flex items-center gap-2 rounded-full border-2 border-[#00ee57] px-8 py-4 font-semibold text-[#00ee57] transition hover:bg-[#00ee57]/10"
						>
							🌐 {$t('hero.cta_secondary')}
						</a>
					</div>

					<!-- Trust Badges -->
					<div class="mt-12 flex items-center gap-4">
						<div class="text-4xl">⭐⭐⭐⭐⭐</div>
						<div>
							<p class="font-semibold text-[#061f12]">{$t('hero.trust_badge')}</p>
						</div>
					</div>
				</div>

				<!-- Right Visual -->
				<div class="relative h-96 lg:h-full">
					<div
						class="absolute inset-0 rounded-3xl bg-gradient-to-br from-[#00ee57] to-[#00a343] shadow-2xl"
					>
						<div class="flex h-full items-center justify-center text-8xl">🌱</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="bg-white/50 px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<div class="mb-16 text-center">
				<h2 class="text-4xl font-bold text-[#061f12] sm:text-5xl">Features</h2>
				<p class="mt-4 text-xl text-[#4b6658]">Smart plant care companion</p>
			</div>

			<div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
				<!-- Feature 1 -->
				<div
					class="rounded-2xl border border-[#00ee57]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">💧</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">
						{$t('features.watering.title')}
					</h3>
					<p class="text-[#4b6658]">
						{$t('features.watering.description')}
					</p>
				</div>

				<!-- Feature 2 -->
				<div
					class="rounded-2xl border border-[#ffb703]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">🎯</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">
						{$t('features.multicare.title')}
					</h3>
					<p class="text-[#4b6658]">
						{$t('features.multicare.description')}
					</p>
				</div>

				<!-- Feature 3 -->
				<div
					class="rounded-2xl border border-[#219ebc]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">📸</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">{$t('features.growth.title')}</h3>
					<p class="text-[#4b6658]">
						{$t('features.growth.description')}
					</p>
				</div>

				<!-- Feature 4 -->
				<div
					class="rounded-2xl border border-[#e63946]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">📊</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">{$t('features.history.title')}</h3>
					<p class="text-[#4b6658]">
						{$t('features.history.description')}
					</p>
				</div>

				<!-- Feature 5 -->
				<div
					class="rounded-2xl border border-[#008f41]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">🔔</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">
						{$t('features.notifications.title')}
					</h3>
					<p class="text-[#4b6658]">
						{$t('features.notifications.description')}
					</p>
				</div>

				<!-- Feature 6 -->
				<div
					class="rounded-2xl border border-[#00ee57]/20 bg-white p-8 shadow-sm transition hover:shadow-lg"
				>
					<div class="mb-4 text-5xl">🌍</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">
						{$t('features.multilang.title')}
					</h3>
					<p class="text-[#4b6658]">
						{$t('features.multilang.description')}
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Stats Section -->
	<section
		id="stats"
		class="bg-gradient-to-r from-[#00ee57] to-[#00a343] px-4 py-20 sm:px-6 lg:px-8"
	>
		<div class="mx-auto max-w-6xl">
			<div class="mb-12 text-center">
				<h2 class="text-4xl font-bold text-white sm:text-5xl">{$t('stats.title')}</h2>
			</div>

			<div class="grid gap-8 sm:grid-cols-3">
				<!-- Stat 1 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">{displayStats.users.toLocaleString()}</div>
					<p class="mt-2 text-lg text-white/90">{$t('stats.users')}</p>
				</div>

				<!-- Stat 2 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">{displayStats.plants.toLocaleString()}</div>
					<p class="mt-2 text-lg text-white/90">{$t('stats.plants')}</p>
				</div>

				<!-- Stat 3 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">
						{displayStats.reminders.toLocaleString()}+
					</div>
					<p class="mt-2 text-lg text-white/90">{$t('stats.reminders')}</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Screenshots/Demo Section -->
	<section class="px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<h2 class="text-center text-4xl font-bold text-[#061f12] sm:text-5xl">
				Beautiful & Intuitive Design
			</h2>
			<p class="mx-auto mt-4 max-w-2xl text-center text-lg text-[#4b6658]">
				Crafted with care for the best user experience. Easy to use, powerful features.
			</p>

			<div class="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
				{#each [{ emoji: '🏠', title: 'Home', desc: 'Quick overview of all your plants' }, { emoji: '💧', title: 'Watering', desc: 'Never miss a watering schedule' }, { emoji: '✏️', title: 'Management', desc: 'Full control over plant details' }] as item (item.title)}
					<div class="overflow-hidden rounded-2xl border border-[#00ee57]/20 bg-white shadow-sm">
						<div
							class="flex h-48 items-center justify-center bg-gradient-to-br from-[#00ee57]/10 to-[#00a343]/10 text-6xl"
						>
							{item.emoji}
						</div>
						<div class="p-6">
							<h3 class="text-xl font-semibold text-[#061f12]">{item.title}</h3>
							<p class="mt-2 text-[#4b6658]">{item.desc}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Download Section -->
	<section id="download" class="bg-white/50 px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-4xl">
			<div
				class="rounded-3xl border-2 border-[#00ee57] bg-gradient-to-br from-[#00ee57]/5 to-[#00a343]/5 p-12 text-center"
			>
				<h2 class="text-4xl font-bold text-[#061f12] sm:text-5xl">{$t('download.title')}</h2>
				<p class="mt-4 text-lg text-[#4b6658]">
					{$t('download.subtitle')}
				</p>

				<div class="mt-8 flex flex-col gap-4 sm:flex-row sm:justify-center">
					<a
						href="https://play.google.com/store/apps/details?id=app.waterme.app"
						class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#061f12] px-8 py-4 font-semibold text-white transition hover:bg-[#061f12]/80"
					>
						<span>🤖</span>
						{$t('download.googleplay')}
					</a>
					<div
						class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#061f12] px-8 py-4 font-semibold text-white transition hover:bg-[#061f12]/80"
					>
						<span>🍎</span>
						{$t('download.appstore')}
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- FAQ Section -->
	<section class="px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-4xl">
			<h2 class="text-center text-4xl font-bold text-[#061f12] sm:text-5xl">{$t('faq.title')}</h2>

			<div class="mt-12 space-y-6">
				{#each ['q1', 'q2', 'q3', 'q4', 'q5', 'q6'] as key (key)}
					<div class="rounded-lg border border-[#00ee57]/20 bg-white p-6">
						<h3 class="font-semibold text-[#061f12]">{$t(`faq.${key}`)}</h3>
						<p class="mt-2 text-[#4b6658]">{$t(`faq.a${key.substring(1)}`)}</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t border-[#00ee57]/10 bg-[#061f12] px-4 py-12 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<div class="grid gap-8 sm:grid-cols-4">
				<div>
					<div class="flex items-center gap-2">
						<div class="text-2xl">🌿</div>
						<span class="text-xl font-bold text-white">Water Me</span>
					</div>
					<p class="mt-2 text-sm text-white/70">Smart plant care companion</p>
				</div>

				<div>
					<h4 class="font-semibold text-white">{$t('footer.product')}</h4>
					<ul class="mt-2 space-y-1">
						<li>
							<a href="#features" class="text-sm text-white/70 transition hover:text-white"
								>{$t('nav.features')}</a
							>
						</li>
						<li>
							<a href="#download" class="text-sm text-white/70 transition hover:text-white"
								>{$t('nav.download')}</a
							>
						</li>
						<li>
							<a
								href="mailto:support@water-me.app"
								class="text-sm text-white/70 transition hover:text-white">Support</a
							>
						</li>
					</ul>
				</div>

				<div>
					<h4 class="font-semibold text-white">{$t('footer.legal')}</h4>
					<ul class="mt-2 space-y-1">
						<li>
							<a
								href={resolve('/privacy')}
								class="text-sm text-white/70 transition hover:text-white">{$t('footer.privacy')}</a
							>
						</li>
						<li>
							<a href={resolve('/terms')} class="text-sm text-white/70 transition hover:text-white"
								>{$t('footer.terms')}</a
							>
						</li>
						<li>
							<a
								href={resolve('/imprint')}
								class="text-sm text-white/70 transition hover:text-white">{$t('footer.imprint')}</a
							>
						</li>
					</ul>
				</div>

				<div>
					<h4 class="font-semibold text-white">{$t('footer.connect')}</h4>
					<ul class="mt-2 space-y-1">
						<li>
							<a
								rel="external"
								target="_blank"
								href={TIKTOK_URL}
								class="text-sm text-white/70 transition hover:text-white">TikTok</a
							>
						</li>
						<li>
							<a
								rel="external"
								target="_blank"
								href={INSTAGRAM_URL}
								class="text-sm text-white/70 transition hover:text-white">Instagram</a
							>
						</li>
						<li>
							<a
								rel="external"
								target="_blank"
								href={GITHUB_URL}
								class="text-sm text-white/70 transition hover:text-white">GitHub</a
							>
						</li>
					</ul>
				</div>
			</div>

			<div class="mt-8 border-t border-white/10 pt-8">
				<p class="text-center text-sm text-white/70">
					{$t('footer.copyright')}
				</p>
			</div>
		</div>
	</footer>
</div>
