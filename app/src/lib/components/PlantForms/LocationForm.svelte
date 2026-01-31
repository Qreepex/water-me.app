<script lang="ts">
	import type { FormData } from '$lib/types/forms';

	interface Props {
		formData: FormData;
	}

	let { formData = $bindable() }: Props = $props();
</script>

<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
	<h2 class="mb-4 text-2xl font-bold text-green-800">📍 Location</h2>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
		<div>
			<label for="room" class="mb-1 block text-sm font-semibold text-gray-700">Room</label>
			<input
				type="text"
				id="room"
				bind:value={formData.room}
				placeholder="e.g., Living Room"
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<div>
			<label for="position" class="mb-1 block text-sm font-semibold text-gray-700">
				Position
			</label>
			<input
				type="text"
				id="position"
				bind:value={formData.position}
				placeholder="e.g., Windowsill"
				class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
			/>
		</div>

		<div>
			<fieldset>
				<legend class="mb-1 block text-sm font-semibold text-gray-700">Environment</legend>
				<label class="flex items-center gap-2">
					<input type="checkbox" bind:checked={formData.isOutdoors} class="h-4 w-4" />
					<span class="text-sm text-gray-700">Outdoors</span>
				</label>
			</fieldset>
		</div>
	</div>
</div>
