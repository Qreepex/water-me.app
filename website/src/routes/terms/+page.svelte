<script>
	import { resolve } from '$app/paths';
</script>

<svelte:head>
	<title>Terms & Conditions - Water Me</title>
	<meta name="description" content="Terms and conditions for Water Me app." />
	<meta property="og:title" content="Terms & Conditions - Water Me" />
	<meta property="og:description" content="Terms and conditions for Water Me app." />
	<link rel="canonical" href="https://water-me.app/terms" />
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
		<h1 class="text-5xl font-bold text-[#061f12]">Terms & Conditions</h1>
		<p class="mt-2 text-[#4b6658]">Last updated: February 2026</p>

		<div class="prose prose-sm mt-12 max-w-none text-[#4b6658]">
			<h2 class="text-3xl font-semibold text-[#061f12]">1. Acceptance of Terms</h2>
			<p>
				By accessing and using the Water Me application ("Service"), you accept and agree to be
				bound by the terms and provision of this agreement. If you do not agree to abide by the
				above, please do not use this service.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">2. Use License</h2>
			<p>
				Permission is granted to temporarily download one copy of the materials (information or
				software) on Water Me's application for personal, non-commercial transitory viewing only.
				This is the grant of a license, not a transfer of title, and under this license you may not:
			</p>
			<ul class="list-inside list-disc space-y-2">
				<li>Modify or copy the materials</li>
				<li>Use the materials for any commercial purpose or for any public display</li>
				<li>Attempt to decompile or reverse engineer any software contained on the Service</li>
				<li>Remove any copyright or other proprietary notations from the materials</li>
				<li>
					Transfer the materials to another person or "mirror" the materials on any other server
				</li>
			</ul>

			<h2 class="text-3xl font-semibold text-[#061f12]">3. Disclaimer</h2>
			<p>
				The materials on Water Me's application are provided on an 'as is' basis. Water Me makes no
				warranties, expressed or implied, and hereby disclaims and negates all other warranties
				including, without limitation, implied warranties or conditions of merchantability, fitness
				for a particular purpose, or non-infringement of intellectual property or other violation of
				rights.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">4. Limitations</h2>
			<p>
				In no event shall Water Me or its suppliers be liable for any damages (including, without
				limitation, damages for loss of data or profit, or due to business interruption) arising out
				of the use or inability to use the materials on Water Me's application.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">5. Accuracy of Materials</h2>
			<p>
				The materials appearing on Water Me's application could include technical, typographical, or
				photographic errors. Water Me does not warrant that any of the materials on its application
				are accurate, complete, or current. Water Me may make changes to the materials contained on
				its application at any time without notice.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">6. Links</h2>
			<p>
				Water Me has not reviewed all of the sites linked to its application and is not responsible
				for the contents of any such linked site. The inclusion of any link does not imply
				endorsement by Water Me of the site. Use of any such linked website is at the user's own
				risk.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">7. Modifications</h2>
			<p>
				Water Me may revise these terms of service for its application at any time without notice.
				By using this application, you are agreeing to be bound by the then current version of these
				terms of service.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">8. Governing Law</h2>
			<p>
				These terms and conditions are governed by and construed in accordance with the laws of the
				jurisdiction in which Water Me operates, and you irrevocably submit to the exclusive
				jurisdiction of the courts in that location.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">9. User Accounts</h2>
			<p>
				When you create an account, you must provide information that is accurate, complete, and
				current at all times. You are responsible for safeguarding the password and for all
				activities that occur under your account. You must notify us immediately of any unauthorized
				use of your account.
			</p>

			<h2 class="text-3xl font-semibold text-[#061f12]">10. Contact Us</h2>
			<p>If you have any questions about these Terms & Conditions, please contact us at:</p>
			<ul class="list-inside list-disc space-y-2">
				<li>Email: legal@water-me.app</li>
				<li>Address: See imprint for contact information</li>
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

	:global(p) {
		margin-bottom: 1rem;
	}

	:global(ul) {
		margin-bottom: 1rem;
	}
</style>
