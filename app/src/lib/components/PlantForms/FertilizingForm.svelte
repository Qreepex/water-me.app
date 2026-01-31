<script lang="ts">
	import type { FormData } from '$lib/types/forms';
	import { FertilizerType } from '$lib/types/api';

	interface Props {
		formData: FormData;
	}

	const { formData }: Props = $props();
</script>

<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
	<h2 class="mb-4 text-xl font-bold text-green-800">🥗 Fertilizing</h2>

	<div class="space-y-4">
		<div>
			<label for="fert-interval" class="mb-1 block text-sm font-semibold text-gray-700">
				Interval (days)
			</label>
			<input
				type="number"
				id="fert-interval"
				min="1"
				bind:value={formData.fertilizingIntervalDays}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<div>
			<label for="fert-type" class="mb-1 block text-sm font-semibold text-gray-700">Type</label>
			<select
				id="fert-type"
				bind:value={formData.fertilizingType}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			>
				{#each Object.values(FertilizerType) as type (type)}
					<option value={type}>{type}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="npk-ratio" class="mb-1 block text-sm font-semibold text-gray-700">
				NPK Ratio
			</label>
			<input
				type="text"
				id="npk-ratio"
				bind:value={formData.npkRatio}
				placeholder="e.g., 10:10:10"
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<div>
			<label for="concentration" class="mb-1 block text-sm font-semibold text-gray-700">
				Concentration (%): <span class="font-bold text-emerald-600"
					>{formData.concentrationPercent}</span
				>
			</label>
			<input
				type="range"
				id="concentration"
				min="1"
				max="100"
				bind:value={formData.concentrationPercent}
				class="w-full accent-emerald-600"
			/>
		</div>

		<label class="flex items-center gap-2">
			<input type="checkbox" bind:checked={formData.activeInWinter} class="h-4 w-4" />
			<span class="text-sm text-gray-700">Fertilize in winter</span>
		</label>
	</div>
</div>
