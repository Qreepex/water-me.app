<script lang="ts">
	import { tStore } from '$lib/i18n';
	import type { SortOption } from '$lib/utils/plant';

	interface Props {
		sortBy?: SortOption;
		onSortChange?: (value: SortOption) => void;
	}

	let { sortBy = 'nameAsc', onSortChange }: Props = $props();

	const sortOptions: { value: SortOption; label: string; icon: string }[] = [
		{ value: 'nameAsc', label: 'plants.sortOptions.nameAsc', icon: '🔤' },
		{ value: 'nameDesc', label: 'plants.sortOptions.nameDesc', icon: '🔤' },
		{ value: 'lastWateredAsc', label: 'plants.sortOptions.lastWateredAsc', icon: '💧' },
		{ value: 'lastWateredDesc', label: 'plants.sortOptions.lastWateredDesc', icon: '💧' },
		{ value: 'speciesAsc', label: 'plants.sortOptions.speciesAsc', icon: '🌿' },
		{ value: 'speciesDesc', label: 'plants.sortOptions.speciesDesc', icon: '🌿' }
	];

	function handleChange(e: Event) {
		const value = (e.currentTarget as HTMLSelectElement).value as SortOption;
		sortBy = value;
		onSortChange?.(value);
	}
</script>


<div class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row sm:items-center sm:gap-3">
	<select
		id="sort"
		value={sortBy}
		onchange={handleChange}
		class="w-full rounded-lg border-2 border-[var(--p-emerald)] bg-[var(--card-light)] px-4 py-2 text-sm font-medium text-[var(--text-light-main)] transition hover:border-[var(--p-emerald-dark)] focus:border-[var(--p-emerald)] focus:outline-none sm:w-auto"
	>
		{#each sortOptions as option (option.value)}
			<option value={option.value}>{option.icon} {$tStore(option.label)}</option>
		{/each}
	</select>
</div>
