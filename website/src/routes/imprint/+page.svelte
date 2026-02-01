<script>
	import { resolve } from '$app/paths';
</script>

<svelte:head>
	<title>Imprint - Water Me</title>
	<meta name="description" content="Imprint and legal information for Water Me app." />
	<meta property="og:title" content="Imprint - Water Me" />
	<meta property="og:description" content="Imprint and legal information for Water Me app." />
	<link rel="canonical" href="https://water-me.app/imprint" />
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
		<h1 class="text-5xl font-bold text-[#061f12]">Imprint</h1>
		<p class="mt-2 text-[#4b6658]">Legal information and contact details</p>

		<div class="prose prose-sm mt-12 max-w-none space-y-8 text-[#4b6658]">
			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Company Information</h2>
				<div class="mt-4 space-y-2">
					<p>
						<strong>Company Name:</strong><br />
						Water Me - Smart Plant Care
					</p>
					<p>
						<strong>Legal Form:</strong><br />
						Private Company
					</p>
					<p>
						<strong>Address:</strong><br />
						123 Garden Street<br />
						Berlin, 10115<br />
						Germany
					</p>
				</div>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Contact Information</h2>
				<div class="mt-4 space-y-2">
					<p>
						<strong>Email:</strong><br />
						<a href="mailto:info@water-me.app" class="text-[#00ee57] transition hover:text-[#00a343]">
							info@water-me.app
						</a>
					</p>
					<p>
						<strong>Support:</strong><br />
						<a href="mailto:support@water-me.app" class="text-[#00ee57] transition hover:text-[#00a343]">
							support@water-me.app
						</a>
					</p>
					<p>
						<strong>Phone:</strong><br />
						+49 (0) 30 123456
					</p>
				</div>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Management</h2>
				<div class="mt-4 space-y-2">
					<p>
						<strong>Management:</strong><br />
						Benjamin Schmitz
					</p>
				</div>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Disclaimer</h2>
				<p>
					The information contained on the Water Me website and application is for general informational purposes
					only. We assume no responsibility for errors or omissions in the contents on the Service.
				</p>
				<p>
					In no event shall Water Me, nor any of its officers, directors, and employees, be held liable for anything
					arising out of or in any way connected with your use of this Service.
				</p>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Liability for Links</h2>
				<p>
					The Service may contain links to external websites and resources. Water Me is not responsible for the
					availability of these external sites and does not endorse and is not responsible or liable for any content,
					advertising, products, or other materials on or available from such sites.
				</p>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Copyright & Intellectual Property</h2>
				<p>
					All content included in the Service, such as text, graphics, logos, images, audio clips, digital downloads,
					and data compilations, are the property of Water Me or its content suppliers and is protected by international
					copyright laws.
				</p>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Data Protection Officer</h2>
				<p>
					For questions regarding data protection and privacy, please contact our data protection officer at:
				</p>
				<p>
					<a href="mailto:privacy@water-me.app" class="text-[#00ee57] transition hover:text-[#00a343]">
						privacy@water-me.app
					</a>
				</p>
			</div>

			<div>
				<h2 class="text-3xl font-semibold text-[#061f12]">Dispute Resolution</h2>
				<p>
					The laws of Germany shall govern these terms and conditions. Any disputes arising under these terms and
					conditions shall be subject to the exclusive jurisdiction of the courts in Berlin, Germany.
				</p>
			</div>
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
		margin-bottom: 0.5rem;
	}
</style>
