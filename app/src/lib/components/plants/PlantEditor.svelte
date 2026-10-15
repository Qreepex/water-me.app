<script lang="ts">
	import { beforeNavigate, goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount, untrack } from 'svelte';
	import { fetchData } from '$lib/auth/fetch.svelte';
	import FormActionBar from '$lib/components/layout/FormActionBar.svelte';
	import PageContent from '$lib/components/layout/PageContent.svelte';
	import PageHeader from '$lib/components/layout/PageHeader.svelte';
	import Scrollable from '$lib/components/layout/Scrollable.svelte';
	import BasicInformationForm from '$lib/components/plants/forms/BasicInformationForm.svelte';
	import FertilizingForm from '$lib/components/plants/forms/FertilizingForm.svelte';
	import LocationForm from '$lib/components/plants/forms/LocationForm.svelte';
	import MetadataForm from '$lib/components/plants/forms/MetadataForm.svelte';
	import MistingForm from '$lib/components/plants/forms/MistingForm.svelte';
	import SeasonalityForm from '$lib/components/plants/forms/SeasonalityForm.svelte';
	import SoilForm from '$lib/components/plants/forms/SoilForm.svelte';
	import WateringForm from '$lib/components/plants/forms/WateringForm.svelte';
	import Alert from '$lib/components/ui/Alert.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Icon from '$lib/components/ui/Icon.svelte';
	import { toneBg, toneText } from '$lib/components/ui/tone';
	import { tStore } from '$lib/i18n';
	import { getPlantsStore } from '$lib/stores/plants.svelte';
	import type { Plant } from '$lib/types/api';
	import type { FormData } from '$lib/types/forms';
	import { createEmptyFormData } from '$lib/types/forms';
	import {
		EDITOR_SECTIONS,
		buildCreatePayload,
		buildFullPayload,
		initializeFormData,
		isBasicInfoValid,
		isFormDirty
	} from '$lib/utils/plantForm';
	import EditorSectionCard from './EditorSectionCard.svelte';
	import PlantImage from './PlantImage.svelte';

	interface Props {
		mode: 'create' | 'edit';
		/** Required in edit mode */
		plant?: Plant;
	}

	const { mode, plant }: Props = $props();
	// Editor instances are keyed per plant, so capturing initial props is intentional.
	const isCreate = untrack(() => mode === 'create');

	const plantsStore = getPlantsStore();

	let formData = $state<FormData>(
		untrack(() => (plant ? initializeFormData(plant) : createEmptyFormData()))
	);
	let snapshot = $state<FormData>(untrack(() => JSON.parse(JSON.stringify(formData))));
	let saving = $state(false);
	let error = $state<string | null>(null);
	let activeSection = $state('basic');
	let toast = $state<string | null>(null);
	let toastTimer: ReturnType<typeof setTimeout> | undefined;
	let allowLeave = false;

	const isDirty = $derived(isFormDirty(formData, snapshot));
	const canSave = $derived(isDirty && isBasicInfoValid(formData));

	// ---- Live per-section summaries (update as the user types) ----
	const summaries = $derived.by((): Record<string, string> => {
		const t = $tStore;
		return {
			basic: [formData.name, formData.species].filter(Boolean).join(' · '),
			location:
				[formData.room, formData.position].filter(Boolean).join(' · ') ||
				(formData.isOutdoors ? t('plants.outdoors') : ''),
			watering: `${t('plants.every')} ${formData.wateringIntervalDays} ${t('plants.days')}`,
			fertilizing: `${t('plants.every')} ${formData.fertilizingIntervalDays} ${t('plants.days')}`,
			humidity:
				`${formData.targetHumidity}%` +
				(formData.requiresMisting
					? ` · ${t('plants.sprayEvery')} ${formData.mistingIntervalDays} ${t('plants.days')}`
					: ''),
			soil: formData.soilType,
			seasonality: `${formData.preferedTemperature}°C`,
			metadata:
				[
					formData.flags.length
						? `${formData.flags.length}× ${t('plants.flags').toLowerCase()}`
						: '',
					formData.notes.length
						? `${formData.notes.length}× ${t('plants.notesTitle').toLowerCase()}`
						: '',
					formData.isToxic ? t('plants.toxic') : ''
				]
					.filter(Boolean)
					.join(' · ') || '—'
		};
	});

	const photoIds = $derived(plant?.photoIds ?? []);
	const photoUrls = $derived((plant as { photoUrls?: string[] } | undefined)?.photoUrls ?? []);

	// ---- Navigation guard for unsaved changes ----
	beforeNavigate((nav) => {
		if (allowLeave || !isDirty) return;
		if (!confirm($tStore('plants.unsavedChangesConfirm'))) {
			nav.cancel();
		}
	});

	onMount(() => {
		// Warn on hard reload/close with unsaved changes
		const handler = (e: BeforeUnloadEvent) => {
			if (isDirty) e.preventDefault();
		};
		window.addEventListener('beforeunload', handler);

		// Deep link: /manage/:id/edit#watering jumps straight to that section
		const hash = window.location.hash.slice(1);
		if (hash) {
			setTimeout(() => scrollToSection(hash, false), 80);
		}

		return () => window.removeEventListener('beforeunload', handler);
	});

	// ---- Scroll spy ----
	function sectionsRoot(node: HTMLDivElement) {
		const sections = Array.from(node.querySelectorAll<HTMLElement>('[data-editor-section]'));
		const observer = new IntersectionObserver(
			(entries) => {
				for (const entry of entries) {
					if (entry.isIntersecting) {
						activeSection = (entry.target as HTMLElement).dataset.editorSection ?? activeSection;
					}
				}
			},
			{ rootMargin: '-20% 0px -65% 0px' }
		);
		sections.forEach((s) => observer.observe(s));
		return { destroy: () => observer.disconnect() };
	}

	function scrollToSection(key: string, smooth = true): void {
		activeSection = key;
		document
			.getElementById(`editor-section-${key}`)
			?.scrollIntoView({ behavior: smooth ? 'smooth' : 'auto', block: 'start' });
	}

	function showToast(text: string): void {
		toast = text;
		clearTimeout(toastTimer);
		toastTimer = setTimeout(() => (toast = null), 2200);
	}

	// ---- Save ----
	async function save(): Promise<void> {
		if (!isBasicInfoValid(formData)) {
			error = $tStore('plants.requiredNameSpecies');
			scrollToSection('basic');
			return;
		}

		saving = true;
		error = null;

		try {
			if (isCreate) {
				const res = await fetchData('/api/plants', {
					method: 'post',
					body: buildCreatePayload(formData) as never
				});

				if (!res.ok) throw new Error(res.error?.message || $tStore('plants.failedToCreatePlant'));
				if (!res.data?.id) throw new Error($tStore('plants.invalidResponse'));

				plantsStore.setPlants([...plantsStore.plants, res.data]);

				allowLeave = true;
				await goto(resolve(`/manage/${res.data.id}/photos`));
			} else if (plant) {
				const res = await fetchData('/api/plants/{id}', {
					method: 'patch',
					params: { id: plant.id },
					body: buildFullPayload(formData)
				});

				if (!res.ok) throw new Error(res.error?.message || $tStore('plants.failedToUpdatePlant'));

				snapshot = JSON.parse(JSON.stringify(formData));

				if (res.data) {
					plantsStore.setPlants(
						plantsStore.plants.map((p) => (p.id === plant.id ? (res.data as Plant) : p))
					);
				}
				showToast($tStore('common.saved'));
			}
		} catch (err) {
			error = err instanceof Error ? err.message : $tStore('plants.failedToUpdatePlant');
		} finally {
			saving = false;
		}
	}

	function cancel(): void {
		if (!isCreate && plant) {
			goto(resolve(`/plant/${plant.id}`));
		} else {
			goto(resolve('/'));
		}
	}
