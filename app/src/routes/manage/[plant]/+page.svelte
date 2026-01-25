<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { Plant } from '$lib/types/api';
	import { SunlightRequirement, WateringMethod, WaterType, FertilizerType } from '$lib/types/api';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/stores';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import type { FormData } from '$lib/types/forms';
	import { createEmptyFormData } from '$lib/types/forms';
	import BasicInformationForm from '$lib/components/PlantForms/BasicInformationForm.svelte';
	import LocationForm from '$lib/components/PlantForms/LocationForm.svelte';
	import WateringForm from '$lib/components/PlantForms/WateringForm.svelte';
	import FertilizingForm from '$lib/components/PlantForms/FertilizingForm.svelte';
	import MistingForm from '$lib/components/PlantForms/MistingForm.svelte';
	import SoilForm from '$lib/components/PlantForms/SoilForm.svelte';
	import SeasonalityForm from '$lib/components/PlantForms/SeasonalityForm.svelte';
	import MetadataForm from '$lib/components/PlantForms/MetadataForm.svelte';

	let plant = $state<Plant | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let success = $state<string | null>(null);
	let submitting = $state(false);
	let newNote = $state('');
	let soilComponentInput = $state('');
	let previewUrls = $state<string[]>([]);

	// Upload state
	const MAX_BYTES = 2 * 1024 * 1024; // 2MB
	const allowedTypes = new Set(['image/jpeg', 'image/png', 'image/webp']);
	type PhotoItem = {
		fileName: string;
		previewUrl: string;
		status: 'pending' | 'compressing' | 'uploading' | 'uploaded' | 'error';
		error?: string;
		key?: string;
	};
	let photos = $state<PhotoItem[]>([]);
	let uploadedPhotoKeys = $state<string[]>([]);
	let uploadTimestamps = $state<Record<string, number>>({}); // Track when each photo was uploaded
	let removedPhotoIds = $state<string[]>([]); // Track removed existing photos

	let formData = $state<FormData>(createEmptyFormData());

	onMount(async () => {
		try {
			const plantId = $page.params.plant ?? '';
			const response = await fetchData('/api/plants/{id}', {
				params: { id: plantId }
			});

			if (!response.ok) {
				error = response.error?.message || 'Failed to load plant';
				return;
			}

			plant = response.data;
			formData = initializeFormData();
			await loadPhotoPreviews();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load plant';
		} finally {
			loading = false;
		}
	});

	async function loadPhotoPreviews(): Promise<void> {
		if (!plant) return;
		const ids = plant.photoIds || [];
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		const urls = ((plant as any)?.photoUrls as string[] | undefined) || [];
		previewUrls = [];
		for (let i = 0; i < ids.length; i++) {
			const id = ids[i];
			const url = urls[i];
			if (!id || !url) continue;
			const objUrl = await getImageObjectURL(id, url);
			if (objUrl) previewUrls.push(objUrl);
		}
	}

	onDestroy(() => {
		previewUrls.forEach((u) => revokeObjectURL(u));
		photos.forEach((p) => p.previewUrl && URL.revokeObjectURL(p.previewUrl));
		// Clean up any unapplied uploads when leaving the page
		cleanupUnappliedUploads();
	});

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
			error = err instanceof Error ? err.message : 'Upload error';
		});
	}

	async function processUploads(files: File[]): Promise<void> {
		for (let i = 0; i < files.length; i++) {
			const file = files[i];
			const item = photos[i];
			if (!allowedTypes.has(file.type)) {
				item.status = 'error';
				item.error = 'Unsupported file type';
				photos = [...photos]; // Trigger reactivity
				continue;
			}
			item.status = 'compressing';
			photos = [...photos]; // Trigger reactivity
			let blob: Blob;
			let contentType: string;
			let outName: string;
			try {
				const result = await compressToUnder2MB(file);
				blob = result.blob;
				contentType = result.contentType;
				outName = result.outName;
			} catch (err) {
				item.status = 'error';
				item.error = err instanceof Error ? err.message : 'Compression failed';
				error = `Upload error for ${item.fileName}: ${item.error}`;
				photos = [...photos]; // Trigger reactivity
				continue;
			}
			item.status = 'uploading';
			photos = [...photos]; // Trigger reactivity
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const presignRes = await (fetchData as any)('/api/uploads/presign', {
				method: 'post',
				body: { filename: outName, contentType, sizeBytes: blob.size }
			});
			if (!presignRes.ok) {
				item.status = 'error';
				item.error = presignRes.error?.message || 'Failed to presign';
				error = `Upload error for ${item.fileName}: ${item.error}`;
				photos = [...photos]; // Trigger reactivity
				continue;
			}
			const { url, headers, key } = presignRes.data as {
				url: string;
				headers: Record<string, string>;
				key: string;
			};
			const putOk = await putToS3(url, headers, blob);
			if (!putOk) {
				item.status = 'error';
				item.error = 'Upload failed';
				error = `Upload error for ${item.fileName}: ${item.error}`;
				photos = [...photos]; // Trigger reactivity
				continue;
			}
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const regRes = await (fetchData as any)('/api/uploads/register', {
				method: 'post',
				body: { key }
			});
			if (!regRes.ok) {
				item.status = 'error';
				item.error = regRes.error?.message || 'Register failed';
				photos = [...photos]; // Trigger reactivity
				continue;
			}
			item.status = 'uploaded';
			item.key = key;
			uploadedPhotoKeys = [...uploadedPhotoKeys, key];
			uploadTimestamps[key] = Date.now(); // Track upload time for 1-hour cleanup
			photos = [...photos]; // Trigger reactivity
		}
	}

	async function putToS3(
		url: string,
		headers: Record<string, string>,
		blob: Blob
	): Promise<boolean> {
		try {
			const res = await fetch(url, {
				method: 'PUT',
				body: blob,
				headers
			});
			return res.ok;
		} catch {
			return false;
		}
	}

	async function blobFromImage(bitmap: ImageBitmap, type: string, quality: number): Promise<Blob> {
		const canvas = document.createElement('canvas');
		canvas.width = bitmap.width;
		canvas.height = bitmap.height;
		const ctx = canvas.getContext('2d');
		if (!ctx) throw new Error('Canvas unsupported');
		ctx.drawImage(bitmap, 0, 0);
		const b = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, type, quality));
		if (!b) throw new Error('Failed to create blob');
		return b;
	}

	async function compressToUnder2MB(
		file: File
	): Promise<{ blob: Blob; contentType: string; outName: string }> {
		const targetType = 'image/webp';
		let bitmap = await createImageBitmap(file);
		const maxDim = 3000;
		if (bitmap.width > maxDim || bitmap.height > maxDim) {
			const scale = Math.min(maxDim / bitmap.width, maxDim / bitmap.height);
			bitmap = await downscaleBitmap(bitmap, scale);
		}
		let quality = 0.92;
		let blob = await blobFromImage(bitmap, targetType, quality);
		let attempts = 0;
		while (blob.size > MAX_BYTES && attempts < 6) {
			quality = Math.max(0.4, quality - 0.15);
			blob = await blobFromImage(bitmap, targetType, quality);
			attempts++;
		}
		if (blob.size > MAX_BYTES) {
			for (let i = 0; i < 3 && blob.size > MAX_BYTES; i++) {
				bitmap = await downscaleBitmap(bitmap, 0.8);
				quality = Math.max(0.5, quality - 0.1);
				blob = await blobFromImage(bitmap, targetType, quality);
			}
		}
		if (blob.size > MAX_BYTES) {
			throw new Error('Unable to compress under 2MB');
		}
		const outName = file.name.replace(/\.[^.]+$/, '') + '.webp';
		return { blob, contentType: targetType, outName };
	}

	async function downscaleBitmap(src: ImageBitmap, scale: number): Promise<ImageBitmap> {
		const canvas = document.createElement('canvas');
		canvas.width = Math.max(1, Math.floor(src.width * scale));
		canvas.height = Math.max(1, Math.floor(src.height * scale));
		const ctx = canvas.getContext('2d');
		if (!ctx) throw new Error('Canvas unsupported');
		ctx.imageSmoothingQuality = 'high';
		ctx.drawImage(src, 0, 0, canvas.width, canvas.height);
		const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve));
		if (!blob) throw new Error('Downscale failed');
		return await createImageBitmap(blob);
	}

	function initializeFormData(): FormData {
		if (!plant) return createEmptyFormData();
		return {
			// Basic info
			name: plant.name,
			species: plant.species,
			isToxic: plant.isToxic,
			sunlight: plant.sunlight as SunlightRequirement,
			preferedTemperature: plant.preferedTemperature,

			// Location
			room: plant.location?.room ?? 'Unknown',
			position: plant.location?.position ?? 'Unknown',
			isOutdoors: plant.location?.isOutdoors ?? false,

			// Watering
			wateringIntervalDays: plant.watering?.intervalDays ?? 7,
			wateringMethod: plant.watering?.method ?? WateringMethod.Top,
			waterType: plant.watering?.waterType ?? WaterType.Tap,

			// Fertilizing
			fertilizingType: plant.fertilizing?.type ?? FertilizerType.Liquid,
			fertilizingIntervalDays: plant.fertilizing?.intervalDays ?? 30,
			npkRatio: plant.fertilizing?.npkRatio ?? '10:10:10',
			concentrationPercent: plant.fertilizing?.concentrationPercent ?? 50,
			activeInWinter: plant.fertilizing?.activeInWinter ?? false,

			// Humidity
			targetHumidity: plant.humidity?.targetHumidityPct ?? 50,
			requiresMisting: plant.humidity?.requiresMisting ?? false,
			mistingIntervalDays: plant.humidity?.mistingIntervalDays ?? 3,
			requiresHumidifier: plant.humidity?.requiresHumidifier ?? false,

			// Soil
			soilType: plant.soil?.type ?? 'Generic',
			repottingCycle: plant.soil?.repottingCycle ?? 2,
			soilComponents: plant.soil?.components ?? [],

			// Seasonality
			winterRestPeriod: plant.seasonality?.winterRestPeriod ?? false,
			winterWaterFactor: plant.seasonality?.winterWaterFactor ?? 0.5,
			minTempCelsius: plant.seasonality?.minTempCelsius ?? 15,

			// Metadata
			flags: plant.flags ?? [],
			notes: plant.notes ?? []
		};
	}

	async function submitForm(): Promise<void> {
		if (!formData.species.trim() || !formData.name.trim()) {
			error = 'Species and name are required';
			return;
		}

		submitting = true;
		error = null;
		success = null;

		try {
			// Combine existing photoIds (excluding removed) with newly uploaded keys
			const existingPhotoIds = (plant?.photoIds || []).filter(
				(id) => !removedPhotoIds.includes(id)
			);
			const allPhotoIds = [...existingPhotoIds, ...uploadedPhotoKeys];

			const updatePayload = {
				name: formData.name,
				species: formData.species,
				isToxic: formData.isToxic,
				sunlight: formData.sunlight,
				preferedTemperature: formData.preferedTemperature,
				location: {
					room: formData.room,
					position: formData.position,
					isOutdoors: formData.isOutdoors
				},
				watering: {
					intervalDays: formData.wateringIntervalDays,
					method: formData.wateringMethod,
					waterType: formData.waterType
				},
				fertilizing: {
					type: formData.fertilizingType,
					intervalDays: formData.fertilizingIntervalDays,
					npkRatio: formData.npkRatio,
					concentrationPercent: formData.concentrationPercent,
					activeInWinter: formData.activeInWinter
				},
				humidity: {
					targetHumidityPct: formData.targetHumidity,
					requiresMisting: formData.requiresMisting,
					mistingIntervalDays: formData.mistingIntervalDays,
					requiresHumidifier: formData.requiresHumidifier
				},
				soil: {
					type: formData.soilType,
					repottingCycle: formData.repottingCycle,
					components: formData.soilComponents
				},
				seasonality: {
					winterRestPeriod: formData.winterRestPeriod,
					winterWaterFactor: formData.winterWaterFactor,
					minTempCelsius: formData.minTempCelsius
				},
				flags: formData.flags,
				notes: formData.notes,
				photoIds: allPhotoIds
			};

			const res = await fetchData('/api/plants/{id}', {
				method: 'patch',
				params: { id: plant?.id ?? '' },
				body: updatePayload
			});

			if (!res.ok) {
				throw new Error(res.error?.message || 'Failed to update plant');
			}

			success = 'Plant updated successfully!';
			if (res.data) {
				plant = res.data;
			}
			// Clear uploaded photos since they were successfully applied
			uploadedPhotoKeys = [];
			uploadTimestamps = {};

			// Invalidate cache for this plant and the plants list
			const plantId = plant?.id ?? '';
			if (navigator.serviceWorker?.controller && plantId) {
				navigator.serviceWorker.controller.postMessage({
					type: 'INVALIDATE_CACHE',
					urls: [`/api/plants/${plantId}`, '/api/plants']
				});
			}

			setTimeout(() => goto(resolve('/manage')), 1500);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			submitting = false;
		}
	}

	function removeExistingPhoto(photoId: string, index: number): void {
		if (!confirm('Remove this photo?')) return;
		removedPhotoIds = [...removedPhotoIds, photoId];
		const newUrls = [...previewUrls];
		const urlToRevoke = newUrls[index];
		newUrls.splice(index, 1);
		previewUrls = newUrls;
		if (urlToRevoke) revokeObjectURL(urlToRevoke);
		// Delete from S3
		deletePhotoFromS3(photoId);
	}

	async function deletePhotoFromS3(key: string): Promise<void> {
		try {
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			await (fetchData as any)(`/api/uploads/${encodeURIComponent(key)}`, {
				method: 'delete'
			});
		} catch (err) {
			console.error('Failed to delete photo from S3:', err);
		}
	}

	async function cleanupUnappliedUploads(): Promise<void> {
		// Delete any uploaded photos that were never saved to the plant
		// This handles the case where user closes/navigates away without saving
		if (uploadedPhotoKeys.length === 0) return;

		for (const key of uploadedPhotoKeys) {
			await deletePhotoFromS3(key);
		}
	}

	function handleBackClick(): void {
		cleanupUnappliedUploads();
		goto(resolve('/manage'));
	}

	function resetForm(): void {
		formData = initializeFormData();
		error = null;
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-emerald-50 via-green-50 to-teal-100 p-6 md:p-10">
	<div class="mx-auto max-w-4xl">
		{#if loading}
			<div class="flex min-h-screen items-center justify-center">
				<div class="text-center">
					<div class="mb-4 animate-spin text-4xl">🌱</div>
					<p class="text-lg text-gray-600">Loading plant details...</p>
				</div>
			</div>
		{:else if !plant}
			<div class="flex min-h-screen items-center justify-center">
				<div class="text-center">
					<p class="mb-4 text-lg text-red-600">{error || 'Plant not found'}</p>
					<a
						href={resolve('/manage')}
						onclick={(e) => {
							e.preventDefault();
							handleBackClick();
						}}
						class="rounded-xl bg-gray-600 px-4 py-2 font-medium text-white shadow-sm transition hover:bg-gray-700"
					>
						← Back to Plants
					</a>
				</div>
			</div>
		{:else}
			<div class="mb-8">
				<div class="mb-4 flex items-center justify-between">
					<div>
						<h1 class="flex items-center gap-3 text-4xl font-bold text-green-900">
							🌿 {plant?.name || 'Plant'}
						</h1>
						<p class="mt-1 text-sm text-emerald-700 italic">{plant?.species || ''}</p>
					</div>
					<a
						href={resolve('/manage')}
						onclick={(e) => {
							e.preventDefault();
							handleBackClick();
						}}
						class="rounded-xl bg-gray-600 px-4 py-2 font-medium text-white shadow-sm transition hover:bg-gray-700"
					>
						← Back
					</a>
				</div>
			</div>

			<!-- Messages -->
			{#if success}
				<div
					class="mb-6 rounded-lg border-2 border-green-400 bg-green-100 px-6 py-4 text-green-800"
				>
					✓ {success}
				</div>
			{/if}

			{#if error}
				<div class="mb-6 rounded-lg border-2 border-red-400 bg-red-100 px-6 py-4 text-red-800">
					✕ {error}
				</div>
			{/if}

			<div class="space-y-6">
				<!-- Images Section -->
				<div class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur">
					<h2 class="mb-4 text-2xl font-bold text-green-800">📸 Photos</h2>
					<div class="space-y-4">
						<label class="block">
							<span class="text-sm font-medium text-green-800"
								>Add new images (JPEG/PNG/WebP, auto-compressed ≤ 2MB)</span
							>
							<input
								type="file"
								accept="image/jpeg,image/png,image/webp"
								multiple
								onchange={onFilesSelected}
								class="mt-2 w-full rounded-lg border border-emerald-200 bg-white p-2 text-sm"
							/>
						</label>

						{#if photos.length}
							<div>
								<p class="mb-2 text-sm font-medium text-green-800">New uploads:</p>
								<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
									{#each photos as p (p.previewUrl)}
										<div class="rounded-md border border-emerald-200 bg-emerald-50 p-2">
											<img
												src={p.previewUrl}
												alt={p.fileName}
												class="h-24 w-full rounded object-cover"
											/>
											<div class="mt-1 text-xs text-emerald-800">
												{p.fileName}
											</div>
											<div class="text-xs">
												{#if p.status === 'pending'}
													<span class="text-gray-600">⏸️ Pending</span>
												{:else if p.status === 'compressing'}
													<span class="text-blue-600">⚙️ Compressing...</span>
												{:else if p.status === 'uploading'}
													<span class="text-emerald-600">📤 Uploading...</span>
												{:else if p.status === 'uploaded'}
													<span class="font-semibold text-green-700">✓ Uploaded!</span>
												{:else}
													<span class="text-red-600">✕ {p.error || 'Error'}</span>
												{/if}
											</div>
										</div>
									{/each}
								</div>
							</div>
						{/if}

						{#if previewUrls.length}
							<div>
								<p class="mb-2 text-sm font-medium text-green-800">Existing photos:</p>
								<div class="grid grid-cols-2 gap-3 md:grid-cols-3">
									{#each previewUrls as u, i (u)}
										<div class="group relative">
											<img src={u} alt="" class="h-32 w-full rounded object-cover" />
											<button
												type="button"
												onclick={() => removeExistingPhoto(plant?.photoIds?.[i] ?? '', i)}
												class="absolute top-1 right-1 flex h-6 w-6 items-center justify-center rounded-full bg-red-600 text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100 hover:bg-red-700"
												title="Remove photo"
											>
												×
											</button>
										</div>
									{/each}
								</div>
							</div>
						{:else if !photos.length}
							<div
								class="flex h-48 items-center justify-center rounded-lg border-2 border-dashed border-emerald-300 bg-emerald-50"
							>
								<div class="text-center">
									<div class="mb-2 text-4xl">🖼️</div>
									<p class="text-sm text-emerald-700">No photos yet</p>
								</div>
							</div>
						{/if}
					</div>
				</div>

				<!-- Form Sections -->
				<BasicInformationForm {formData} />
				<LocationForm {formData} />

				<!-- Watering & Fertilizing -->
				<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
					<WateringForm {formData} />
					<FertilizingForm {formData} />
				</div>

				<MistingForm {formData} />

				<!-- Advanced Settings -->
				<SoilForm {formData} bind:soilComponentInput />
				<SeasonalityForm {formData} />
				<MetadataForm {formData} bind:newNote />

				<!-- Action Buttons -->
				<div class="flex justify-between gap-3">
					<button
						onclick={resetForm}
						class="rounded-lg bg-gray-200 px-6 py-3 font-semibold text-gray-800 transition hover:bg-gray-300"
					>
						Reset
					</button>
					<button
						onclick={submitForm}
						disabled={submitting}
						class="rounded-lg bg-gradient-to-r from-emerald-600 to-green-600 px-8 py-3 font-semibold text-white shadow-md transition hover:from-emerald-700 hover:to-green-700 disabled:opacity-50"
					>
						{submitting ? 'Saving...' : 'Save Changes'}
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
</style>
