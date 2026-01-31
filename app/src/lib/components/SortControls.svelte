<script lang="ts">
	type SortOption =
		| 'name'
		| 'lastWatered'
		| 'lastFertilized'
		| 'wateringIntervalDays'
		| 'mistingIntervalDays';

	interface Props {
		sortBy?: SortOption;
		onSortChange?: (value: SortOption) => void;
	}

	let { sortBy = 'name', onSortChange }: Props = $props();

	function handleChange(e: Event) {
		const value = (e.currentTarget as HTMLSelectElement).value as SortOption;
		sortBy = value;
		onSortChange?.(value);
	}
</script>

<div class="flex items-center gap-3">
	<label for="sort" class="font-semibold text-[var(--text-light-main)]">Sort by:</label>
	<select
		id="sort"
		value={sortBy}
		onchange={handleChange}
		class="rounded-lg border-2 border-[var(--p-emerald)] bg-[var(--card-light)] px-4 py-2 font-medium text-[var(--text-light-main)] transition hover:border-[var(--p-emerald-dark)] focus:border-[var(--p-emerald)] focus:outline-none"
	>
		<option value="name">Plant Name</option>
		<option value="lastWatered">Last Watered</option>
		<option value="lastFertilized">Last Fertilized</option>
		<option value="wateringIntervalDays">Watering Frequency</option>
		<option value="mistingIntervalDays">Spray Frequency</option>
	</select>
</div>
