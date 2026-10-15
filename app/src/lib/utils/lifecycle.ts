import { App } from '@capacitor/app';
import { Capacitor } from '@capacitor/core';

/**
 * Minimum time between two resume-triggered refreshes. Returning to the app can
 * fire several events in a row (visibilitychange + focus + appStateChange), and
 * quick tab switching should not hammer the API.
 */
const MIN_INTERVAL_MS = 15_000;

/**
 * Call `onResume` whenever the user returns to the app: the tab becomes visible
 * again (web) or the native app is brought back to the foreground (Capacitor).
 *
 * Used to pick up changes made on another device, because the plant store is
 * only loaded once per session.
 *
 * Returns a cleanup function that removes all listeners.
 */
export function onAppResume(onResume: () => void): () => void {
	let lastRun = 0;

	const trigger = (): void => {
		const now = Date.now();
		if (now - lastRun < MIN_INTERVAL_MS) return;
		lastRun = now;
		onResume();
	};

	const cleanups: Array<() => void> = [];

	if (typeof document !== 'undefined') {
		const onVisibilityChange = (): void => {
			if (document.visibilityState === 'visible') trigger();
		};
		document.addEventListener('visibilitychange', onVisibilityChange);
		cleanups.push(() => document.removeEventListener('visibilitychange', onVisibilityChange));

		// Some browsers/WebViews fire focus instead of (or in addition to) the
		// visibility event above.
		window.addEventListener('focus', trigger);
		cleanups.push(() => window.removeEventListener('focus', trigger));
	}

	if (Capacitor.isNativePlatform()) {
		// Android's WebView does not reliably fire visibilitychange, so the
		// native lifecycle event is the dependable signal there.
		void App.addListener('appStateChange', ({ isActive }) => {
			if (isActive) trigger();
		}).then((listener) => {
			cleanups.push(() => void listener.remove());
		});
	}

	return () => {
		for (const cleanup of cleanups) cleanup();
		cleanups.length = 0;
	};
}
