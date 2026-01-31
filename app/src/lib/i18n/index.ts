import { get, writable } from 'svelte/store';
import { currentLanguage, languageStore } from '../stores/language.js';
import type { Record as TranslationRecord } from './types.js';

const locales = ['en', 'de', 'es'] as const;
type Locale = (typeof locales)[number];
const concerns = ['common', 'menu', 'auth', 'plants', 'time'] as const;

const translations: Partial<Record<Locale, TranslationRecord>> = {};
const loadedLocales = new Set<Locale>();
const loadingLocales = new Map<Locale, Promise<void>>();

let activeLocale: Locale = 'en';

function normalizeLocale(locale: string): Locale {
	return locales.includes(locale as Locale) ? (locale as Locale) : 'en';
}

function translate(key: string, locale: Locale): string {
	const parts = key.split('.');

	// Allow single-part keys mapped to 'common'
	let concern: string;
	let rest: string[];

	if (parts.length === 1) {
		concern = 'common';
		rest = parts.slice();
	} else if (parts.length >= 2) {
		[concern, ...rest] = parts;
	} else {
		return key;
	}

	const concernTranslation = translations[locale]?.[concern];
	if (!concernTranslation) {
		if (locale !== 'en') {
			return translate(key, 'en');
		}

		// Concern not loaded yet
		return key;
	}

	let concernTranslationValue: any = concernTranslation;

	for (const part of rest) {
		if (typeof concernTranslationValue === 'object' && concernTranslationValue !== null) {
			concernTranslationValue = concernTranslationValue[part];
			if (concernTranslationValue === undefined) {
				break;
			}
		} else {
			break;
		}
	}

	// If not found and not already using English fallback, try English
	if (concernTranslationValue === undefined && locale !== 'en') {
		return translate(key, 'en');
	}

	return typeof concernTranslationValue === 'string' ? concernTranslationValue : key;
}

// Create writable store for reactive translations
const tStoreInternal = writable<(key: string) => string>((key: string) =>
	translate(key, activeLocale)
);
export const tStore = { subscribe: tStoreInternal.subscribe };

async function loadLocale(locale: Locale) {
	if (loadedLocales.has(locale)) return;

	const inflight = loadingLocales.get(locale);
	if (inflight) return inflight;

	const loadPromise = (async () => {
		const loaded: TranslationRecord = {};
		for (const concern of concerns) {
			const module = await import(`./translations/${locale}/${concern}.json`);
			loaded[concern] = module.default;
		}
		translations[locale] = loaded;
		loadedLocales.add(locale);
		// Trigger store update after loading
		tStoreInternal.set((key: string) => translate(key, activeLocale));
	})()
		.catch((error) => {
			console.error(`Failed to load locale: ${locale}`, error);
		})
		.finally(() => {
			loadingLocales.delete(locale);
		});

	loadingLocales.set(locale, loadPromise);
	return loadPromise;
}

languageStore.subscribe((nextLocale) => {
	const normalized = normalizeLocale(nextLocale);
	activeLocale = normalized;

	// Load locales and update store
	Promise.all([
		loadLocale(normalized),
		normalized !== 'en' ? loadLocale('en') : Promise.resolve()
	]).then(() => {
		// Update store after both locales are loaded
		tStoreInternal.set((key: string) => translate(key, activeLocale));
	});
});

export function t(key: string, overrideLocale?: string): string {
	const locale = normalizeLocale(overrideLocale ?? activeLocale);
	return translate(key, locale);
}

export function tSync(key: string, overrideLocale?: string): string {
	const locale = normalizeLocale(overrideLocale ?? activeLocale);
	return translate(key, locale);
}

export async function initializeI18n() {
	const locale = normalizeLocale(get(currentLanguage));
	await loadLocale(locale);
	if (locale !== 'en') {
		await loadLocale('en');
	}
	// Update store after initialization
	tStoreInternal.set((key: string) => translate(key, activeLocale));
}

// Eager-load English to avoid raw keys on first render
void loadLocale('en');
