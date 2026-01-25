const CACHE_NAME = 'plants-images-v1';
const API_CACHE_NAME = 'plants-api-v1';
const CACHE_TTL = 15 * 60 * 1000; // 15 minutes in milliseconds
const MAX_ENTRIES = 200;

self.addEventListener('install', (event) => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});

// Listen for cache invalidation messages from the client
self.addEventListener('message', (event) => {
  if (event.data.type === 'INVALIDATE_CACHE') {
    const urlsToInvalidate = event.data.urls || [];
    (async () => {
      const cache = await caches.open(API_CACHE_NAME);
      
      // Get all cached requests and delete matching ones
      const keys = await cache.keys();
      for (const key of keys) {
        const keyUrl = new URL(key.url);
        // Check if this cached request matches any of the URLs to invalidate
        for (const invalidateUrl of urlsToInvalidate) {
          if (keyUrl.pathname === invalidateUrl || keyUrl.pathname.includes(invalidateUrl)) {
            await cache.delete(key);
            console.log('Cache invalidated for:', key.url);
          }
        }
      }
    })();
  }
});

// Helper to check if cached response is expired
async function isCacheExpired(cache, req) {
  const cached = await cache.match(req, { ignoreSearch: true });
  if (!cached) return true;
  
  const cacheTime = cached.headers.get('X-Cache-Time');
  if (!cacheTime) return false; // No timestamp, assume valid
  
  const age = Date.now() - parseInt(cacheTime, 10);
  return age > CACHE_TTL;
}

// Helper to add timestamp to response headers
function addCacheTimestamp(response) {
  const cloned = response.clone();
  const headers = new Headers(cloned.headers);
  headers.set('X-Cache-Time', Date.now().toString());
  return new Response(cloned.body, {
    status: cloned.status,
    statusText: cloned.statusText,
    headers: headers
  });
}

self.addEventListener('fetch', (event) => {
  const req = event.request;
  const url = new URL(req.url);
  
  // Cache images with cache-first strategy
  const isImage = req.destination === 'image' || /\.(png|jpg|jpeg|webp|gif|svg)(\?.*)?$/.test(url.pathname);
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
    return;
  }

  // Cache API GET requests with cache-first strategy (with 15min TTL)
  const isAPI = url.pathname.startsWith('/api/');
  if (isAPI && req.method === 'GET') {
    event.respondWith(
      (async () => {
        const cache = await caches.open(API_CACHE_NAME);
        
        // Check cache first - if found and not expired, return immediately
        const cached = await cache.match(req, { ignoreSearch: true });
        if (cached && !await isCacheExpired(cache, req)) {
          return cached;
        }
        
        // Cache expired or not found - fetch from network
        try {
          const res = await fetch(req);
          if (res && res.ok) {
            // Add timestamp to response before caching
            const withTimestamp = addCacheTimestamp(res);
            cache.put(req, withTimestamp.clone());
          }
          return res;
        } catch (err) {
          // No network - return cached response even if expired
          if (cached) return cached;
          return new Response(JSON.stringify({ error: { message: 'Network unavailable' } }), {
            status: 503,
            headers: { 'Content-Type': 'application/json' }
          });
        }
      })()
    );
    return;
  }

  // Don't cache POST/PATCH requests - just pass through to network
  if (isAPI && (req.method === 'POST' || req.method === 'PATCH')) {
    return;
  }
});

