<script lang="ts">
	import { tStore } from '$lib/i18n';

	export let text: string = 'Message';
	export let type: 'info' | 'error' | 'success' = 'info';

	const variants: Record<typeof type, string> = {
		info: 'border-blue-400 bg-blue-100 text-blue-800',
		success: 'border-green-400 bg-green-100 text-green-800',
		error: 'border-red-400 bg-red-100 text-red-800'
	};

	$: classes = `mb-5 rounded-lg border-2 px-4 py-3 ${variants[type]}`;
</script>

<div class={classes}>
	<p class="text-sm font-semibold">{$tStore(text)}</p>
</div>
