// Service worker: image caching only.
//
// API responses are deliberately NOT cached. The previous cache-first strategy
// with a 15-minute TTL caused hard-to-debug inconsistencies (stale plants after
// watering, changes from another device not showing up, invalidation races), so
// every /api/ request now goes straight to the network.
const CACHE_NAME = 'plants-images-v1';

self.addEventListener('install', () => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      // Remove the API cache left behind by earlier service worker versions.
      await caches.delete('plants-api-v1');
      await self.clients.claim();
    })()
  );
});

self.addEventListener('fetch', (event) => {
  const req = event.request;
  const url = new URL(req.url);

  // API requests: always hit the network, never the cache. (An old service
  // worker version may still hold cached API responses, so drop them.)
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(fetch(req));
    return;
  }

  // Cache images with a cache-first strategy.
  const isImage =
    req.destination === 'image' || /\.(png|jpg|jpeg|webp|gif|svg)(\?.*)?$/.test(url.pathname);
  if (isImage && req.method === 'GET') {
    event.respondWith(
      caches.open(CACHE_NAME).then(async (cache) => {
        const cached = await cache.match(req);
        const fetchPromise = fetch(req)
          .then((res) => {
            if (res && res.ok) {
              cache.put(req, res.clone()).catch(() => {});
            }
            return res;
          })
          .catch(() => cached);
        return cached || fetchPromise;
      })
    );
  }
});
