import { writable } from 'svelte/store';

type Language = 'en' | 'de' | 'es';

interface Translations {
	[key: string]: string | Translations;
}

const translations: Record<Language, Translations> = {
	en: {},
	de: {},
	es: {}
};

let translationsLoaded = false;

// Load translations dynamically
async function loadTranslations() {
	if (translationsLoaded) return;

	try {
		translations.en = (await import('./i18n/translations/en.json')).default as Translations;
		translations.de = (await import('./i18n/translations/de.json')).default as Translations;
		translations.es = (await import('./i18n/translations/es.json')).default as Translations;
		translationsLoaded = true;
	} catch (error) {
		console.error('Failed to load translations:', error);
	}
}

// Initialize translations
loadTranslations();

// Detect language from URL or browser
function detectLanguage(): Language {
	// Check URL path
	if (typeof window !== 'undefined') {
		const path = window.location.pathname;
		if (path.startsWith('/de')) return 'de';
		if (path.startsWith('/es')) return 'es';
		if (path.startsWith('/en')) return 'en';

		// Fallback to browser language
		const browserLang = navigator.language.split('-')[0];
		if (browserLang === 'de') return 'de';
		if (browserLang === 'es') return 'es';
	}

	return 'en';
}

function createLanguageStore() {
	const { subscribe, set } = writable<Language>(detectLanguage());

	return {
		subscribe,
		set,
		setLanguage: (lang: Language) => set(lang)
	};
}

export const languageStore = createLanguageStore();

function getNestedTranslation(obj: Translations, keys: string[]): string | undefined {
	let current: string | Translations = obj;
	for (const key of keys) {
		if (typeof current === 'object' && current && key in current) {
			current = current[key];
		} else {
			return undefined;
		}
	}
	return typeof current === 'string' ? current : undefined;
}

export function tStore(key: string, args?: string[]): string {
	let result = key;
	let currentLang: Language = 'en';

	// Get current language synchronously
	const unsubscribe = languageStore.subscribe((lang) => {
		currentLang = lang;
	});
	unsubscribe();

	const keys = key.split('.');
	const translation = getNestedTranslation(translations[currentLang], keys);

	result = translation || key;

	// Replace %s placeholders with args
	if (args) {
		for (const arg of args) {
			result = result.replace('%s', arg);
		}
	}

	return result;
}

export type { Language };
