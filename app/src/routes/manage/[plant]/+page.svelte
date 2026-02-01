<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { Plant } from '$lib/types/api';
	import { SunlightRequirement, WateringMethod, WaterType, FertilizerType } from '$lib/types/api';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/stores';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import { getImageObjectURL, revokeObjectURL } from '$lib/utils/imageCache';
	import { invalidateApiCache } from '$lib/utils/cache';
	import { tStore } from '$lib/i18n';
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
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import LoadingSpinner from '$lib/components/ui/LoadingSpinner.svelte';
	import Alert from '$lib/components/ui/Message.svelte';

	let plant = $state<Plant | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let success = $state<string | null>(null);
	let submitting = $state(false);
	let newNote = $state('');
	let soilComponentInput = $state('');

	// Accordion state for mobile
	let expandedSections = $state<Record<string, boolean>>({
		basic: true,
		location: false,
		watering: false,
		fertilizing: false,
		humidity: false,
		soil: false,
		seasonality: false,
		metadata: false,
		photos: true
	});
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
	let originalFormData = $state<FormData>(createEmptyFormData());

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
			originalFormData = JSON.parse(JSON.stringify(formData));
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

			// Helper to check if value changed
			const hasChanged = (key: keyof FormData): boolean => {
				return JSON.stringify(formData[key]) !== JSON.stringify(originalFormData[key]);
			};

			// Build updatePayload with only changed fields
			type UpdatePayload = Record<string, unknown>;
			const updatePayload: UpdatePayload = {};

			// Basic info
			if (hasChanged('name')) updatePayload.name = formData.name;
			if (hasChanged('species')) updatePayload.species = formData.species;
			if (hasChanged('isToxic')) updatePayload.isToxic = formData.isToxic;
			if (hasChanged('sunlight')) updatePayload.sunlight = formData.sunlight;
			if (hasChanged('preferedTemperature'))
				updatePayload.preferedTemperature = formData.preferedTemperature;

			// Location
			if (hasChanged('room') || hasChanged('position') || hasChanged('isOutdoors')) {
				updatePayload.location = {
					room: formData.room,
					position: formData.position,
					isOutdoors: formData.isOutdoors
				};
			}

			// Watering
			if (
				hasChanged('wateringIntervalDays') ||
				hasChanged('wateringMethod') ||
				hasChanged('waterType')
			) {
				updatePayload.watering = {
					intervalDays: formData.wateringIntervalDays,
					method: formData.wateringMethod,
					waterType: formData.waterType
				};
			}

			// Fertilizing
			if (
				hasChanged('fertilizingType') ||
				hasChanged('fertilizingIntervalDays') ||
				hasChanged('npkRatio') ||
				hasChanged('concentrationPercent') ||
				hasChanged('activeInWinter')
			) {
				updatePayload.fertilizing = {
					type: formData.fertilizingType,
					intervalDays: formData.fertilizingIntervalDays,
					npkRatio: formData.npkRatio,
					concentrationPercent: formData.concentrationPercent,
					activeInWinter: formData.activeInWinter
				};
			}

			// Humidity
			if (
				hasChanged('targetHumidity') ||
				hasChanged('requiresMisting') ||
				hasChanged('mistingIntervalDays') ||
				hasChanged('requiresHumidifier')
			) {
				updatePayload.humidity = {
					targetHumidityPct: formData.targetHumidity,
					requiresMisting: formData.requiresMisting,
					mistingIntervalDays: formData.mistingIntervalDays,
					requiresHumidifier: formData.requiresHumidifier
				};
			}

			// Soil
			if (hasChanged('soilType') || hasChanged('repottingCycle') || hasChanged('soilComponents')) {
				updatePayload.soil = {
					type: formData.soilType,
					repottingCycle: formData.repottingCycle,
					components: formData.soilComponents
				};
			}

			// Seasonality
			if (
				hasChanged('winterRestPeriod') ||
				hasChanged('winterWaterFactor') ||
				hasChanged('minTempCelsius')
			) {
				updatePayload.seasonality = {
					winterRestPeriod: formData.winterRestPeriod,
					winterWaterFactor: formData.winterWaterFactor,
					minTempCelsius: formData.minTempCelsius
				};
			}

			// Metadata
			if (hasChanged('flags')) updatePayload.flags = formData.flags;
			if (hasChanged('notes')) updatePayload.notes = formData.notes;

			// Photos - only send if changed
			if (uploadedPhotoKeys.length > 0 || removedPhotoIds.length > 0) {
				updatePayload.photoIds = allPhotoIds;
			}

			console.log('Sending update payload:', updatePayload);

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
				formData = initializeFormData();
				originalFormData = JSON.parse(JSON.stringify(formData));
			}
			// Clear uploaded photos since they were successfully applied
			uploadedPhotoKeys = [];
			uploadTimestamps = {};
			removedPhotoIds = [];

			// Invalidate cache for this plant and the plants list BEFORE redirecting
			const plantId = plant?.id ?? '';
			if (plantId) {
				await invalidateApiCache([`/api/plants/${plantId}`, '/api/plants'], {
					waitForAck: true,
					timeoutMs: 100
				});
			}

			goto(resolve(plant ? `/plant/${plant.id}` : '/'));
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			submitting = false;
		}
	}

	function removeExistingPhoto(photoId: string, index: number): void {
		if (!confirm($tStore('plants.deletePhotoConfirm') || 'Remove this photo?')) return;
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
		goto(resolve(plant ? `/plant/${plant.id}` : '/'));
	}

	function toggleSection(section: string): void {
		expandedSections[section] = !expandedSections[section];
	}

	function resetForm(): void {
		formData = initializeFormData();
		originalFormData = JSON.parse(JSON.stringify(formData));
		error = null;
	}
