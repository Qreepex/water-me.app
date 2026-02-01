<script>
	import { resolve } from '$app/paths';
	import { GITHUB_URL, INSTAGRAM_URL, TIKTOK_URL } from '$lib';

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
	<link rel="canonical" href="https://water-me.app" />
	<script type="application/ld+json">
		{
			"@context": "https://schema.org",
			"@type": "WebApplication",
			"name": "Water Me",
			"description": "Smart plant care reminder app",
			"url": "https://water-me.app",
			"applicationCategory": "Productivity",
			"operatingSystem": "Android",
			"offers": {
				"@type": "Offer",
				"price": "0",
				"priceCurrency": "EUR"
			}
		}
	</script>
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
				<div class="hidden gap-8 sm:flex">
					<a href="#features" class="text-[#4b6658] transition hover:text-[#00ee57]">
						Features
					</a>
					<a href="#stats" class="text-[#4b6658] transition hover:text-[#00ee57]">
						Stats
					</a>
					<a href="#download" class="text-[#4b6658] transition hover:text-[#00ee57]">
						Download
					</a>
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
					<h1 class="text-5xl font-bold leading-tight text-[#061f12] sm:text-6xl">
						Never Forget to Water Your Plants Again
					</h1>
					<p class="mt-6 text-xl text-[#4b6658]">
						Water Me is your smart plant care companion. Get timely reminders for watering, spraying, and fertilizing. Track your entire plant collection with detailed care schedules and growth history.
					</p>

					<!-- CTA Buttons -->
					<div class="mt-8 flex gap-4">
						<a
							href="#download"
							class="inline-flex items-center gap-2 rounded-full bg-[#00ee57] px-8 py-4 font-semibold text-[#061f12] shadow-lg shadow-[#00ee57]/40 transition hover:shadow-xl hover:shadow-[#00ee57]/50"
						>
							📥 Download Now
						</a>
						<a
							href="#features"
							class="inline-flex items-center gap-2 rounded-full border-2 border-[#00ee57] px-8 py-4 font-semibold text-[#00ee57] transition hover:bg-[#00ee57]/10"
						>
							✨ Learn More
						</a>
					</div>

					<!-- Trust Badges -->
					<div class="mt-12 flex items-center gap-4">
						<div class="text-4xl">⭐⭐⭐⭐⭐</div>
						<div>
							<p class="font-semibold text-[#061f12]">Highly Rated</p>
							<p class="text-sm text-[#4b6658]">Loved by plant parents everywhere</p>
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
				<h2 class="text-4xl font-bold text-[#061f12] sm:text-5xl">Powerful Features</h2>
				<p class="mt-4 text-xl text-[#4b6658]">Everything you need to care for your plants</p>
			</div>

			<div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
				<!-- Feature 1 -->
				<div class="rounded-2xl border border-[#00ee57]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">💧</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Smart Watering Reminders</h3>
					<p class="text-[#4b6658]">
						Get timely reminders for watering based on your plant's specific needs and the last watering date.
					</p>
				</div>

				<!-- Feature 2 -->
				<div class="rounded-2xl border border-[#ffb703]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">🎯</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Multi-Care Tracking</h3>
					<p class="text-[#4b6658]">
						Track watering, fertilizing, spraying, and misting all in one place. Set intervals for each care type.
					</p>
				</div>

				<!-- Feature 3 -->
				<div class="rounded-2xl border border-[#219ebc]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">📸</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Growth Tracking</h3>
					<p class="text-[#4b6658]">
						Take photos of your plants and track their growth over time. Create detailed growth history.
					</p>
				</div>

				<!-- Feature 4 -->
				<div class="rounded-2xl border border-[#e63946]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">📊</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Detailed Care History</h3>
					<p class="text-[#4b6658]">
						Maintain comprehensive records including soil type, humidity, temperature, and seasonal variations.
					</p>
				</div>

				<!-- Feature 5 -->
				<div class="rounded-2xl border border-[#008f41]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">🔔</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Push Notifications</h3>
					<p class="text-[#4b6658]">
						Receive timely push notifications so you never miss a care task, no matter where you are.
					</p>
				</div>

				<!-- Feature 6 -->
				<div class="rounded-2xl border border-[#00ee57]/20 bg-white p-8 shadow-sm transition hover:shadow-lg">
					<div class="mb-4 text-5xl">🌍</div>
					<h3 class="mb-2 text-2xl font-semibold text-[#061f12]">Multi-Language Support</h3>
					<p class="text-[#4b6658]">
						Available in English, German, and Spanish. Care for your plants in your preferred language.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- Stats Section -->
	<section id="stats" class="bg-gradient-to-r from-[#00ee57] to-[#00a343] px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<div class="mb-12 text-center">
				<h2 class="text-4xl font-bold text-white sm:text-5xl">Our Community</h2>
				<p class="mt-4 text-lg text-white/90">Thousands of plant parents trust Water Me</p>
			</div>

			<div class="grid gap-8 sm:grid-cols-3">
				<!-- Stat 1 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">{displayStats.users.toLocaleString()}</div>
					<p class="mt-2 text-lg text-white/90">Active Users</p>
				</div>

				<!-- Stat 2 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">{displayStats.plants.toLocaleString()}</div>
					<p class="mt-2 text-lg text-white/90">Plants Tracked</p>
				</div>

				<!-- Stat 3 -->
				<div class="rounded-2xl bg-white/20 p-8 backdrop-blur-sm">
					<div class="text-5xl font-bold text-white">{displayStats.reminders.toLocaleString()}+</div>
					<p class="mt-2 text-lg text-white/90">Reminders Sent</p>
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
				{#each [
					{ emoji: '🏠', title: 'Home', desc: 'Quick overview of all your plants' },
					{ emoji: '💧', title: 'Watering', desc: 'Never miss a watering schedule' },
					{ emoji: '✏️', title: 'Management', desc: 'Full control over plant details' }
				] as item (item.title)}
					<div class="overflow-hidden rounded-2xl border border-[#00ee57]/20 bg-white shadow-sm">
						<div class="flex h-48 items-center justify-center bg-gradient-to-br from-[#00ee57]/10 to-[#00a343]/10 text-6xl">
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
			<div class="rounded-3xl border-2 border-[#00ee57] bg-gradient-to-br from-[#00ee57]/5 to-[#00a343]/5 p-12 text-center">
				<h2 class="text-4xl font-bold text-[#061f12] sm:text-5xl">Ready to Care for Your Plants?</h2>
				<p class="mt-4 text-lg text-[#4b6658]">
					Download Water Me today and join thousands of plant parents who never forget to water their plants.
				</p>

				<div class="mt-8 flex flex-col gap-4 sm:flex-row sm:justify-center">
					<a
						href="https://play.google.com/store/apps/details?id=app.waterme.app"
						class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#061f12] px-8 py-4 font-semibold text-white transition hover:bg-[#061f12]/80"
					>
						<span>🤖</span> Get on Google Play
					</a>
					<div
						class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#061f12] px-8 py-4 font-semibold text-white transition hover:bg-[#061f12]/80"
					>
						<span>🍎</span> Coming to iOS
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- FAQ Section -->
	<section class="px-4 py-20 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-4xl">
			<h2 class="text-center text-4xl font-bold text-[#061f12] sm:text-5xl">Frequently Asked Questions</h2>

			<div class="mt-12 space-y-6">
				{#each [
					{
						q: 'Is Water Me free?',
						a: 'Yes! Water Me is completely free to download and use. No hidden charges or premium upgrades required.'
					},
					{
						q: 'Does it work offline?',
						a: 'Yes, the app works offline. Your data syncs when you are back online.'
					},
					{
						q: 'What languages are supported?',
						a: 'Water Me is available in English, German, and Spanish. More languages coming soon!'
					},
					{
						q: 'Is my data private?',
						a: 'Your privacy is our priority. All your plant data is securely stored and never shared with third parties.'
					},
					{
						q: 'Can I use it for multiple plants?',
						a: 'Absolutely! Track unlimited plants with individual care schedules and reminders for each.'
					},
					{
						q: 'How do I get reminders?',
						a: 'Enable push notifications in the app settings and we will send you timely reminders for each care task.'
					}
				] as faq (faq.q)}
					<div class="rounded-lg border border-[#00ee57]/20 bg-white p-6">
						<h3 class="font-semibold text-[#061f12]">{faq.q}</h3>
						<p class="mt-2 text-[#4b6658]">{faq.a}</p>
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
					<h4 class="font-semibold text-white">Product</h4>
					<ul class="mt-2 space-y-1">
						<li><a href="#features" class="text-sm text-white/70 transition hover:text-white">Features</a></li>
						<li><a href="#download" class="text-sm text-white/70 transition hover:text-white">Download</a></li>
						<li><a href="mailto:support@water-me.app" class="text-sm text-white/70 transition hover:text-white">Support</a></li>
					</ul>
				</div>

				<div>
					<h4 class="font-semibold text-white">Legal</h4>
					<ul class="mt-2 space-y-1">
						<li><a href={resolve('/privacy')} class="text-sm text-white/70 transition hover:text-white">Privacy</a></li>
						<li><a href={resolve('/terms')} class="text-sm text-white/70 transition hover:text-white">Terms</a></li>
						<li><a href={resolve('/imprint')} class="text-sm text-white/70 transition hover:text-white">Imprint</a></li>
					</ul>
				</div>

				<div>
					<h4 class="font-semibold text-white">Connect</h4>
					<ul class="mt-2 space-y-1">
						<li><a rel="external" target="_blank" href={TIKTOK_URL} class="text-sm text-white/70 transition hover:text-white">TikTok</a></li>
						<li><a rel="external" target="_blank" href={INSTAGRAM_URL} class="text-sm text-white/70 transition hover:text-white">Instagram</a></li>
						<li><a  rel="external" target="_blank" href={GITHUB_URL} class="text-sm text-white/70 transition hover:text-white">GitHub</a></li>
					</ul>
				</div>
			</div>

			<div class="mt-8 border-t border-white/10 pt-8">
				<p class="text-center text-sm text-white/70">
					© {new Date().getFullYear()} Water Me. All rights reserved.
				</p>
			</div>
		</div>
	</footer>
</div>
