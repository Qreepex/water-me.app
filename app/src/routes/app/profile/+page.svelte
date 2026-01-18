<!-- <script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { setLanguage } from '$lib/stores/language';
	import { tSync } from '$lib/i18n';
	import { API_BASE_URL } from '$lib/constants';
	import { resolve } from '$app/paths';

	let token = $state($authStore.token);
	let language = $state('en' as 'en' | 'de' | 'es');

	let username = $state('');
	let email = $state('');
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');

	let message = $state('');
	let messageType = $state<'success' | 'error' | ''>('');
	let loading = $state(false);

	let initialized = $state(false);

	let checkAuthInterval: NodeJS.Timeout;

	onMount(() => {
		checkAuthInterval = setInterval(() => {
			if ($authStore.initialized) {
				clearInterval(checkAuthInterval);
				initialized = true;

				if (!$authStore.isAuthenticated) {
					goto(resolve('/'));
					return;
				}

				token = $authStore.token;
				if ($authStore.user) {
					username = $authStore.user.username || '';
					email = $authStore.user.email || '';
					if ($authStore.user.language) {
						language = $authStore.user.language as 'en' | 'de' | 'es';
					}
				}
			}
		}, 50);
	});

	onDestroy(() => {
		if (checkAuthInterval) {
			clearInterval(checkAuthInterval);
		}
	});

	async function handleUpdateProfile() {
		if (!token) {
			message = tSync('common.error', language);
			messageType = 'error';
			return;
		}

		loading = true;
		message = '';

		try {
			const response = await fetch(API_BASE_URL + '/api/user', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({
					username,
					email,
					language
				})
			});

			if (!response.ok) {
				if (response.status === 401) {
					await authStore.logout();
					goto(resolve('/'));
					return;
				}
				const error = await response.json();
				throw new Error(error.error || tSync('common.error', language));
			}

			const updatedUser = await response.json();
			await authStore.setUser(updatedUser);

			if (updatedUser.language) {
				language = updatedUser.language as 'en' | 'de' | 'es';
				await setLanguage(language);
			}

			message = tSync('profile.profileUpdated', language);
			messageType = 'success';
		} catch (err) {
			message = err instanceof Error ? err.message : tSync('common.error', language);
			messageType = 'error';
		} finally {
			loading = false;
			setTimeout(() => {
				message = '';
			}, 3000);
		}
	}

	async function handleChangePassword() {
		if (!token) {
			message = tSync('common.error', language);
			messageType = 'error';
			return;
		}

		if (newPassword !== confirmPassword) {
			message = tSync('profile.passwordsDoNotMatch', language);
			messageType = 'error';
			return;
		}

		if (newPassword.length < 6) {
			message = tSync('profile.passwordTooShort', language);
			messageType = 'error';
			return;
		}

		loading = true;
		message = '';

		try {
			const response = await fetch(API_BASE_URL + '/api/user/password', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({
					currentPassword,
					newPassword
				})
			});

			if (!response.ok) {
				if (response.status === 401) {
					const error = await response.json();
					throw new Error(error.error || tSync('profile.currentPasswordIncorrect', language));
				}
				const error = await response.json();
				throw new Error(error.error || tSync('common.error', language));
			}

			currentPassword = '';
			newPassword = '';
			confirmPassword = '';

			message = tSync('profile.passwordChanged', language);
			messageType = 'success';
		} catch (err) {
			message = err instanceof Error ? err.message : tSync('common.error', language);
			messageType = 'error';
		} finally {
			loading = false;
			setTimeout(() => {
				message = '';
			}, 3000);
		}
	}
</script>

{#if initialized}
	<div class="mx-auto max-w-2xl">
		<h1 class="mb-8 text-3xl font-bold text-emerald-900">
			{tSync('profile.userProfile', language)}
		</h1>

		{#if message}
			<div
				class={`mb-6 rounded-lg p-4 ${
					messageType === 'success'
						? 'border border-green-300 bg-green-100 text-green-800'
						: 'border border-red-300 bg-red-100 text-red-800'
				}`}
			>
				{message}
			</div>
		{/if}

		<div class="mb-8 rounded-lg border border-emerald-200 bg-white p-6 shadow-md">
			<h2 class="mb-6 text-xl font-bold text-emerald-900">
				{tSync('profile.profileInformation', language)}
			</h2>

			<div class="space-y-4">
				<div>
					<label for="username" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.username', language)}
					</label>
					<input
						id="username"
						type="text"
						bind:value={username}
						placeholder="Your username"
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="email" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.email', language)}
					</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						placeholder="your@email.com"
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="language" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.language', language)}
					</label>
					<select
						id="language"
						bind:value={language}
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					>
						<option value="en">English</option>
						<option value="de">Deutsch</option>
						<option value="es">Español</option>
					</select>
				</div>

				<button
					onclick={handleUpdateProfile}
					disabled={loading}
					class="w-full rounded-lg bg-emerald-600 py-2 font-semibold text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{loading ? tSync('profile.updating', language) : tSync('profile.updateProfile', language)}
				</button>
			</div>
		</div>

		<div class="rounded-lg border border-emerald-200 bg-white p-6 shadow-md">
			<h2 class="mb-6 text-xl font-bold text-emerald-900">
				{tSync('profile.changePassword', language)}
			</h2>

			<div class="space-y-4">
				<div>
					<label for="currentPassword" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.currentPassword', language)}
					</label>
					<input
						id="currentPassword"
						type="password"
						bind:value={currentPassword}
						placeholder="Enter current password"
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="newPassword" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.newPassword', language)}
					</label>
					<input
						id="newPassword"
						type="password"
						bind:value={newPassword}
						placeholder="Enter new password"
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="confirmPassword" class="mb-2 block text-sm font-semibold text-emerald-800">
						{tSync('profile.confirmPassword', language)}
					</label>
					<input
						id="confirmPassword"
						type="password"
						bind:value={confirmPassword}
						placeholder="Confirm new password"
						class="w-full rounded-lg border border-emerald-300 px-4 py-2 focus:ring-2 focus:ring-emerald-500 focus:outline-none"
					/>
				</div>

				<button
					onclick={handleChangePassword}
					disabled={loading}
					class="w-full rounded-lg bg-emerald-600 py-2 font-semibold text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{loading
						? tSync('profile.changing', language)
						: tSync('profile.changePasswordButton', language)}
				</button>
			</div>
		</div>
	</div>
{/if} -->
