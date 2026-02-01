<script>
	import { resolve } from '$app/paths';
</script>

<svelte:head>
	<title>Privacy Policy - Water Me</title>
	<meta name="description" content="Privacy policy for Water Me app. Learn how we protect your data." />
	<meta property="og:title" content="Privacy Policy - Water Me" />
	<meta
		property="og:description"
		content="Privacy policy for Water Me app. Learn how we protect your data."
	/>
	<link rel="canonical" href="https://water-me.app/privacy" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-[#f2faf5] via-white to-[#e8f9f0]">
	<!-- Navigation -->
	<nav class="sticky top-0 z-50 border-b border-[#00ee57]/10 bg-white/80 backdrop-blur-md">
		<div class="mx-auto max-w-6xl px-4 py-4 sm:px-6 lg:px-8">
			<div class="flex items-center gap-3">
				<a href={resolve('/')} class="flex items-center gap-3 transition hover:opacity-80">
					<div class="text-3xl">🌿</div>
					<span class="text-2xl font-bold text-[#061f12]">Water Me</span>
				</a>
			</div>
		</div>
	</nav>

	<!-- Content -->
	<div class="mx-auto max-w-4xl px-4 py-20 sm:px-6 lg:px-8">
		<h1 class="text-5xl font-bold text-[#061f12]">Privacy Policy</h1>
		<p class="mt-2 text-[#4b6658]">Last updated: February 2026</p>

		<div class="prose prose-sm mt-12 max-w-none text-[#4b6658]">
			<h2 class="text-3xl font-semibold text-[#061f12]">1. Introduction</h2>
			<p>
				Water Me ("we", "us", "our" or "Company") operates the Water Me mobile application (the "Service").
			</p>
			<p>
				This page informs you of our policies regarding the collection, use, and disclosure of personal data
				when you use our Service and the choices you have associated with that data.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">2. Information Collection and Use</h2>
			<p>We collect several different types of information for various purposes to provide and improve our Service to you.</p>

			<h3 class="text-2xl font-semibold text-[#061f12]">Types of Data Collected:</h3>
			<ul class="list-inside list-disc space-y-2">
				<li>
					<strong>Personal Data:</strong> When you create an account, we may collect information such as your email address
					and authentication credentials.
				</li>
				<li>
					<strong>Plant Data:</strong> Information about your plants, including names, species, care schedules, watering frequency,
					and photos.
				</li>
				<li>
					<strong>Usage Data:</strong> Information about how you interact with the app, including feature usage and timestamps.
				</li>
				<li>
					<strong>Device Information:</strong> Information about your device such as device type, operating system version,
					and unique device identifiers.
				</li>
				<li>
					<strong>Location Data:</strong> We do not collect precise location data without your explicit consent.
				</li>
			</ul>

			<h2 class="text-3xl font-semibold text-[#061f12]">3. Use of Data</h2>
			<p>Water Me uses the collected data for various purposes:</p>
			<ul class="list-inside list-disc space-y-2">
				<li>To provide and maintain the Service</li>
				<li>To notify you about changes to our Service</li>
				<li>To send you push notifications for plant care reminders</li>
				<li>To gather analysis or valuable information so that we can improve the Service</li>
				<li>To monitor the usage of the Service</li>
				<li>To detect, prevent and address technical issues</li>
			</ul>

			<h2 class="text-3xl font-semibold text-[#061f12]">4. Security of Data</h2>
			<p>
				The security of your data is important to us but remember that no method of transmission over the Internet
				or method of electronic storage is 100% secure. While we strive to use commercially acceptable means to
				protect your Personal Data, we cannot guarantee its absolute security.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">5. Changes to This Privacy Policy</h2>
			<p>
				We may update our Privacy Policy from time to time. We will notify you of any changes by posting the new
				Privacy Policy on this page and updating the "effective date" at the top of this Privacy Policy.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">6. Contact Us</h2>
			<p>If you have any questions about this Privacy Policy, please contact us at:</p>
			<ul class="list-inside list-disc space-y-2">
				<li>Email: privacy@water-me.app</li>
				<li>Address: See imprint for contact information</li>
			</ul>

			<h2 class="text-3xl font-semibold text-[#061f12]">7. Your Rights</h2>
			<p>You have the right to:</p>
			<ul class="list-inside list-disc space-y-2">
				<li>Access your personal data</li>
				<li>Correct inaccurate data</li>
				<li>Request deletion of your data</li>
				<li>Object to processing of your data</li>
				<li>Request restriction of processing</li>
				<li>Data portability</li>
			</ul>
		</div>
	</div>

	<!-- Footer -->
	<footer class="border-t border-[#00ee57]/10 bg-[#061f12] px-4 py-12 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-6xl">
			<div class="text-center text-sm text-white/70">
				© {new Date().getFullYear()} Water Me. All rights reserved.
			</div>
		</div>
	</footer>
</div>

<style>
	:global(h2) {
		margin: 2rem 0 1rem;
	}

	:global(h3) {
		margin: 1.5rem 0 0.75rem;
	}

	:global(p) {
		margin-bottom: 1rem;
	}

	:global(ul) {
		margin-bottom: 1rem;
	}
</style>
