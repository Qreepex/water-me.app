<script lang="ts">
	import type { FormData } from '$lib/types/forms';

	interface Props {
		formData: FormData;
	}

	let { formData = $bindable() }: Props = $props();
</script>

<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
	<h2 class="mb-4 text-xl font-bold text-green-800">💦 Misting & Humidity</h2>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<label class="flex items-center gap-2">
			<input type="checkbox" bind:checked={formData.requiresMisting} class="h-4 w-4" />
			<span class="text-sm font-semibold text-gray-700">Requires misting</span>
		</label>

		<label class="flex items-center gap-2">
			<input type="checkbox" bind:checked={formData.requiresHumidifier} class="h-4 w-4" />
			<span class="text-sm font-semibold text-gray-700">Needs humidifier</span>
		</label>

		{#if formData.requiresMisting}
			<div>
				<label for="mist-interval" class="mb-1 block text-sm font-semibold text-gray-700">
					Misting Interval (days)
				</label>
				<input
					type="number"
					id="mist-interval"
					min="1"
					bind:value={formData.mistingIntervalDays}
					class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
				/>
			</div>
		{/if}
	</div>
</div>
