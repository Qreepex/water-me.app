<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import FormActionBar from '$lib/components/layout/FormActionBar.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Spinner from '$lib/components/ui/Spinner.svelte';
	import { tStore } from '$lib/i18n';
	import type { Plant } from '$lib/types/api';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import { ALLOWED_PHOTO_TYPES, compressPhoto, putToPresignedUrl } from '$lib/utils/uploads';
	import { onDestroy } from 'svelte';

	type PhotoItem = {
		fileName: string;
		previewUrl: string;
		status: 'pending' | 'compressing' | 'uploading' | 'uploaded' | 'error';
		error?: string;
		key?: string;
	};

	let plant = $state<Plant | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let previewUrls = $state<string[]>([]);
	let photos = $state<PhotoItem[]>([]);
	let uploadedPhotoKeys = $state<string[]>([]);
	let removedPhotoIds = $state<string[]>([]);

	const existingPhotoIds = $derived(plant?.photoIds ?? []);
	const hasPhotoChanges = $derived(uploadedPhotoKeys.length > 0 || removedPhotoIds.length > 0);

	onMount(async () => {
		try {
			const plantId = page.params.plant ?? '';
			const response = await fetchData('/api/plants/{id}', {
				params: { id: plantId }
			});

			if (!response.ok) {
				error = response.error?.message || $tStore('plants.failedToFetchPlants');
				return;
			}

			plant = response.data;
			await loadPhotoPreviews();
		} catch (err) {
			error = err instanceof Error ? err.message : $tStore('plants.failedToFetchPlants');
		} finally {
			loading = false;
		}
	});

	onDestroy(() => {
		previewUrls.forEach((u) => revokeObjectURL(u));
		photos.forEach((p) => p.previewUrl && URL.revokeObjectURL(p.previewUrl));
	});

	async function loadPhotoPreviews(): Promise<void> {
		if (!plant) return;
		const ids = plant.photoIds || [];
		const urls = (plant as { photoUrls?: string[] }).photoUrls ?? [];
		previewUrls = [];
		for (let i = 0; i < ids.length; i++) {
			const id = ids[i];
			const url = urls[i];
			if (!id || !url) continue;
			const objUrl = await getImageObjectURL(id, url);
			if (objUrl) previewUrls.push(objUrl);
		}
	}

	function onFilesSelected(e: Event): void {
		const input = e.target as HTMLInputElement;
		const files = input.files ? Array.from(input.files) : [];
		if (!files.length) return;
		photos = files.map((f) => ({
			fileName: f.name,
			previewUrl: URL.createObjectURL(f),
			status: 'pending'
		}));
		processUploads(files).catch((err) => {
			error = err instanceof Error ? err.message : $tStore('common.uploadError');
		});
	}

	function setItemStatus(
		index: number,
		status: PhotoItem['status'],
		errorMessage?: string,
		key?: string
	): void {
		const item = photos[index];
		if (!item) return;
		item.status = status;
		item.error = errorMessage;
		if (key) item.key = key;
		photos = [...photos];
	}

	async function processUploads(files: File[]): Promise<void> {
		for (let i = 0; i < files.length; i++) {
			const file = files[i];
			if (!ALLOWED_PHOTO_TYPES.has(file.type)) {
				setItemStatus(i, 'error', $tStore('common.unsupportedFileType'));
				continue;
			}

			setItemStatus(i, 'compressing');

			let compressed;
			try {
				compressed = await compressPhoto(file);
			} catch (err) {
				setItemStatus(
					i,
					'error',
					err instanceof Error ? err.message : $tStore('common.compressionFailed')
				);
				continue;
			}

			setItemStatus(i, 'uploading');

			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const presignRes = await (fetchData as any)('/api/uploads/presign', {
				method: 'post',
				body: {
					filename: compressed.outName,
					contentType: compressed.contentType,
					sizeBytes: compressed.blob.size
				}
			});
			if (!presignRes.ok) {
				setItemStatus(i, 'error', presignRes.error?.message || $tStore('common.failed'));
				continue;
			}

			const { url, headers, key } = presignRes.data as {
				url: string;
				headers: Record<string, string>;
				key: string;
			};

			const putOk = await putToPresignedUrl(url, headers, compressed.blob);
			if (!putOk) {
				setItemStatus(i, 'error', $tStore('common.failed'));
				continue;
			}

			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const regRes = await (fetchData as any)('/api/uploads/register', {
				method: 'post',
				body: { key }
			});
			if (!regRes.ok) {
				setItemStatus(i, 'error', regRes.error?.message || $tStore('common.failed'));
				continue;
			}

			setItemStatus(i, 'uploaded', undefined, key);
			uploadedPhotoKeys = [...uploadedPhotoKeys, key];
		}
	}

	function removeExistingPhoto(photoId: string, index: number): void {
		if (!confirm($tStore('plants.deletePhotoConfirm'))) return;
		removedPhotoIds = [...removedPhotoIds, photoId];
		const newUrls = [...previewUrls];
		const urlToRevoke = newUrls[index];
		newUrls.splice(index, 1);
		previewUrls = newUrls;
		if (urlToRevoke) revokeObjectURL(urlToRevoke);
		void deletePhotoFromS3(photoId);
	}

	async function deletePhotoFromS3(key: string): Promise<void> {
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			await (fetchData as any)(`/api/uploads/${encodeURIComponent(key)}`, {
				method: 'delete'
			});
		} catch {
			// ignore
		}
	}

	async function saveAndContinue(): Promise<void> {
		if (!plant) return;
		saving = true;
		error = null;

		try {
			const keepIds = (plant.photoIds || []).filter((id) => !removedPhotoIds.includes(id));
			const allPhotoIds = [...keepIds, ...uploadedPhotoKeys];

			if (allPhotoIds.length !== (plant.photoIds || []).length || removedPhotoIds.length > 0) {
				const res = await fetchData('/api/plants/{id}', {
					method: 'patch',
					params: { id: plant.id },
					body: { photoIds: allPhotoIds }
				});

				if (!res.ok) {
					throw new Error(res.error?.message || $tStore('plants.failedToSavePhotos'));
				}
			}

			await goto(resolve(`/manage/${plant.id}/edit`));
		} catch (err) {
			error = err instanceof Error ? err.message : $tStore('plants.failedToSavePhotos');
		} finally {
			saving = false;
		}
	}

	function backToEditor(): void {
		if (!plant) return;
		goto(resolve(`/manage/${plant.id}/edit`));
	}