</script>

<PageHeader
	icon={isCreate ? 'sprout' : 'pencil'}
	title={isCreate ? 'plants.newPlant' : (plant?.name ?? 'plants.editPlant')}
	description="plants.editorDescription"
	backHref={!isCreate && plant ? resolve(`/plant/${plant.id}`) : resolve('/')}
/>

<PageContent>
	<!-- Section navigation (scroll-spy pills) -->
	<nav
		class="flex shrink-0 [scrollbar-width:none] gap-2 overflow-x-auto px-2 pb-3 [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden"
		aria-label={$tStore('plants.manageSections')}
	>
		{#each EDITOR_SECTIONS as section (section.key)}
			<button
				type="button"
				onclick={() => scrollToSection(section.key)}
				aria-current={activeSection === section.key ? 'true' : undefined}
				class="flex shrink-0 cursor-pointer items-center gap-1.5 rounded-full px-3.5 py-2 text-sm font-medium whitespace-nowrap shadow-sm transition-all active:scale-95 {activeSection ===
				section.key
					? 'bg-brand text-onbrand shadow-[0_4px_14px_rgba(0,238,87,0.35)]'
					: 'bg-surface text-ink-soft hover:bg-brand/10 hover:text-ink'}"
			>
				<Icon name={section.icon} size={16} />
				{$tStore(section.titleKey)}
			</button>
		{/each}
	</nav>

	{#if error}
		<div class="mx-2 mb-3 shrink-0">
			<Alert type="error" title="common.error" description={error} />
		</div>
	{/if}

	<Scrollable noPadding>
		<!-- pb reserves space for the fixed FormActionBar -->
		<div class="space-y-4 px-2 pb-48" use:sectionsRoot>
			<EditorSectionCard
				id="basic"
				icon="leaf"
				titleKey="plants.basicInformation"
				summary={summaries.basic}
			>
				<BasicInformationForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="photos"
				icon="camera"
				titleKey="plants.photos"
				tone="info"
				summary={!isCreate ? `${photoIds.length}` : undefined}
			>
				{#if isCreate}
					<p class="text-sm text-ink-soft">{$tStore('plants.photosAfterCreateHint')}</p>
				{:else if plant}
					<div class="flex flex-wrap items-center gap-3">
						{#each photoIds.slice(0, 4) as photoId, i (photoId)}
							<PlantImage
								{photoId}
								remoteUrl={photoUrls[i]}
								alt={plant.name}
								class="h-16 w-16 rounded-xl object-cover"
								fallbackClass="h-16 w-16 rounded-xl text-2xl"
							/>
						{/each}
						{#if photoIds.length > 4}
							<span class="text-sm font-medium text-ink-soft">+{photoIds.length - 4}</span>
						{/if}
					</div>
					<div class="mt-3">
						<Button
							variant="ghost"
							size="sm"
							icon="camera"
							text="plants.managePhotos"
							onclick={() => goto(resolve(`/manage/${plant.id}/photos`))}
						/>
					</div>
				{/if}
			</EditorSectionCard>

			<EditorSectionCard
				id="location"
				icon="map-pin"
				titleKey="plants.location"
				summary={summaries.location}
			>
				<LocationForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="watering"
				icon="droplet"
				titleKey="plants.wateringTitle"
				tone="info"
				summary={summaries.watering}
			>
				<WateringForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="fertilizing"
				icon="flask-conical"
				titleKey="plants.fertilizingTitle"
				tone="warn"
				summary={summaries.fertilizing}
			>
				<FertilizingForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="humidity"
				icon="cloud-drizzle"
				titleKey="plants.humidityTitle"
				tone="info"
				summary={summaries.humidity}
			>
				<MistingForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="soil"
				icon="layers"
				titleKey="plants.soilTitle"
				summary={summaries.soil}
			>
				<SoilForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="seasonality"
				icon="snowflake"
				titleKey="plants.seasonalityTitle"
				tone="warn"
				summary={summaries.seasonality}
			>
				<SeasonalityForm bind:formData />
			</EditorSectionCard>

			<EditorSectionCard
				id="metadata"
				icon="tag"
				titleKey="plants.metadata"
				tone="neutral"
				summary={summaries.metadata}
			>
				<MetadataForm bind:formData />
			</EditorSectionCard>
		</div>
	</Scrollable>

	<FormActionBar
		onCancel={cancel}
		onSave={save}
		{saving}
		{canSave}
		cancelText={isCreate ? 'common.cancel' : 'common.close'}
		saveText={isCreate ? 'plants.createPlant' : 'common.save'}
	/>

	<!-- Save confirmation toast -->
	{#if toast}
		<div
			class="fixed bottom-44 left-1/2 z-50 -translate-x-1/2 rounded-full px-5 py-2.5 text-sm font-semibold shadow-lg {toneBg[
				'ok'
			]} {toneText['ok']}"
			role="status"
		>
			<span class="inline-flex items-center gap-2"
				><Icon name="check-circle" size={16} />{toast}</span
			>
		</div>
	{/if}
</PageContent>
