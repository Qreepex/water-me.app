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
	<path
		stroke-linecap="round"
		stroke-linejoin="round"
		stroke-width="2"
		d="M3 12l9-9 9 9M5 10v10a1 1 0 001 1h12a1 1 0 001-1V10H5z"
	></path>
</svg>
