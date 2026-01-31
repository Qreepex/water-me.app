<script lang="ts">
	const { class: className = 'h-5 w-5' } = $props<{ class?: string }>();
</script>

<svg
	class={className}
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 32 32"
	width="32"
	height="32"
	fill="none"
>
	<!-- Körper -->
	<rect
		x="10"
		y="14"
		width="14"
		height="12"
		rx="3"
		fill="none"
		stroke="currentColor"
		stroke-width="2.2"
		stroke-linejoin="round"
	/>
	<!-- Deckel -->
	<line
		x1="13"
		y1="14"
		x2="21"
		y2="14"
		stroke="currentColor"
		stroke-width="2.2"
		stroke-linecap="round"
	/>
	<!-- Griff -->
	<path
		d="M 24 17 Q 29 17 29 22"
		fill="none"
		stroke="currentColor"
		stroke-width="2.2"
		stroke-linecap="round"
	/>
	<!-- Ausguss -->
	<line
		x1="10"
		y1="18"
		x2="4"
		y2="12"
		stroke="currentColor"
		stroke-width="2.2"
		stroke-linecap="round"
	/>
	<!-- Tropfen -->
	<circle cx="2" cy="9" r="1.2" fill="currentColor" />
	<circle cx="5" cy="6" r="1" fill="currentColor" />
</svg>
