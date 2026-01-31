import { t } from '$lib/i18n';
import { SvelteDate } from 'svelte/reactivity';

export function formatPastTimestamp(timestamp: Date): string {
	const now = new SvelteDate();

	const clockTime = timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

	if (isSameDay(timestamp, now)) {
		return t('time.today:time').replace('{time}', clockTime);
	} else if (isYesterday(timestamp, now)) {
		return t('time.yesterday:time').replace('{time}', clockTime);
	} else {
		const weekday = timestamp.toLocaleDateString('en', { weekday: 'long' }).toLowerCase();
		return t(`time.${weekday}:time`).replace('{time}', clockTime);
	}
}

function isSameDay(date1: Date, date2: Date): boolean {
	return (
		date1.getFullYear() === date2.getFullYear() &&
		date1.getMonth() === date2.getMonth() &&
		date1.getDate() === date2.getDate()
	);
}

function isYesterday(date: Date, reference: Date): boolean {
	const yesterday = reference;
	yesterday.setDate(reference.getDate() - 1);
	return isSameDay(date, yesterday);
}

function isTomorrow(date: Date, reference: Date): boolean {
	const tomorrow = reference;
	tomorrow.setDate(reference.getDate() + 1);
	return isSameDay(date, tomorrow);
}

export function formatFutureTimestamp(timestamp: Date): string {
	const now = new SvelteDate();

	const clockTime = timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

	if (isSameDay(timestamp, now)) {
		return t('time.today:time').replace('{time}', clockTime);
	} else if (isTomorrow(timestamp, now)) {
		return t('time.tomorrow:time').replace('{time}', clockTime);
	} else {
		const weekday = timestamp.toLocaleDateString('en', { weekday: 'long' }).toLowerCase();
		return t(`time.${weekday}:time`).replace('{time}', clockTime);
	}
}
