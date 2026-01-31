<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/stores';
	import Burger from '$lib/assets/Burger.svg.svelte';
	import Can from '$lib/assets/Can.svg.svelte';
	import Home from '$lib/assets/Home.svg.svelte';
	import PenPaper from '$lib/assets/PenPaper.svg.svelte';
	import SunFlower from '$lib/assets/SunFlower.svg.svelte';
	import { tStore } from '$lib/i18n';
	import BurgerMenu from './BurgerMenu.svelte';

	let showMenu = $state(false);

	function isActive(path: '/' | '/water' | '/manage'): boolean {
		return (
			$page.url.pathname === resolve(path) || $page.url.pathname.startsWith(resolve(path) + '/')
		);
	}

	function navigate(path: '/' | '/water' | '/manage'): void {
		goto(resolve(path));
	}

	function toggleMenu() {
		showMenu = !showMenu;
	}
</script>

<!-- Bottom Navigation Bar -->
<div class="fixed right-0 bottom-0 left-0 z-40 border-t border-emerald-200 bg-white shadow-lg">
	<div class="pb-safe flex h-20 items-center justify-around">
		<!-- Home -->
		<button
			onclick={() => navigate('/')}
			class="flex flex-1 flex-col items-center justify-center gap-1 py-2 transition-colors {isActive(
				'/'
			)
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Home"
		>
			<Home isActive={isActive('/')} />
			<span class="text-xs font-medium">{$tStore('menu.home')}</span>
		</button>

		<!-- Water -->
		<button
			onclick={() => navigate('/water')}
			class="flex flex-1 flex-col items-center justify-center gap-1 py-2 transition-colors {isActive(
				'/water'
			)
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Water"
		>
			<Can />
			<span class="text-xs font-medium">{$tStore('menu.water')}</span>
		</button>

		<!-- Manage -->
		<button
			onclick={() => navigate('/manage')}
			class="flex flex-1 flex-col items-center justify-center gap-1 py-2 transition-colors {isActive(
				'/manage'
			)
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Manage"
		>
			<SunFlower isActive={isActive('/manage')} />
			<span class="text-xs font-medium">{$tStore('menu.garden')}</span>
		</button>

		<!-- Menu -->
		<button
			onclick={toggleMenu}
			class="flex flex-1 flex-col items-center justify-center gap-1 py-2 transition-colors {showMenu
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Menu"
		>
			<Burger isActive={showMenu} />
			<span class="text-xs font-medium">{$tStore('menu.menu')}</span>
		</button>
	</div>
</div>

<!-- Menu Overlay -->
{#if showMenu}
	<div class="fixed inset-0 bottom-20 z-50 overflow-y-auto bg-white">
		<div class="pt-safe pb-safe flex items-center justify-between border-b border-emerald-200 p-6">
			<h2 class="text-2xl font-bold text-emerald-700">Settings</h2>
			<button
				onclick={() => (showMenu = false)}
				class="rounded-full p-2 transition-colors hover:bg-emerald-100"
				aria-label="Close menu"
			>
				<Burger isActive={showMenu} />
			</button>
		</div>
		<BurgerMenu
			onClose={() => {
				showMenu = false;
			}}
		/>
	</div>
{/if}

<style>
	.pb-safe {
		padding-bottom: env(safe-area-inset-bottom);
	}

	.pt-safe {
		padding-top: env(safe-area-inset-top);
	}
</style>