</script>

<div class="flex h-full min-h-0 flex-col overflow-hidden">
	<div class="flex-shrink-0">
		<PageHeader
			icon="✏️"
			title={plant?.name || $tStore('plants.editPlant')}
			description={plant?.species || ''}
		>
			<Button
				variant="ghost"
				size="sm"
				text="common.back"
				icon="←"
				onclick={() => handleBackClick()}
			/>
		</PageHeader>
	</div>

	<PageContent>
		{#if loading}
			<LoadingSpinner message="Loading plant details..." icon="🌱" />
		{:else if !plant}
			<div class="flex flex-col items-center justify-center gap-6 py-12">
				<p class="text-lg text-red-600">{error || 'Plant not found'}</p>
				<Button variant="secondary" onclick={() => handleBackClick()} text="common.back" />
			</div>
		{:else}
			<div class="h-full min-h-0 flex-1 overflow-y-auto pb-24">
				<!-- Messages -->
				{#if success}
					<Alert type="success" title="common.success" description={success} />
				{/if}

				{#if error}
					<Alert type="error" title="common.error" description={error} />
				{/if}

				<div class="space-y-3 px-4 py-4">
					<!-- Photos Section - Expanded by default -->
					<div class="overflow-hidden rounded-lg border border-emerald-200 bg-white">
						<button
							onclick={() => toggleSection('photos')}
							aria-expanded={expandedSections.photos}
							aria-label={`${$tStore('plants.photos')} section, ${expandedSections.photos ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-gradient-to-r from-emerald-50 to-emerald-100/50 px-4 py-3 text-left font-semibold text-emerald-900 transition-colors hover:bg-emerald-100 focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:outline-none"
						>
							<span class="flex items-center gap-2">
								<span>📸 {$tStore('plants.photos')}</span>
								<span class="rounded-full bg-emerald-200 px-2 py-0.5 text-xs font-normal">
									{previewUrls.length + photos.filter((p) => p.status === 'uploaded').length}
								</span>
							</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.photos ? '−' : '+'}</span>
						</button>

						{#if expandedSections.photos}
							<div class="space-y-4 border-t border-emerald-100 p-4">
								<label class="block">
									<span class="text-sm font-medium text-emerald-900"
										>{$tStore('plants.addImages')}</span
									>
									<input
										type="file"
										accept="image/jpeg,image/png,image/webp"
										multiple
										onchange={onFilesSelected}
										aria-label={$tStore('plants.addImages')}
										class="mt-2 w-full cursor-pointer touch-manipulation rounded-lg border border-emerald-300 bg-white p-3 text-sm"
									/>
								</label>

								{#if photos.length}
									<div>
										<p class="mb-2 text-xs font-semibold text-emerald-700">
											{$tStore('plants.newUploads')}
										</p>
										<div class="grid grid-cols-2 gap-2">
											{#each photos as p (p.previewUrl)}
												<div
													class="overflow-hidden rounded-lg border border-emerald-200 bg-emerald-50"
												>
													<img
														src={p.previewUrl}
														alt={p.fileName}
														class="h-20 w-full object-cover"
													/>
													<div class="p-2 text-xs">
														<div class="mb-1 truncate font-medium text-emerald-900">
															{p.fileName}
														</div>
														<div class="text-xs">
															{#if p.status === 'pending'}
																<span class="text-gray-600">⏸️ Pending</span>
															{:else if p.status === 'compressing'}
																<span class="text-blue-600">⚙️ Compressing</span>
															{:else if p.status === 'uploading'}
																<span class="text-emerald-600">📤 Uploading</span>
															{:else if p.status === 'uploaded'}
																<span class="font-semibold text-green-700">✓ Uploaded</span>
															{:else}
																<span class="text-red-600">✕ {p.error || 'Error'}</span>
															{/if}
														</div>
													</div>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if previewUrls.length}
									<div>
										<p class="mb-2 text-xs font-semibold text-emerald-700">
											{$tStore('plants.existingPhotos')}
										</p>
										<div class="grid grid-cols-2 gap-2">
											{#each previewUrls as u, i (u)}
												<div
													class="group relative overflow-hidden rounded-lg border border-emerald-200"
												>
													<img
														src={u}
														alt={plant?.name || 'Plant'}
														class="h-20 w-full object-cover"
													/>
													<button
														onclick={() => removeExistingPhoto(plant?.photoIds?.[i] ?? '', i)}
														aria-label={$tStore('plants.deletePhoto')}
														class="absolute inset-0 flex cursor-pointer items-center justify-center bg-red-600/80 text-sm font-bold text-white opacity-0 transition-all group-hover:opacity-100 hover:bg-red-700/90 focus:opacity-100 focus:outline-none"
													>
														{$tStore('common.delete')}
													</button>
												</div>
											{/each}
										</div>
									</div>
								{:else if !photos.length}
									<div
										class="flex h-24 items-center justify-center rounded-lg border-2 border-dashed border-emerald-300 bg-emerald-50"
									>
										<div class="text-center">
											<div class="text-2xl">🖼️</div>
											<p class="text-xs text-emerald-700">{$tStore('plants.noPhotosYet')}</p>
										</div>
									</div>
								{/if}
							</div>
						{/if}
					</div>

					<!-- Accordion Sections -->
					<!-- Basic Information -->
					<div class="overflow-hidden rounded-lg border border-gray-200 bg-white">
						<button
							onclick={() => toggleSection('basic')}
							aria-expanded={expandedSections.basic}
							aria-label={`${$tStore('plants.basicInformation')} section, ${expandedSections.basic ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-gray-50 px-4 py-3 text-left font-semibold text-gray-900 transition-colors hover:bg-gray-100 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>📋 {$tStore('plants.basicInformation')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.basic ? '−' : '+'}</span>
						</button>
						{#if expandedSections.basic}
							<div class="border-t border-gray-100 p-4">
								<BasicInformationForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Location -->
					<div class="overflow-hidden rounded-lg border border-gray-200 bg-white">
						<button
							onclick={() => toggleSection('location')}
							aria-expanded={expandedSections.location}
							aria-label={`${$tStore('plants.location')} section, ${expandedSections.location ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-gray-50 px-4 py-3 text-left font-semibold text-gray-900 transition-colors hover:bg-gray-100 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>📍 {$tStore('plants.location')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.location ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.location}
							<div class="border-t border-gray-100 p-4">
								<LocationForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Watering -->
					<div class="overflow-hidden rounded-lg border border-blue-200 bg-white">
						<button
							onclick={() => toggleSection('watering')}
							aria-expanded={expandedSections.watering}
							aria-label={`${$tStore('plants.wateringTitle')} section, ${expandedSections.watering ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-blue-50 px-4 py-3 text-left font-semibold text-blue-900 transition-colors hover:bg-blue-100 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>💧 {$tStore('plants.wateringTitle')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.watering ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.watering}
							<div class="border-t border-blue-100 p-4">
								<WateringForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Fertilizing -->
					<div class="overflow-hidden rounded-lg border border-yellow-200 bg-white">
						<button
							onclick={() => toggleSection('fertilizing')}
							aria-expanded={expandedSections.fertilizing}
							aria-label={`${$tStore('plants.fertilizingTitle')} section, ${expandedSections.fertilizing ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-yellow-50 px-4 py-3 text-left font-semibold text-yellow-900 transition-colors hover:bg-yellow-100 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>🍯 {$tStore('plants.fertilizingTitle')}</span>
							<span class="text-lg" aria-hidden="true"
								>{expandedSections.fertilizing ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.fertilizing}
							<div class="border-t border-yellow-100 p-4">
								<FertilizingForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Humidity & Misting -->
					<div class="overflow-hidden rounded-lg border border-purple-200 bg-white">
						<button
							onclick={() => toggleSection('humidity')}
							aria-expanded={expandedSections.humidity}
							aria-label={`${$tStore('plants.humidityTitle')} section, ${expandedSections.humidity ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-purple-50 px-4 py-3 text-left font-semibold text-purple-900 transition-colors hover:bg-purple-100 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>💨 {$tStore('plants.humidityTitle')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.humidity ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.humidity}
							<div class="space-y-4 border-t border-purple-100 p-4">
								<MistingForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Soil -->
					<div class="overflow-hidden rounded-lg border border-amber-200 bg-white">
						<button
							onclick={() => toggleSection('soil')}
							aria-expanded={expandedSections.soil}
							aria-label={`${$tStore('plants.soilTitle')} section, ${expandedSections.soil ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-amber-50 px-4 py-3 text-left font-semibold text-amber-900 transition-colors hover:bg-amber-100 focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>🌍 {$tStore('plants.soilTitle')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.soil ? '−' : '+'}</span>
						</button>
						{#if expandedSections.soil}
							<div class="border-t border-amber-100 p-4">
								<SoilForm {formData} bind:soilComponentInput />
							</div>
						{/if}
					</div>

					<!-- Seasonality -->
					<div class="overflow-hidden rounded-lg border border-orange-200 bg-white">
						<button
							onclick={() => toggleSection('seasonality')}
							aria-expanded={expandedSections.seasonality}
							aria-label={`${$tStore('plants.seasonalityTitle')} section, ${expandedSections.seasonality ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-orange-50 px-4 py-3 text-left font-semibold text-orange-900 transition-colors hover:bg-orange-100 focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>❄️ {$tStore('plants.seasonalityTitle')}</span>
							<span class="text-lg" aria-hidden="true"
								>{expandedSections.seasonality ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.seasonality}
							<div class="border-t border-orange-100 p-4">
								<SeasonalityForm {formData} />
							</div>
						{/if}
					</div>

					<!-- Metadata -->
					<div class="overflow-hidden rounded-lg border border-red-200 bg-white">
						<button
							onclick={() => toggleSection('metadata')}
							aria-expanded={expandedSections.metadata}
							aria-label={`${$tStore('plants.metadata')} section, ${expandedSections.metadata ? 'expanded' : 'collapsed'}`}
							class="flex w-full cursor-pointer items-center justify-between bg-red-50 px-4 py-3 text-left font-semibold text-red-900 transition-colors hover:bg-red-100 focus:ring-2 focus:ring-red-500 focus:ring-offset-2 focus:outline-none"
						>
							<span>🏷️ {$tStore('plants.metadata')}</span>
							<span class="text-lg" aria-hidden="true">{expandedSections.metadata ? '−' : '+'}</span
							>
						</button>
						{#if expandedSections.metadata}
							<div class="border-t border-red-100 p-4">
								<MetadataForm {formData} bind:newNote />
							</div>
						{/if}
					</div>

					<!-- Action Buttons -->
					<div
						class="sticky bottom-0 flex gap-3 border-t border-gray-200 bg-gradient-to-t from-white via-white to-white/80 px-4 py-4"
					>
						<Button
							variant="secondary"
							size="md"
							onclick={resetForm}
							text="common.reset"
							class="w-full cursor-pointer"
						/>
						<Button
							variant="primary"
							size="md"
							disabled={submitting}
							onclick={submitForm}
							text={submitting ? 'common.saving' : 'common.save'}
							class="w-full cursor-pointer"
						/>
					</div>
				</div>
			</div>
		{/if}
	</PageContent>
</div>
