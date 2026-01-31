<script lang="ts">
	interface Props {
		isActive: boolean;
	}

	let { isActive }: Props = $props();
</script>

<svg
	class="h-6 w-6"
	fill={isActive ? 'currentColor' : 'none'}
	stroke="currentColor"
	viewBox="0 0 24 24"
>
	<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"
	></path>
</svg>
