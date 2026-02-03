import { Capacitor } from '@capacitor/core';
import { DefaultSystemBrowserOptions, InAppBrowser } from '@capacitor/inappbrowser';

export async function openExternalLink(url: string): Promise<void> {
	if (Capacitor.isNativePlatform()) {
		await InAppBrowser.openInSystemBrowser({ url, options: DefaultSystemBrowserOptions });
	} else {
		window.open(url, '_blank', 'noopener');
	}
}
