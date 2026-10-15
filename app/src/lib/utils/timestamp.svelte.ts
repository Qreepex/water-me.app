import { SvelteDate } from 'svelte/reactivity';

/**
 * A translatable timestamp: the caller resolves it with `$tStore(key, args)` so
 * the text stays reactive (a plain `t()` call computed once would freeze the raw
 * key if the locale had not finished loading yet).
 */
export interface TimestampText {
	key: string;
	args: string[];
}

function clockTime(timestamp: Date): string {
	return timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function weekdayKey(timestamp: Date): string {
	return timestamp.toLocaleDateString('en', { weekday: 'long' }).toLowerCase();
}

export function pastTimestamp(timestamp: Date): TimestampText {
	const now = new SvelteDate();
	if (isSameDay(timestamp, now)) {
		return { key: 'time.today:time', args: [clockTime(timestamp)] };
	}
	if (isYesterday(timestamp, now)) {
		return { key: 'time.yesterday:time', args: [clockTime(timestamp)] };
	}
	return { key: `time.${weekdayKey(timestamp)}:time`, args: [clockTime(timestamp)] };
}

export function futureTimestamp(timestamp: Date): TimestampText {
	const now = new SvelteDate();
	if (isSameDay(timestamp, now)) {
		return { key: 'time.today:time', args: [clockTime(timestamp)] };
	}
	if (isTomorrow(timestamp, now)) {
		return { key: 'time.tomorrow:time', args: [clockTime(timestamp)] };
	}
	return { key: `time.${weekdayKey(timestamp)}:time`, args: [clockTime(timestamp)] };
}

function isSameDay(date1: Date, date2: Date): boolean {
	return (
		date1.getFullYear() === date2.getFullYear() &&
		date1.getMonth() === date2.getMonth() &&
		date1.getDate() === date2.getDate()
	);
}

function isYesterday(date: Date, reference: Date): boolean {
	const yesterday = new SvelteDate(reference);
	yesterday.setDate(reference.getDate() - 1);
	return isSameDay(date, yesterday);
}

function isTomorrow(date: Date, reference: Date): boolean {
	const tomorrow = new SvelteDate(reference);
	tomorrow.setDate(reference.getDate() + 1);
	return isSameDay(date, tomorrow);
}
