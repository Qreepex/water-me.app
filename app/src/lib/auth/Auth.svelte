<script lang="ts">
	import Spinner from '$lib/components/Spinner.svelte';
	import { FirebaseAuthentication } from '@capacitor-firebase/authentication';
	import { Capacitor } from '@capacitor/core';
	import { initializeApp } from 'firebase/app';
	import {
		getAuth,
		setPersistence,
		browserLocalPersistence,
		onAuthStateChanged,
		type User
	} from 'firebase/auth';
	import { onMount } from 'svelte';
	import { FIREBASE_CONFIG } from './firebase';
	import { tStore } from '$lib/i18n';
	import Button from '$lib/components/Button.svelte';
	import Message from '$lib/components/Message.svelte';

	const app = initializeApp(FIREBASE_CONFIG);
	const auth = getAuth(app);

	let user: User | null = null;
	let loading = false;
	let initializing = true;
	let error: string | null = null;

	const platform = Capacitor.getPlatform();

	onMount(async () => {
		if (platform === 'web') {
			setPersistence(auth, browserLocalPersistence).catch(console.error);

			const unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
				user = firebaseUser ?? null;
				initializing = false;
			});

			return unsubscribe;
		} else {
			const result = await FirebaseAuthentication.getCurrentUser();
			user = result.user ?? null;
			initializing = false;

			const listener = await FirebaseAuthentication.addListener('authStateChange', (res) => {
				user = res.user;
			});

			return async () => listener.remove();
		}
	});

	async function loginWithGoogle() {
		try {
			loading = true;
			error = null;
			// Startet den nativen Google-Dialog auf dem Handy
			const result = await FirebaseAuthentication.signInWithGoogle();

			// Check if user cancelled (result exists but no user)
			if (!result || !result.user) {
				error = 'auth.signInCancelled';
				return;
			}

			user = result.user;

			// Das ist das Token für dein Backend!
			const idToken = await FirebaseAuthentication.getIdToken();
			console.log('Token für API:', idToken.token);

			// Hier kannst du das Token an dein Backend senden
		} catch (err: any) {
			console.error('Login fehlgeschlagen', err);
			// Check if user cancelled the popup
			if (
				err?.message?.includes('popup_closed_by_user') ||
				err?.code === 'popup-closed-by-user' ||
				err?.message?.includes('cancelled')
			) {
				error = 'auth.signInCancelled';
			} else {
				error = 'auth.signInError';
			}
		} finally {
			loading = false;
		}
	}
</script>

{#if initializing}
	<Spinner />
{:else if user}
	<slot />
{:else}
	<div
		class="flex min-h-screen items-center justify-center bg-gradient-to-br from-green-50 via-emerald-50 to-teal-50 p-4"
	>
		<div class="w-full max-w-md">
			<!-- Logo/Title -->
			<div class="mb-8 text-center">
				<h1 class="mb-2 text-5xl font-bold text-green-800">{$tStore('common.app')}</h1>
				<p class="text-green-700">{$tStore('common.appDescription')}</p>
			</div>

			<!-- Card -->
			<div class="rounded-2xl bg-white p-8 shadow-lg">
				<!-- Mode Indicator -->
				<div class="mb-8">
					<h2 class="mb-4 text-2xl font-bold text-green-800">{$tStore('auth.signIn')}</h2>
					<p class="text-gray-600">{$tStore('auth.signInToContinue')}</p>
				</div>

				<!-- Error Message -->
				{#if error}
					<Message text={error} type="error" />
				{/if}

				{#if loading}
					<Message text="auth.signingIn" />
				{/if}

				<!-- Submit Button -->
				<Button
					disabled={loading}
					{loading}
					onclick={loginWithGoogle}
					text="auth.signInWithGoogle"
					loadingText="auth.signingIn"
				/>
			</div>
		</div>
	</div>
{/if}
