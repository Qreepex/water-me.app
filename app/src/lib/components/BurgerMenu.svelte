<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';

	let isOpen = $state(false);
	let selectedLanguage = $state('en');

	function toggleMenu() {
		isOpen = !isOpen;
	}

	function closeMenu() {
		isOpen = false;
	}

	async function handleLogout() {
		await authStore.logout();
		closeMenu();
	}

	function handleLanguageChange(lang: string) {
		selectedLanguage = lang;
		// TODO: Implement language switching
		console.log('Language changed to:', lang);
	}

	function navigateTo(path: string) {
		goto(path);
		closeMenu();
	}
</script>

<div class="relative">
	<!-- Burger Icon Button -->
	<button
		onclick={toggleMenu}
		class="p-3 rounded-full bg-white shadow-lg hover:shadow-xl transition-all hover:scale-105 border border-emerald-200"
		aria-label="Menu"
	>
		<svg
			class="w-6 h-6 text-emerald-700"
			fill="none"
			stroke="currentColor"
			viewBox="0 0 24 24"
		>
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
		<div class="fixed inset-0 bg-white z-[100] overflow-y-auto safe-area">
			<!-- Close Button -->
			<div class="flex justify-between items-center p-6 border-b border-emerald-200">
				<h2 class="text-2xl font-bold text-emerald-700">Menu</h2>
				<button
					onclick={closeMenu}
					class="p-2 rounded-full hover:bg-emerald-100 transition-colors"
					aria-label="Close menu"
				>
					<svg class="w-6 h-6 text-emerald-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>

			<!-- User Info -->
			{#if $authStore.user}
				<div class="p-6 bg-emerald-50 border-b border-emerald-200">
					<div class="flex items-center gap-3">
						<div class="w-12 h-12 bg-emerald-600 rounded-full flex items-center justify-center">
							<svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
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
			{/if}

			<!-- Menu Items -->
			<div class="p-6 space-y-2">
				<button
					onclick={() => navigateTo('/overview')}
					class="w-full text-left px-6 py-4 hover:bg-emerald-50 rounded-xl transition-colors flex items-center gap-3 text-lg"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path 
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
              d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
						></path>
					</svg>
					<span class="text-emerald-900 font-medium">My Plants</span>
				</button>

				<button
					onclick={() => navigateTo('/profile')}
					class="w-full text-left px-6 py-4 hover:bg-emerald-50 rounded-xl transition-colors flex items-center gap-3 text-lg"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
						></path>
					</svg>
					<span class="text-emerald-900 font-medium">User Profile</span>
				</button>

				<button
					onclick={() => navigateTo('/manage')}
					class="w-full text-left px-6 py-4 hover:bg-emerald-50 rounded-xl transition-colors flex items-center gap-3 text-lg"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 6v6m0 0v6m0-6h6m-6 0H6"
						></path>
					</svg>
					<span class="text-emerald-900 font-medium">Manage Plants</span>
				</button>
			</div>

			<!-- Language Picker -->
			<div class="px-6 py-4 border-t border-emerald-200">
				<p class="text-sm font-semibold text-emerald-700 mb-3">Language</p>
				<div class="flex gap-3">
						<button
							onclick={() => handleLanguageChange('en')}
							class={`px-6 py-3 rounded-lg text-base font-medium transition-colors ${
								selectedLanguage === 'en'
									? 'bg-emerald-600 text-white'
									: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
							}`}
						>
							EN
						</button>
						<button
							onclick={() => handleLanguageChange('de')}
							class={`px-6 py-3 rounded-lg text-base font-medium transition-colors ${
								selectedLanguage === 'de'
									? 'bg-emerald-600 text-white'
									: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
							}`}
						>
							DE
						</button>
						<button
							onclick={() => handleLanguageChange('es')}
							class={`px-6 py-3 rounded-lg text-base font-medium transition-colors ${
								selectedLanguage === 'es'
									? 'bg-emerald-600 text-white'
									: 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
							}`}
						>
							ES
						</button>
					</div>
			</div>

			<!-- Footer Links -->
			<div class="px-6 py-4 border-t border-emerald-200 space-y-2">
				<a
					href="/"
					class="block px-4 py-3 text-base text-emerald-700 hover:bg-emerald-50 rounded-lg transition-colors"
					onclick={closeMenu}
				>
					Website
				</a>
				<a
					href="/privacy"
					class="block px-4 py-3 text-base text-emerald-700 hover:bg-emerald-50 rounded-lg transition-colors"
					onclick={closeMenu}
				>
					Privacy Policy
				</a>
				<a
					href="/impressum"
					class="block px-4 py-3 text-base text-emerald-700 hover:bg-emerald-50 rounded-lg transition-colors"
					onclick={closeMenu}
				>
					Impressum
				</a>
			</div>

			<!-- Build Info -->
			<div class="px-6 py-4 border-t border-emerald-200 text-sm text-emerald-600">
				<p>Build: {new Date().toISOString().split('T')[0]}</p>
				<p>Version: 1.0.0</p>
			</div>

			<!-- Logout Button -->
			<div class="px-6 py-4 border-t border-emerald-200">
				<button
					onclick={handleLogout}
					class="w-full text-center px-6 py-4 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors font-semibold text-lg"
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