</script>

<PageHeader
	icon="camera"
	title="plants.photos"
	description="plants.managePhotosDescription"
	backHref={plant ? resolve(`/manage/${plant.id}/edit`) : undefined}
/>

<PageContent>
	{#if loading}
		<Spinner message="common.loadingPlantDetails" />
	{:else if !plant}
		<Alert
			type="error"
			title="common.error"
			description={error || $tStore('common.plantNotFound')}
		/>
	{:else}
		{#if error}
			<div class="mx-2 mb-3">
				<Alert type="error" title="common.error" description={error} />
			</div>
		{/if}

		<!-- pb reserves space for the fixed FormActionBar -->
		<div class="min-h-0 flex-1 space-y-4 overflow-y-auto px-2 pb-48">
			<Card padded>
				<label class="block cursor-pointer">
					<span class="text-sm font-semibold text-ink">{$tStore('plants.addImages')}</span>
					<input
						type="file"
						accept="image/jpeg,image/png,image/webp"
						multiple
						onchange={onFilesSelected}
						aria-label={$tStore('plants.addImages')}
						class="mt-2 w-full cursor-pointer rounded-xl border-2 border-dashed border-brand/30 bg-canvas p-3 text-sm text-ink file:mr-3 file:rounded-lg file:border-0 file:bg-brand file:px-4 file:py-2 file:text-sm file:font-semibold file:text-onbrand"
					/>
				</label>
			</Card>

			{#if photos.length}
				<Card padded>
					<p class="mb-2 text-sm font-semibold text-ink">{$tStore('plants.newUploads')}</p>
					<div class="grid grid-cols-2 gap-2">
						{#each photos as p (p.previewUrl)}
							<div class="overflow-hidden rounded-xl bg-canvas">
								<img src={p.previewUrl} alt={p.fileName} class="h-24 w-full object-cover" />
								<div class="p-2 text-sm">
									<div class="mb-1 truncate font-medium text-ink">{p.fileName}</div>
									{#if p.status === 'uploaded'}
										<span class="inline-flex items-center gap-1.5 font-semibold text-ok">
											<Icon name="check" size={14} />
											{$tStore('plants.uploadUploaded')}
										</span>
									{:else if p.status === 'error'}
										<span class="inline-flex items-center gap-1.5 text-danger">
											<Icon name="x" size={14} />
											{p.error || $tStore('common.error')}
										</span>
									{:else}
										<span class="text-ink-soft">
											{p.status === 'pending'
												? $tStore('plants.uploadPending')
												: p.status === 'compressing'
													? $tStore('plants.uploadCompressing')
													: $tStore('plants.uploadUploading')}
										</span>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</Card>
			{/if}

			{#if previewUrls.length}
				<Card padded>
					<p class="mb-2 text-sm font-semibold text-ink">{$tStore('plants.existingPhotos')}</p>
					<div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
						{#each previewUrls as u, i (u)}
							<div class="group relative overflow-hidden rounded-xl">
								<img
									src={u}
									alt={plant.name || $tStore('common.plant')}
									class="h-24 w-full object-cover"
								/>
								<button
									onclick={() => removeExistingPhoto(existingPhotoIds[i] ?? '', i)}
									aria-label={$tStore('plants.deletePhoto')}
									class="absolute inset-0 flex cursor-pointer items-center justify-center bg-danger/80 text-sm font-bold text-white opacity-0 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100"
								>
									{$tStore('common.delete')}
								</button>
							</div>
						{/each}
					</div>
				</Card>
			{/if}
		</div>

		<FormActionBar
			onCancel={backToEditor}
			onSave={saveAndContinue}
			{saving}
			canSave={hasPhotoChanges}
			cancelText="common.close"
		/>
	{/if}
</PageContent>
