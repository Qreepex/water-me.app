type InvalidateOptions = {
	waitForAck?: boolean;
	timeoutMs?: number;
};

export async function invalidateApiCache(
	urls: string[],
	options: InvalidateOptions = {}
): Promise<void> {
	if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) return;
	const controller = navigator.serviceWorker.controller;
	if (!controller) return;

	const { waitForAck = false, timeoutMs = 200 } = options;

	if (!waitForAck) {
		controller.postMessage({
			type: 'INVALIDATE_CACHE',
			urls
		});
		return;
	}

	await new Promise<void>((resolve) => {
		const channel = new MessageChannel();
		channel.port1.onmessage = () => resolve();
		controller.postMessage(
			{
				type: 'INVALIDATE_CACHE',
				urls
			},
			[channel.port2]
		);
		setTimeout(() => resolve(), timeoutMs);
	});
}
