<script lang="ts">
	import type { FormData } from '$lib/types/forms';
	import { PlantFlag } from '$lib/types/api';

	interface Props {
		formData: FormData;
		newNote: string;
	}

	let { formData, newNote = $bindable() }: Props = $props();

	function toggleFlag(flag: PlantFlag): void {
		if (formData.flags.includes(flag)) {
			formData.flags = formData.flags.filter((f) => f !== flag);
		} else {
			formData.flags = [...formData.flags, flag];
		}
	}

	function addNote(): void {
		if (newNote.trim()) {
			formData.notes = [...formData.notes, newNote.trim()];
			newNote = '';
		}
	}

	function removeNote(index: number): void {
		formData.notes = formData.notes.filter((_, i) => i !== index);
	}
</script>

<details class="space-y-4">
	<summary class="cursor-pointer font-semibold text-green-700 select-none">
		<span>🏷️ Flags & Notes</span>
		<span class="ml-2 text-sm text-gray-600">▶</span>
	</summary>

	<div class="ml-2 space-y-4">
		<!-- Flags -->
		<div>
			<span class="mb-2 block text-sm font-semibold text-gray-700">Flags</span>
			<div class="space-y-2">
				{#each Object.values(PlantFlag) as flag (flag)}
					<label class="flex items-center gap-2">
						<input
							type="checkbox"
							checked={formData.flags.includes(flag)}
							onchange={() => toggleFlag(flag)}
							class="h-4 w-4"
						/>
						<span class="text-sm text-gray-700">{flag}</span>
					</label>
				{/each}
			</div>
		</div>

		<!-- Notes -->
		<fieldset>
			<legend class="mb-2 block text-sm font-semibold text-gray-700">Notes</legend>
			<div class="mb-2 flex gap-2">
				<input
					type="text"
					bind:value={newNote}
					onkeydown={(e) => e.key === 'Enter' && addNote()}
					placeholder="Add a note..."
					class="flex-1 rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
				/>
				<button
					type="button"
					onclick={addNote}
					class="rounded-lg bg-emerald-600 px-4 py-2 text-white shadow-sm transition hover:bg-emerald-700"
				>
					Add
				</button>
			</div>

			{#if formData.notes.length > 0}
				<div class="space-y-2">
					{#each formData.notes as note, i (i)}
						<div class="flex items-start justify-between rounded-lg bg-blue-50 p-3">
							<p class="flex-1 text-sm text-gray-800">{note}</p>
							<button
								type="button"
								onclick={() => removeNote(i)}
								class="ml-2 font-bold text-red-500 hover:text-red-700"
							>
								×
							</button>
						</div>
					{/each}
				</div>
			{:else}
				<p class="text-sm text-gray-500 italic">No notes yet</p>
			{/if}
		</fieldset>
	</div>
</details>
