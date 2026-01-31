<script lang="ts">
	interface Props {
		icon?: string;
		message?: string;
	}

	const { icon = '🌱', message = 'Loading...' }: Props = $props();
</script>

<div class="flex min-h-96 items-center justify-center">
	<div class="text-center">
		<div class="mb-4 animate-spin text-5xl">{icon}</div>
		<p class="text-lg text-[var(--text-light-main)]">{message}</p>
	</div>
</div>
