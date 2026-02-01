<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { invalidateApiCache } from '$lib/utils/cache';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { tStore } from '$lib/i18n';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Message from '$lib/components/ui/Message.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';

	let error: string | null = null;
	let creating = false;

	async function createEmptyPlant(): Promise<void> {
		creating = true;
		error = null;

		try {
			const createPayload = {
				name: $tStore('plants.newPlant') || 'New Plant',
				species: '',
				isToxic: false,
				preferedTemperature: 0,
				photoIds: [],
				flags: [],
				notes: [],
				pestHistory: [],
				growthHistory: []
			};

			const createRes = await fetchData('/api/plants', {
				method: 'post' as const,
				body: createPayload
			});

			if (!createRes.ok) {
				const errMsg = createRes.error?.message || 'Failed to create plant';
				console.error('Create plant error:', createRes.error);
				throw new Error(errMsg);
			}

			const newPlant = createRes.data;
			if (!newPlant || !newPlant.id) {
				console.error('Invalid response data:', newPlant);
				throw new Error('Invalid response from server');
			}

			console.log('Created new plant:', newPlant);

			// Invalidate cache so new plant appears in list and detail page
			await invalidateApiCache(['/api/plants', `/api/plants/${newPlant.id}`], {
				waitForAck: true,
				timeoutMs: 200
			});

			// Navigate to edit page
			goto(resolve(`/manage/${newPlant.id}`));
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
			console.error('Create plant exception:', err);
			creating = false;
		}
	}
</script>

<div class="flex h-full min-h-0 flex-col overflow-hidden">
	<div class="flex-shrink-0">
		<PageHeader icon="🌱" title="plants.newPlant" description="plants.startAddingPlants">
			<Button
				variant="ghost"
				size="sm"
				text="common.back"
				icon="←"
				onclick={() => goto(resolve('/manage'))}
			/>
		</PageHeader>
	</div>

	<PageContent>
		<div class="flex min-h-0 flex-1 items-center justify-center">
			<div class="flex max-w-sm flex-col items-center gap-6">
				{#if error}
					<Message type="error" title="common.error" description={error} />
				{/if}

				{#if creating}
					<LoadingSpinner message="plants.creating" icon="🌿" />
				{:else}
					<div class="text-center">
						<div class="mb-4 text-6xl">🌱</div>
						<h1 class="mb-2 text-2xl font-bold text-gray-900">
							{$tStore('plants.newPlant')}
						</h1>
						<p class="mb-8 text-gray-600">
							{$tStore('plants.createDescription') ||
								'Create a new plant and fill in the details in the next step.'}
						</p>
					</div>

					<Button
						variant="primary"
						size="lg"
						onclick={createEmptyPlant}
						text="plants.createPlant"
						class="w-full cursor-pointer"
					/>
				{/if}
			</div>
		</div>
	</PageContent>
</div>
