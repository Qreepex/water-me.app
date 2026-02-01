<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Burger from '$lib/assets/Burger.svg.svelte';
	import Can from '$lib/assets/Can.svg.svelte';
	import SunFlower from '$lib/assets/SunFlower.svg.svelte';
	import { tStore } from '$lib/i18n';
	import BurgerMenu from './BurgerMenu.svelte';

	let showMenu = $state(false);

	function isActive(path: '/' | '/water' | '/manage'): boolean {
		return page.url.pathname === resolve(path) || page.url.pathname.startsWith(resolve(path) + '/');
	}

	function navigate(path: '/' | '/water' | '/manage'): void {
		showMenu = false;
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
			class="flex flex-1 cursor-pointer flex-col items-center justify-center gap-1 py-2 transition-colors {!showMenu &&
			isActive('/')
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Home"
		>
			<SunFlower isActive={isActive('/')} />
			<span class="text-xs font-medium">{$tStore('menu.garden')}</span>
		</button>

		<!-- Water -->
		<button
			onclick={() => navigate('/water')}
			class="flex flex-1 cursor-pointer flex-col items-center justify-center gap-1 py-2 transition-colors {!showMenu &&
			isActive('/water')
				? 'text-emerald-600'
				: 'text-gray-600'}"
			aria-label="Water"
		>
			<Can isActive={isActive('/water')} />
			<span class="text-xs font-medium">{$tStore('menu.water')}</span>
		</button>

		<!-- Menu -->
		<button
			onclick={toggleMenu}
			class="flex flex-1 cursor-pointer flex-col items-center justify-center gap-1 py-2 transition-colors {showMenu
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
	<div class="pt-safe fixed inset-0 bottom-20 z-50 bg-white">
		<div class="flex h-full flex-col overflow-hidden px-3 pt-2 md:px-10 md:pt-10 xl:px-32 xl:pt-14">
			<BurgerMenu />
		</div>
	</div>
{/if}

<style>
	.pt-safe {
		padding-top: env(safe-area-inset-top);
	}

	.pb-safe {
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
