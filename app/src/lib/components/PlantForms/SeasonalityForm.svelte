<script lang="ts">
	import type { FormData } from '$lib/types/forms';

	interface Props {
		formData: FormData;
	}

	const { formData }: Props = $props();
</script>

<details class="space-y-4">
	<summary class="cursor-pointer font-semibold text-green-700 select-none">
		<span>❄️ Seasonality</span>
		<span class="ml-2 text-sm text-gray-600">▶</span>
	</summary>

	<div class="ml-2 space-y-4">
		<label class="flex items-center gap-2">
			<input type="checkbox" bind:checked={formData.winterRestPeriod} class="h-4 w-4" />
			<span class="text-sm font-semibold text-gray-700">Has winter rest period</span>
		</label>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="winter-factor" class="mb-1 block text-sm font-semibold text-gray-700">
					Winter Water Factor: <span class="font-bold text-emerald-600"
						>{(formData.winterWaterFactor * 100).toFixed(0)}%</span
					>
				</label>
				<input
					type="range"
					id="winter-factor"
					min="0"
					max="1"
					step="0.1"
					bind:value={formData.winterWaterFactor}
					class="w-full accent-blue-600"
				/>
				<p class="mt-1 text-xs text-gray-500">Multiply normal watering by this factor</p>
			</div>

			<div>
				<label for="min-temp" class="mb-1 block text-sm font-semibold text-gray-700">
					Minimum Temperature (°C)
				</label>
				<input
					type="number"
					id="min-temp"
					bind:value={formData.minTempCelsius}
					class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
				/>
			</div>
		</div>
	</div>
</details>
