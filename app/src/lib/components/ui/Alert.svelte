<script lang="ts">
	type AlertType = 'error' | 'success' | 'warning' | 'info';

	interface Props {
		type?: AlertType;
		title?: string;
		message?: string;
		description?: string;
	}

	const {
		type = 'info',
		title = undefined,
		message = undefined,
		description = undefined
	}: Props = $props();

	const styles: Record<AlertType, { border: string; bg: string; text: string }> = {
		error: {
			border: 'border-[var(--status-error)]',
			bg: 'bg-[var(--status-error)]/10',
			text: 'text-[var(--status-error)]'
		},
		success: {
			border: 'border-[var(--status-success)]',
			bg: 'bg-[var(--status-success)]/10',
			text: 'text-[var(--status-success)]'
		},
		warning: {
			border: 'border-[var(--status-warn)]',
			bg: 'bg-[var(--status-warn)]/10',
			text: 'text-[var(--status-warn)]'
		},
		info: {
			border: 'border-[var(--status-info)]',
			bg: 'bg-[var(--status-info)]/10',
			text: 'text-[var(--status-info)]'
		}
	};

	const style = $derived(styles[type]);
	const content = $derived(description || message);
</script>

<div class="rounded-lg border-2 {style.border} {style.bg} px-6 py-4 {style.text}">
	{#if title}
		<p class="font-bold">{title}</p>
	{/if}
	{#if content}
		<p>{content}</p>
	{/if}
</div>
