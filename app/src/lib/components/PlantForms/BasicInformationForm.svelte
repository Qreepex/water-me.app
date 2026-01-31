<script lang="ts">
	import type { FormData } from '$lib/types/forms';
	import { SunlightRequirement } from '$lib/types/api';

	interface Props {
		formData: FormData;
	}

	let { formData = $bindable() }: Props = $props();
</script>

<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
	<h2 class="mb-4 text-2xl font-bold text-green-800">ℹ️ Basic Information</h2>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<!-- Name -->
		<div>
			<label for="name" class="mb-1 block text-sm font-semibold text-gray-700">
				Plant Name *
			</label>
			<input
				type="text"
				id="name"
				required
				bind:value={formData.name}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<!-- Species -->
		<div>
			<label for="species" class="mb-1 block text-sm font-semibold text-gray-700"> Species *</label>
			<input
				type="text"
				id="species"
				required
				bind:value={formData.species}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<!-- Sunlight -->
		<div>
			<label for="sunlight" class="mb-1 block text-sm font-semibold text-gray-700">
				Sunlight Requirements
			</label>
			<select
				id="sunlight"
				bind:value={formData.sunlight}
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			>
				{#each Object.values(SunlightRequirement) as req (req)}
					<option value={req}>{req}</option>
				{/each}
			</select>
		</div>

		<!-- Toxic -->
		<div>
			<fieldset>
				<legend class="mb-1 block text-sm font-semibold text-gray-700">Safety</legend>
				<label class="flex items-center gap-2">
					<input type="checkbox" bind:checked={formData.isToxic} class="h-4 w-4" />
					<span class="text-sm text-gray-700">Toxic to pets/children</span>
				</label>
			</fieldset>
		</div>

		<!-- Temperature -->
		<div>
			<label for="temp" class="mb-1 block text-sm font-semibold text-gray-700">
				Preferred Temperature (°C): <span class="font-bold text-emerald-600"
					>{formData.preferedTemperature}</span
				>
			</label>
			<input
				type="range"
				id="temp"
				min="-50"
				max="100"
				bind:value={formData.preferedTemperature}
				class="w-full accent-emerald-600"
			/>
		</div>

		<!-- Target Humidity -->
		<div>
			<label for="humidity" class="mb-1 block text-sm font-semibold text-gray-700">
				Target Humidity (%): <span class="font-bold text-emerald-600"
					>{formData.targetHumidity}</span
				>
			</label>
			<input
				type="range"
				id="humidity"
				min="0"
				max="100"
				bind:value={formData.targetHumidity}
				class="w-full accent-emerald-600"
			/>
		</div>
	</div>
</div>
