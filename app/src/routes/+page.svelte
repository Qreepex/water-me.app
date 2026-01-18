<script lang="ts">
	import { onMount } from 'svelte';

	onMount(() => {
		// // Wait for auth store to initialize from Capacitor preferences
		// const unsubscribe = authStore.subscribe((state) => {
		// 	if (state.isAuthenticated) {
		// 		goto(resolve('/app'));
		// 	}
		// });
		// return unsubscribe;
	});

	// async function handleSubmit(e: Event) {
	// 	e.preventDefault();
	// 	loading = true;
	// 	error = null;

	// 	try {
	// 		const endpoint = mode === 'login' ? '/api/login' : '/api/signup';
	// 		const response = await fetch(API_BASE_URL + `${endpoint}`, {
	// 			method: 'POST',
	// 			headers: { 'Content-Type': 'application/json' },
	// 			body: JSON.stringify({ email, password })
	// 		});

	// 		const data = await response.json();

	// 		if (!response.ok) {
	// 			error = data.error || 'Authentication failed';
	// 			return;
	// 		}

	// 		// Store auth data and redirect
	// 		const user: User = data.user;
	// 		const token: string = data.token;

	// 		authStore.login(user, token);
	// 		goto(resolve('/app'));
	// 	} catch (err) {
	// 		error = err instanceof Error ? err.message : 'An error occurred';
	// 	} finally {
	// 		loading = false;
	// 	}
	// }

	// function toggleMode() {
	// 	mode = mode === 'login' ? 'signup' : 'login';
	// 	error = null;
	// }
</script>

<h1>Welcome to Plants</h1>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
</style>
