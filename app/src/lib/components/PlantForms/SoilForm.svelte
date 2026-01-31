<script lang="ts">
	import type { FormData } from '$lib/types/forms';
	import Button from '../ui/Button.svelte';

	interface Props {
		formData: FormData;
		soilComponentInput: string;
	}

	let { formData = $bindable(), soilComponentInput = $bindable() }: Props = $props();

	function addSoilComponent(): void {
		if (soilComponentInput.trim()) {
			formData.soilComponents = [...formData.soilComponents, soilComponentInput.trim()];
			soilComponentInput = '';
		}
	}

	function removeSoilComponent(index: number): void {
		formData.soilComponents = formData.soilComponents.filter((_, i) => i !== index);
	}
</script>

<details class="space-y-4">
	<summary class="cursor-pointer font-semibold text-green-700 select-none">
		<span>🪴 Soil</span>
		<span class="ml-2 text-sm text-gray-600">▶</span>
	</summary>

	<div class="ml-2 space-y-4">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="soil-type" class="mb-1 block text-sm font-semibold text-gray-700">
					Soil Type
				</label>
				<input
					type="text"
					id="soil-type"
					bind:value={formData.soilType}
					placeholder="e.g., Peat moss"
					class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
				/>
			</div>

			<div>
				<label for="repot-cycle" class="mb-1 block text-sm font-semibold text-gray-700">
					Repotting Cycle (years)
				</label>
				<input
					type="number"
					id="repot-cycle"
					min="1"
					bind:value={formData.repottingCycle}
					class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
				/>
			</div>
		</div>

		<div>
			<fieldset>
				<legend class="mb-2 block text-sm font-semibold text-gray-700">Soil Components</legend>
				<div class="mb-2 flex gap-2">
					<input
						type="text"
						bind:value={soilComponentInput}
						placeholder="e.g., Perlite, Orchid bark"
						class="flex-1 rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
					/>
					<Button onclick={addSoilComponent} text="Add" variant="primary" />
				</div>

				{#if formData.soilComponents.length > 0}
					<div class="space-y-2">
						{#each formData.soilComponents as component, i (i)}
							<div class="flex items-center justify-between rounded-lg bg-blue-50 p-2">
								<span class="text-sm text-gray-800">{component}</span>
								<Button
									text="Remove"
									variant="danger"
									size="sm"
									onclick={() => removeSoilComponent(i)}
								/>
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-gray-500 italic">No components added</p>
				{/if}
			</fieldset>
		</div>
	</div>
</details>
