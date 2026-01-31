<script lang="ts">
	import type { FormData } from '$lib/types/forms';
	import { WateringMethod, WaterType } from '$lib/types/api';

	interface Props {
		formData: FormData;
	}

	let { formData = $bindable() }: Props = $props();
</script>

<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
	<h2 class="mb-4 text-xl font-bold text-green-800">💧 Watering</h2>

	<div class="space-y-4">
		<div>
			<label for="water-interval" class="mb-1 block text-sm font-semibold text-gray-700">
				Interval (days)
			</label>
			<input
				type="number"
				id="water-interval"
				min="1"
				bind:value={formData.wateringIntervalDays}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<div>
			<label for="water-method" class="mb-1 block text-sm font-semibold text-gray-700">
				Method
			</label>
			<select
				id="water-method"
				bind:value={formData.wateringMethod}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			>
				{#each Object.values(WateringMethod) as method (method)}
					<option value={method}>{method}</option>
				{/each}
			</select>
		</div>

		<div>
			<label for="water-type" class="mb-1 block text-sm font-semibold text-gray-700">
				Water Type
			</label>
			<select
				id="water-type"
				bind:value={formData.waterType}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			>
				{#each Object.values(WaterType) as type (type)}
					<option value={type}>{type}</option>
				{/each}
			</select>
		</div>
	</div>
</div>
