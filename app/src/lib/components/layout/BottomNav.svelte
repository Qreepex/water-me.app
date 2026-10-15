<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Burger from '$lib/assets/Burger.svg.svelte';
	import Can from '$lib/assets/Can.svg.svelte';
	import SunFlower from '$lib/assets/SunFlower.svg.svelte';
	import { tStore } from '$lib/i18n';
	import type { Component } from 'svelte';
	import AppMenu from './AppMenu.svelte';

	let menuOpen = $state(false);

	const tabs = [
		{ path: '/', icon: SunFlower, labelKey: 'menu.garden' },
		{ path: '/water', icon: Can, labelKey: 'tasks.title' }
	] as const satisfies readonly {
		path: '/' | '/water';
		icon: Component<{ isActive: boolean }>;
		labelKey: string;
	}[];

	function isActive(path: '/' | '/water'): boolean {
		return page.url.pathname === resolve(path) || page.url.pathname.startsWith(resolve(path) + '/');
	}

	function navigate(path: '/' | '/water'): void {
		menuOpen = false;
		goto(resolve(path));
	}
</script>

<!-- Bottom Navigation Bar -->
<nav
	class="fixed right-0 bottom-0 left-0 z-40 border-t border-ink/5 bg-surface/90 shadow-[0_-6px_24px_rgba(0,0,0,0.08)] backdrop-blur-xl select-none"
>
	<div class="pb-safe flex h-20 items-stretch justify-around">
		{#each tabs as tab (tab.path)}
			{@const active = !menuOpen && isActive(tab.path)}
			<button
				onclick={() => navigate(tab.path)}
				class="flex flex-1 cursor-pointer flex-col items-center justify-center gap-0.5 py-2 transition-colors"
				aria-label={$tStore(tab.labelKey)}
				aria-current={active ? 'page' : undefined}
			>
				<span
					class="flex h-8 w-16 items-center justify-center rounded-full transition-colors duration-200 {active
						? 'bg-brand/25'
						: ''}"
				>
					<tab.icon isActive={active} />
				</span>
				<span
					class="text-xs transition-colors {active
						? 'font-semibold text-brand-dark dark:text-brand'
						: 'font-medium text-ink-soft'}"
				>
					{$tStore(tab.labelKey)}
				</span>
			</button>
		{/each}

		<button
			onclick={() => (menuOpen = !menuOpen)}
			class="flex flex-1 cursor-pointer flex-col items-center justify-center gap-0.5 py-2 transition-colors"
			aria-label={$tStore('menu.menu')}
			aria-expanded={menuOpen}
		>
			<span
				class="flex h-8 w-16 items-center justify-center rounded-full transition-colors duration-200 {menuOpen
					? 'bg-brand/25'
					: ''}"
			>
				<Burger isActive={menuOpen} />
			</span>
			<span
				class="text-xs transition-colors {menuOpen
					? 'font-semibold text-brand-dark dark:text-brand'
					: 'font-medium text-ink-soft'}"
			>
				{$tStore('menu.menu')}
			</span>
		</button>
	</div>
</nav>

<!-- Menu Overlay -->
{#if menuOpen}
	<div class="pt-safe fixed inset-0 bottom-20 z-50 bg-canvas">
		<div class="flex h-full flex-col overflow-hidden px-2 pt-2 md:px-10 md:pt-10 xl:px-32 xl:pt-14">
			<AppMenu onClose={() => (menuOpen = false)} />
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
