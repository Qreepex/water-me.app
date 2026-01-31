const CACHE_NAME = 'plants-images-v1';
const API_CACHE_NAME = 'plants-api-v1';
const CACHE_TTL = 15 * 60 * 1000; // 15 minutes in milliseconds

self.addEventListener('install', () => {
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim());
});

// Listen for cache invalidation messages from the client
self.addEventListener('message', (event) => {
  if (event.data.type === 'INVALIDATE_CACHE') {
    const urlsToInvalidate = event.data.urls || [];
    console.log('📢 Cache invalidation requested for:', urlsToInvalidate);
    (async () => {
      const cache = await caches.open(API_CACHE_NAME);
      
      // Get all cached requests and delete matching ones
      const keys = await cache.keys();
      console.log('📦 Cached URLs before invalidation:', keys.map(k => k.url));
      
      let deletedCount = 0;
      for (const key of keys) {
        const keyUrl = new URL(key.url);
        const keyPath = keyUrl.pathname;
        console.log('🔍 Checking:', keyPath);
        
        // Check if this cached request matches any of the URLs to invalidate
        for (const invalidateUrl of urlsToInvalidate) {
          // Exact match or exact match with trailing slash
          if (keyPath === invalidateUrl || keyPath === invalidateUrl + '/') {
            await cache.delete(key);
            console.log('✓ Deleted exact match:', keyPath);
            deletedCount++;
          }
          // For /api/plants/ID patterns, check if the invalidateUrl is the base path
          else if (invalidateUrl.endsWith('/ID')) {
            // This shouldn't happen - placeholder paths should not be in invalidation URLs
          }
          // Check if it's a detail URL that should be deleted when invalidating list URL
          else if (invalidateUrl === '/api/plants' && keyPath.match(/^\/api\/plants\/[a-f0-9]{24}$/)) {
            // When invalidating /api/plants, also invalidate all /api/plants/ID endpoints
            await cache.delete(key);
            console.log('✓ Deleted detail URL when invalidating list:', keyPath);
            deletedCount++;
          }
          // Check if it's a list URL that should be deleted when invalidating a detail URL
          else if (invalidateUrl.match(/^\/api\/plants\/[a-f0-9]{24}$/) && keyPath === '/api/plants') {
            // When invalidating /api/plants/ID, also invalidate /api/plants list
            await cache.delete(key);
            console.log('✓ Deleted list URL when invalidating detail:', keyPath);
            deletedCount++;
          }
        }
      }
      
      console.log(`✓ Cache invalidation complete. Deleted ${deletedCount} entries.`);
      const keysAfter = await cache.keys();
      console.log('📦 Cached URLs after invalidation:', keysAfter.map(k => k.url));
      
      // Notify the client that cache invalidation is complete (if MessageChannel port provided)
      if (event.ports && event.ports[0]) {
        event.ports[0].postMessage({ success: true });
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
        } catch (error) {
          console.error('Fetch failed; returning cached data if available.', error);

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

