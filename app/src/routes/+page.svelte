<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth';
	import type { User } from '$lib/auth/auth';
	import { API_BASE_URL } from '$lib/constants';

	let mode: 'login' | 'signup' = 'login';
	let email = '';
	let password = '';
	let loading = false;
	let error: string | null = null;
	let isInitialized = false;

	onMount(async () => {
		// Wait for auth store to initialize from Capacitor preferences
		const unsubscribe = authStore.subscribe((state) => {
			isInitialized = state.initialized;
			if (state.isAuthenticated) {
				goto('/overview');
			}
		});

		return unsubscribe;
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = null;

		try {
			const endpoint = mode === 'login' ? '/api/login' : '/api/signup';
			const response = await fetch(API_BASE_URL + `${endpoint}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			const data = await response.json();

			if (!response.ok) {
				error = data.error || 'Authentication failed';
				return;
			}

			// Store auth data and redirect
			const user: User = data.user;
			const token: string = data.token;

			authStore.login(user, token);
			goto('/overview');
		} catch (err) {
			error = err instanceof Error ? err.message : 'An error occurred';
		} finally {
			loading = false;
		}
	}

	function toggleMode() {
		mode = mode === 'login' ? 'signup' : 'login';
		error = null;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-green-50 via-emerald-50 to-teal-50 flex items-center justify-center p-4">
	<div class="w-full max-w-md">
		<!-- Logo/Title -->
		<div class="text-center mb-8">
			<h1 class="text-5xl font-bold text-green-800 mb-2">🌱 PlantCare</h1>
			<p class="text-green-700">Take care of your green friends</p>
		</div>

		<!-- Card -->
		<div class="bg-white rounded-2xl shadow-lg p-8">
			<!-- Mode Indicator -->
			<div class="mb-8">
				<h2 class="text-2xl font-bold text-green-800 mb-4">
					{mode === 'login' ? 'Welcome Back' : 'Create Account'}
				</h2>
				<p class="text-gray-600">
					{mode === 'login'
						? 'Sign in to manage your plants'
						: 'Join us to start tracking your plants'}
				</p>
			</div>

			<!-- Form -->
			<form on:submit={handleSubmit} class="space-y-5">
				<!-- Email Input -->
				<div>
					<label for="email" class="block text-sm font-semibold text-green-800 mb-2">
						Email
					</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						placeholder="you@example.com"
						required
						disabled={loading}
						class="w-full rounded-lg border-2 border-green-300 px-4 py-3 transition focus:border-green-500 focus:outline-none disabled:bg-gray-100 hover:border-green-400"
					/>
				</div>

				<!-- Password Input -->
				<div>
					<label for="password" class="block text-sm font-semibold text-green-800 mb-2">
						Password
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						placeholder={mode === 'signup' ? 'At least 6 characters' : '••••••'}
						required
						disabled={loading}
						class="w-full rounded-lg border-2 border-green-300 px-4 py-3 transition focus:border-green-500 focus:outline-none disabled:bg-gray-100 hover:border-green-400"
					/>
				</div>

				<!-- Error Message -->
				{#if error}
					<div class="rounded-lg border-2 border-red-400 bg-red-100 px-4 py-3 text-red-800">
						<p class="text-sm font-semibold">{error}</p>
					</div>
				{/if}

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={loading}
					class="w-full rounded-lg bg-green-600 px-4 py-3 font-semibold text-white transition hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if loading}
						<span class="inline-flex items-center gap-2">
							<svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							{mode === 'login' ? 'Signing in...' : 'Creating account...'}
						</span>
					{:else}
						{mode === 'login' ? 'Sign In' : 'Sign Up'}
					{/if}
				</button>
			</form>

			<!-- Toggle Mode -->
			<div class="mt-6 text-center">
				<p class="text-gray-600">
					{mode === 'login' ? "Don't have an account?" : 'Already have an account?'}
					<button
						type="button"
						on:click={toggleMode}
						class="font-semibold text-green-600 hover:text-green-700 transition"
					>
						{mode === 'login' ? 'Sign Up' : 'Sign In'}
					</button>
				</p>
			</div>

			<!-- Demo Note -->
			<div class="mt-8 rounded-lg bg-emerald-50 border border-emerald-300 p-4">
				<p class="text-xs text-emerald-800">
					💡 <strong>Demo tip:</strong> Use any email and password (min 6 chars) to get started!
				</p>
			</div>
		</div>

		<!-- Footer -->
		<div class="mt-8 text-center text-sm text-gray-600">
			<p>Made with 🌿 for plant lovers</p>
		</div>
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
</style>
