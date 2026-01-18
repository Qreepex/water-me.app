<script lang="ts">
	import { onMount } from 'svelte';
	import type { Plant } from '$lib/types/types';
	import { PlantFlag, SunlightRequirement } from '$lib/types/types';
	import { goto } from '$app/navigation';
	import { API_BASE_URL } from '$lib/constants';
	import { resolve } from '$app/paths';

	interface FormData {
		id?: string;
		species: string;
		name: string;
		sunLight: SunlightRequirement;
		preferedTemperature: number;
		wateringIntervalDays: number;
		fertilizingIntervalDays: number;
		preferedHumidity: number;
		sprayIntervalDays?: number;
		notes: string[];
		flags: PlantFlag[];
		photoIds: string[];
	}

	let plants: Plant[] = [];
	let loading = true;
	let error: string | null = null;
	let success: string | null = null;
	let submitting = false;
	let token: string | null = null;
	let isInitialized = false;

	// authStore.subscribe((state) => {
	// 	token = state.token;
	// 	isInitialized = state.initialized;
	// });

	let showForm = false;
	let editingId: string | null = null;

	let formData: FormData = {
		species: '',
		name: '',
		sunLight: SunlightRequirement.INDIRECT_SUN,
		preferedTemperature: 20,
		wateringIntervalDays: 7,
		fertilizingIntervalDays: 30,
		preferedHumidity: 50,
		notes: [],
		flags: [],
		photoIds: []
	};

	let newNote = '';
	let photoPreview: { data: string; key?: string }[] = []; // data URL or S3 key
	let fileInput: HTMLInputElement;
	let photoUploadProgress: Record<number, number> = {};

	function handleLogout() {
		// authStore.logout();
		goto(resolve('/'));
	}

	// Compress image to max 2MB using canvas
	async function compressImage(file: File): Promise<Blob> {
		return new Promise((resolve, reject) => {
			const reader = new FileReader();
			reader.onload = (e) => {
				const img = new Image();
				img.onload = () => {
					const canvas = document.createElement('canvas');
					let { width, height } = img;

					// Scale down if too large
					const maxDimension = 2048;
					if (width > maxDimension || height > maxDimension) {
						const scale = Math.min(maxDimension / width, maxDimension / height);
						width = Math.floor(width * scale);
						height = Math.floor(height * scale);
					}

					canvas.width = width;
					canvas.height = height;
					const ctx = canvas.getContext('2d')!;
					ctx.drawImage(img, 0, 0, width, height);

					// Compress with decreasing quality until under 2MB
					let quality = 0.9;
					const targetSize = 2 * 1024 * 1024; // 2MB
					canvas.toBlob(
						(blob) => {
							if (blob && blob.size <= targetSize) {
								resolve(blob);
							} else if (quality > 0.1) {
								quality -= 0.1;
								canvas.toBlob((b) => b && resolve(b), 'image/jpeg', quality);
							} else {
								resolve(blob || new Blob());
							}
						},
						'image/jpeg',
						quality
					);
				};
				img.onerror = () => reject(new Error('Failed to load image'));
				img.src = e.target?.result as string;
			};
			reader.onerror = () => reject(new Error('Failed to read file'));
			reader.readAsDataURL(file);
		});
	}

	// Request presigned URLs from backend
	async function getPresignedUrls(
		files: File[]
	): Promise<{ key: string; url: string; headers: Record<string, string> }[]> {
		if (!token) throw new Error('Unauthorized');

		const body = {
			files: files.map((f) => ({
				filename: f.name,
				contentType: f.type
			}))
		};

		const response = await fetch(API_BASE_URL + '/api/uploads/presign', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(body)
		});

		if (!response.ok) {
			const data = await response.json();
			throw new Error(data.error || `HTTP ${response.status}`);
		}

		const result = await response.json();
		return result.items;
	}

	// Upload file to S3 via presigned URL
	async function uploadToS3(
		blob: Blob,
		url: string,
		headers: Record<string, string>,
		index: number
	): Promise<string> {
		return new Promise((resolve, reject) => {
			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (e) => {
				if (e.lengthComputable) {
					photoUploadProgress[index] = Math.round((e.loaded / e.total) * 100);
				}
			});

			xhr.addEventListener('load', () => {
				if (xhr.status >= 200 && xhr.status < 300) {
					resolve(xhr.responseText || '');
				} else {
					reject(new Error(`Upload failed: ${xhr.status}`));
				}
			});

			xhr.addEventListener('error', () => reject(new Error('Network error during upload')));

			xhr.open('PUT', url);
			Object.entries(headers).forEach(([key, value]) => xhr.setRequestHeader(key, value));
			xhr.send(blob);
		});
	}

	onMount(() => {
		// Wait for auth to initialize
		const checkAuth = setInterval(() => {
			if (isInitialized) {
				clearInterval(checkAuth);
				if (!token) {
					goto(resolve('/'));
					return;
				}
				loadPlants();
			}
		}, 50);
	});

	async function loadPlants(): Promise<void> {
		try {
			const response = await fetch(API_BASE_URL + '/api/plants', {
				headers: { Authorization: `Bearer ${token}` }
			});

			if (response.status === 401) {
				// authStore.logout();
				goto(resolve('/'));
				return;
			}

			if (!response.ok) throw new Error('Failed to fetch plants');
			plants = await response.json();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			loading = false;
		}
	}

	function resetForm(): void {
		formData = {
			species: '',
			name: '',
			sunLight: SunlightRequirement.INDIRECT_SUN,
			preferedTemperature: 20,
			wateringIntervalDays: 7,
			fertilizingIntervalDays: 30,
			preferedHumidity: 50,
			notes: [],
			flags: [],
			photoIds: []
		};
		newNote = '';
		photoPreview = [];
		editingId = null;
	}

	function startCreate(): void {
		resetForm();
		showForm = true;
		editingId = null;
	}

	function startEdit(plant: Plant): void {
		formData = {
			...plant,
			photoIds: plant.photoIds // S3 keys from DB
		};
		photoPreview = plant.photoIds.map((key) => ({ data: '', key })) as {
			data: string;
			key?: string;
		}[];
		showForm = true;
		editingId = plant.id;
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

	function toggleFlag(flag: PlantFlag): void {
		if (formData.flags.includes(flag)) {
			formData.flags = formData.flags.filter((f) => f !== flag);
		} else {
			formData.flags = [...formData.flags, flag];
		}
	}

	function handlePhotoUpload(event: Event): void {
		const target = event.target as HTMLInputElement;
		const files = target.files;
		if (!files) return;

		// Process files: compress and upload to S3
		const fileArray = Array.from(files);
		error = null;

		(async () => {
			try {
				// Compress all files
				const compressed: Blob[] = [];
				for (const file of fileArray) {
					const blob = await compressImage(file);
					if (blob.size > 2 * 1024 * 1024) {
						throw new Error(`${file.name} is too large even after compression`);
					}
					compressed.push(blob);
				}

				// Get presigned URLs
				success = `Uploading ${fileArray.length} photo(s)...`;
				const presigns = await getPresignedUrls(fileArray);

				// Upload each file to S3
				const uploadedKeys: string[] = [];
				for (let i = 0; i < compressed.length; i++) {
					await uploadToS3(compressed[i], presigns[i].url, presigns[i].headers, i);
					uploadedKeys.push(presigns[i].key);
					delete photoUploadProgress[i];
				}

				// Add to preview and formData
				photoPreview = [...photoPreview, ...uploadedKeys.map((key) => ({ data: '', key }))];
				formData.photoIds = [...formData.photoIds, ...uploadedKeys];
				success = `Uploaded ${uploadedKeys.length} photo(s) successfully!`;

				// Reset file input
				if (fileInput) fileInput.value = '';
			} catch (err) {
				error = err instanceof Error ? err.message : 'Upload failed';
				photoUploadProgress = {};
			}
		})();
	}

	function removePhoto(index: number): void {
		photoPreview = photoPreview.filter((_, i) => i !== index);
		formData.photoIds = photoPreview.map((p) => p.key || p.data).filter(Boolean);
	}

	async function submitForm(): Promise<void> {
		if (!token) {
			// authStore.logout();
			goto(resolve('/'));
			return;
		}

		if (!formData.species.trim() || !formData.name.trim()) {
			error = 'Species and name are required';
			return;
		}

		submitting = true;
		error = null;
		success = null;

		try {
			const url = editingId ? `/api/plants/${editingId}` : '/api/plants';
			const method = editingId ? 'PUT' : 'POST';

			const payload = {
				...formData,
				photoIds: formData.photoIds.filter(Boolean) // S3 keys
			};

			const response = await fetch(API_BASE_URL + url, {
				method,
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify(payload)
			});

			if (!response.ok) {
				const data = await response.json();
				throw new Error(data.error || `HTTP error! status: ${response.status}`);
			}

			success = editingId ? 'Plant updated successfully!' : 'Plant created successfully!';
			showForm = false;
			await loadPlants();
			resetForm();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			submitting = false;
		}
	}

	async function deletePlant(id: string): Promise<void> {
		if (!token || !confirm('Are you sure you want to delete this plant?')) return;

		try {
			const response = await fetch(API_BASE_URL + `/api/plants/${id}`, {
				method: 'DELETE',
				headers: { Authorization: `Bearer ${token}` }
			});
			if (!response.ok) throw new Error('Failed to delete plant');
			success = 'Plant deleted successfully!';
			await loadPlants();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		}
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-emerald-50 via-green-50 to-teal-100 p-6 md:p-10">
	<div class="mx-auto max-w-6xl">
		<!-- Header -->
		<div class="mb-8">
			<div class="mb-4 flex items-center justify-between">
				<h1 class="flex items-center gap-3 text-4xl font-bold text-green-900">🌿 Manage Plants</h1>
				<div class="flex items-center gap-3">
					<a
						href={resolve('/app')}
						class="rounded-xl bg-gradient-to-r from-green-600 to-emerald-600 px-4 py-2 font-medium text-white shadow-sm transition hover:from-green-700 hover:to-emerald-700"
					>
						← Back to Overview
					</a>
					<button
						on:click={handleLogout}
						class="rounded-xl bg-red-600 px-4 py-2 font-medium text-white shadow-sm transition hover:bg-red-700"
					>
						Logout
					</button>
				</div>
			</div>
			<p class="text-emerald-800">Create, update, or delete your plants</p>
		</div>

		<!-- Messages -->
		{#if success}
			<div class="mb-6 rounded-lg border-2 border-green-400 bg-green-100 px-6 py-4 text-green-800">
				{success}
			</div>
		{/if}

		{#if error}
			<div class="mb-6 rounded-lg border-2 border-red-400 bg-red-100 px-6 py-4 text-red-800">
				{error}
			</div>
		{/if}

		<div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
			<!-- Form Column -->
			<div class="lg:col-span-1">
				{#if !showForm}
					<button
						on:click={startCreate}
						class="w-full rounded-xl bg-green-600 px-6 py-4 text-lg font-bold text-white transition hover:bg-green-700"
					>
						+ Add New Plant
					</button>
				{:else}
					<div
						class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur"
					>
						<h2 class="mb-6 text-2xl font-bold text-green-800">
							{editingId ? 'Edit Plant' : 'New Plant'}
						</h2>

						<!-- Form Fields -->
						<div class="space-y-4">
							<!-- Name -->
							<div>
								<label for="plant-name" class="mb-1 block text-sm font-semibold text-gray-700"
									>Plant Name *</label
								>
								<input
									type="text"
									id="plant-name"
									bind:value={formData.name}
									placeholder="e.g., My Monstera"
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								/>
							</div>

							<!-- Species -->
							<div>
								<label for="species" class="mb-1 block text-sm font-semibold text-gray-700"
									>Species *</label
								>
								<input
									type="text"
									id="species"
									bind:value={formData.species}
									placeholder="e.g., Monstera deliciosa"
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								/>
							</div>

							<!-- Sun Light -->
							<div>
								<label for="sunlight" class="mb-1 block text-sm font-semibold text-gray-700"
									>Sunlight</label
								>
								<select
									id="sunlight"
									bind:value={formData.sunLight}
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								>
									{#each Object.values(SunlightRequirement) as req (req)}
										<option value={req}>{req}</option>
									{/each}
								</select>
							</div>

							<!-- Temperature -->
							<div>
								<label for="temperature" class="mb-1 block text-sm font-semibold text-gray-700">
									Temperature (°C): {formData.preferedTemperature}
								</label>
								<input
									type="range"
									id="temperature"
									min="-50"
									max="100"
									bind:value={formData.preferedTemperature}
									class="w-full accent-emerald-600"
								/>
							</div>

							<!-- Humidity -->
							<div>
								<label for="humidity" class="mb-1 block text-sm font-semibold text-gray-700">
									Humidity (%): {formData.preferedHumidity}
								</label>
								<input
									type="range"
									id="humidity"
									min="0"
									max="100"
									bind:value={formData.preferedHumidity}
									class="w-full accent-emerald-600"
								/>
							</div>

							<!-- Watering -->
							<div>
								<label
									for="watering-interval"
									class="mb-1 block text-sm font-semibold text-gray-700"
								>
									Watering Interval (days)
								</label>
								<input
									type="number"
									id="watering-interval"
									min="1"
									bind:value={formData.wateringIntervalDays}
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								/>
							</div>

							<!-- Fertilizing -->
							<div>
								<label
									for="fertilizing-interval"
									class="mb-1 block text-sm font-semibold text-gray-700"
								>
									Fertilizing Interval (days)
								</label>
								<input
									type="number"
									id="fertilizing-interval"
									min="1"
									bind:value={formData.fertilizingIntervalDays}
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								/>
							</div>

							<!-- Spray -->
							<div>
								<label for="spray-interval" class="mb-1 block text-sm font-semibold text-gray-700">
									Spray Interval (days, optional)
								</label>
								<input
									type="number"
									id="spray-interval"
									min="1"
									bind:value={formData.sprayIntervalDays}
									placeholder="Leave empty if not needed"
									class="w-full rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
								/>
							</div>

							<!-- Flags -->
							<div>
								<span class="mb-2 block text-sm font-semibold text-gray-700">Flags</span>
								<div class="space-y-2">
									{#each Object.values(PlantFlag) as flag (flag)}
										<label class="flex items-center">
											<input
												type="checkbox"
												checked={formData.flags.includes(flag)}
												on:change={() => toggleFlag(flag)}
												class="h-4 w-4 accent-emerald-600"
											/>
											<span class="ml-2 text-sm text-emerald-900">{flag}</span>
										</label>
									{/each}
								</div>
							</div>

							<!-- Buttons -->
							<div class="flex gap-2">
								<button
									on:click={submitForm}
									disabled={submitting}
									class="flex-1 rounded-lg bg-gradient-to-r from-emerald-600 to-green-600 px-4 py-2 font-semibold text-white shadow-sm transition hover:from-emerald-700 hover:to-green-700 disabled:opacity-50"
								>
									{submitting ? 'Saving...' : 'Save Plant'}
								</button>
								<button
									on:click={() => {
										showForm = false;
										resetForm();
									}}
									class="flex-1 rounded-lg bg-gray-200 px-4 py-2 font-semibold text-gray-800 transition hover:bg-gray-300"
								>
									Cancel
								</button>
							</div>
						</div>
					</div>
				{/if}
			</div>

			<!-- Form Extended Content -->
			{#if showForm}
				<div class="space-y-6 lg:col-span-2">
					<!-- Photos Section -->
					<div
						class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur"
					>
						<h3 class="mb-4 text-xl font-bold text-green-800">📸 Photos</h3>

						<div class="mb-4">
							<span class="mb-2 block text-sm font-semibold text-gray-700">
								Upload Photos (auto-compressed to ≤2MB, first will show on overview)
							</span>
							<input
								type="file"
								bind:this={fileInput}
								on:change={handlePhotoUpload}
								multiple
								accept="image/*"
								class="w-full rounded-lg border-2 border-dashed border-emerald-300 bg-white/70 px-3 py-2"
							/>
						</div>

						{#if photoPreview.length > 0}
							<div class="grid grid-cols-2 gap-4">
								{#each photoPreview as photo, i (i)}
									<div class="relative">
										{#if photo.key}
											<!-- S3 key: show placeholder or presigned URL preview -->
											<div
												class="flex h-32 w-full items-center justify-center rounded-lg bg-green-100 text-center text-sm text-green-700"
											>
												📸 {photo.key.substring(0, 20)}...
											</div>
										{:else if photo.data}
											<!-- Local data URL -->
											<img
												src={photo.data}
												alt="Preview {i + 1}"
												class="h-32 w-full rounded-lg object-cover"
											/>
										{/if}
										<button
											on:click={() => removePhoto(i)}
											class="absolute top-1 right-1 flex h-6 w-6 items-center justify-center rounded-full bg-red-500 text-white shadow hover:bg-red-600"
										>
											×
										</button>
										{#if i === 0}
											<span
												class="absolute top-1 left-1 rounded bg-green-500 px-2 py-1 text-xs text-white"
											>
												Primary
											</span>
										{/if}
										{#if i in photoUploadProgress}
											<div
												class="absolute inset-0 flex items-center justify-center rounded-lg bg-black/50"
											>
												<span class="text-sm text-white">{photoUploadProgress[i]}%</span>
											</div>
										{/if}
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-sm text-gray-500 italic">No photos yet</p>
						{/if}
					</div>

					<!-- Notes Section -->
					<div
						class="rounded-2xl border border-emerald-100 bg-white/90 p-6 shadow-md backdrop-blur"
					>
						<h3 class="mb-4 text-xl font-bold text-green-800">📝 Notes</h3>

						<div class="mb-4 flex gap-2">
							<input
								type="text"
								bind:value={newNote}
								on:keydown={(e) => e.key === 'Enter' && addNote()}
								placeholder="Add a note..."
								class="flex-1 rounded-lg border-2 border-emerald-200 px-3 py-2 shadow-sm focus:border-emerald-500 focus:outline-none"
							/>
							<button
								on:click={addNote}
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
											on:click={() => removeNote(i)}
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
					</div>
				</div>
			{/if}
		</div>

		<!-- Plants List -->
		{#if !loading && plants.length > 0}
			<div class="mt-12">
				<h2 class="mb-6 text-2xl font-bold text-green-800">Your Plants</h2>
				<div class="space-y-3">
					{#each plants as plant (plant.id)}
						<div
							class="flex items-center justify-between rounded-xl border border-emerald-100 bg-white/90 p-4 shadow backdrop-blur transition hover:shadow-lg"
						>
							<div class="flex flex-1 items-center gap-4">
								{#if plant.photoIds.length > 0}
									<div
										class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-200 text-lg font-bold"
									>
										📸
									</div>
								{:else}
									<div
										class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-200 text-xl"
									>
										🌱
									</div>
								{/if}
								<div>
									<h3 class="font-bold text-emerald-900">{plant.name}</h3>
									<p class="text-sm text-emerald-700">{plant.species}</p>
								</div>
							</div>
							<div class="flex gap-2">
								<button
									on:click={() => startEdit(plant)}
									class="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-emerald-700"
								>
									Edit
								</button>
								<button
									on:click={() => deletePlant(plant.id)}
									class="rounded-lg bg-red-500 px-4 py-2 text-sm font-semibold text-white shadow-sm transition hover:bg-red-600"
								>
									Delete
								</button>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		{#if loading}
			<div class="py-12 text-center">
				<div class="mb-4 animate-bounce text-6xl">🌿</div>
				<p class="text-lg font-medium text-green-700">Loading...</p>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
</style>
