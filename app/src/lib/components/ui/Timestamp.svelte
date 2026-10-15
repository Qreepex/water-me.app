<script lang="ts">
	import { tStore } from '$lib/i18n';
	import { futureTimestamp, pastTimestamp } from '$lib/utils/timestamp.svelte';

	interface Props {
		date: Date;
		/** 'past' (default) or 'future' */
		mode?: 'past' | 'future';
	}

	let { date, mode = 'past' }: Props = $props();

	// Resolved on every render so the text stays reactive: computing this once
	// with t() would freeze the raw i18n key if the locale had not loaded yet.
	const text = $derived(mode === 'future' ? futureTimestamp(date) : pastTimestamp(date));
</script>

{$tStore(text.key, text.args)}
