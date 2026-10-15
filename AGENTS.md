# AGENTS.md — Water Me! (water-me.app)

Guidance for all coding agents working in this repository. This is the single,
central instruction file for the project — read it first, then the per-project
READMEs. Everything here was verified against the code; when in doubt, trust the
code over older docs (`README.md`, `DESIGN.md`, `TODO.md`).

> **Authoritative docs, in order of freshness:** source code → `TODO.md` (live
> backlog) → this file → `backend/README.md` / `backend/openapi.yaml` → `DESIGN.md`
> (design tokens).

---

## 1. What this project is (function level)

**Water Me!** is a plant-care / watering-reminder app:

- Users sign in with Google (Firebase Auth) and manage a personal collection of
  **plants** (≤ 50 per user).
- Each plant stores care config: watering, fertilizing, misting/humidity, soil &
  repotting, seasonality, sunlight, location, toxicity, flags, notes, photos, and
  care/growth history.
- A **"Water" view** ranks plants by urgency (overdue → due soon → ok) and lets the
  user water one or many plants with a hold-to-confirm gesture.
- A **background worker** sends FCM push notifications when plants need water /
  fertilizer / misting / repotting.
- A **marketing website** (landing page in EN/DE/ES, legal pages, live stats)
  advertises the app and links to the web app / stores.

The three deliverables live side by side in this monorepo. There is **no root
`package.json`** — each part is built independently.

| Dir         | What                                            | Stack                                                          | Output / Deploy                       |
| ----------- | ----------------------------------------------- | -------------------------------------------------------------- | ------------------------------------- |
| `backend/`  | REST API + notification worker + cleanup worker | Go 1.24, gorilla/mux, MongoDB, Firebase Admin, AWS SDK v2 (S3) | Docker image → k3s (`api.water-me.app`) |
| `app/`      | Product app (Android via Capacitor + web PWA)   | SvelteKit 2 / Svelte 5 runes, Tailwind 4, Capacitor 8          | static `www/` → Cloudflare (web) + APK |
| `website/`  | Marketing / landing / legal site                | SvelteKit 2 / Svelte 5, Tailwind 4                             | static `www/` → Cloudflare Pages        |

Domains: app `https://my.water-me.app`, API `https://api.water-me.app`, site
`https://water-me.app`.

---

## 2. Repo map — where to look

```
backend/
  cmd/api/                    # API server entrypoint (port 8080)
  cmd/notification-worker/    # 5-min loop: send care reminders
  cmd/cleanup-worker/         # 30-min loop: delete orphaned S3 uploads
  cmd/setup-cors/             # one-off S3/CORS bootstrap
  routes/                     # HTTP handlers: plants.go, uploads.go, notifications.go, stats.go, routes.go
  services/                   # mongo.go (client), database.go (queries), firebase.go, s3.go, uploads.go, notification_worker.go, ratelimit.go
  middlewares/auth.go         # Firebase ID-token verification → userID in context
  validation/                 # Go validation for plants, limits, notification config
  types/                      # Go structs (plants.go is the domain model)
  constants/constants.go      # collection names, upload limits, allowed MIME types
  messages/*.json             # notification message templates (watering/fertilizing/misting/repotting)
  util/helpers.go             # JSON request/response helpers (JSON decode is strict)
  openapi.yaml                # API spec (INCOMPLETE — see gotchas)
  infra/k8s/*.yaml            # namespace, api (2 replicas + ingress), notification-worker (1 replica)
  Dockerfile, build.sh        # multi-component image build (build-arg COMPONENT=api|notification-worker)
  VERSION                     # image tag used by build.sh / k8s
app/
  src/routes/                 # SvelteKit pages (see §5)
  src/lib/auth/               # Auth.svelte (login gate), fetch.svelte.ts (typed API client), firebase.ts (config)
  src/lib/stores/             # runes singletons: plants, watering, imageCache (+ legacy language.ts)
  src/lib/components/ui/      # design-system primitives (Button, Card, Icon, tone.ts, forms/*)
  src/lib/components/layout/  # PageHeader, PageContent, Scrollable, List, BottomNav, AppMenu, FormActionBar
  src/lib/components/plants/  # PlantList/Card/Editor/Image + forms/* (one per editor section)
  src/lib/utils/              # plant.ts (watering math/sorting), plantForm.ts (form↔payload), imageCache.ts, uploads.ts, cache.ts
  src/lib/i18n/               # custom loader + translations/{de,en,es}/{common,menu,auth,plants,time}.json
  src/lib/types/api.ts        # AUTO-GENERATED from backend/openapi.yaml (do not edit)
  src/routes/layout.css       # Tailwind 4 theme + design tokens (semantic colors)
  static/sw.js                # service worker: image + API cache, invalidation
  capacitor.config.json       # Android app id app.waterme.app
  android/                    # Capacitor Android project
website/
  src/routes/                 # / (lang redirect), /en /de /es, /privacy /terms /imprint
  src/lib/components/landing/ # Hero, Features, Screenshots, Stats, Download, Faq, Footer, Nav, LegalShell
  src/lib/i18n.svelte.ts      # synchronous JSON translations
  src/lib/index.ts            # URLs + API_BASE_URL
DESIGN.md                     # "Emerald Oasis" palette (source of truth for design)
TODO.md                       # live backlog (Fixes / Features / Future)
```

---

## 3. Design system ("Emerald Oasis")

`DESIGN.md` is the canonical palette. It is already encoded as CSS variables +
Tailwind 4 theme in **`app/src/routes/layout.css`** (and mirrored in
`website/src/routes/layout.css`).

- Brand emerald `#00EE57` (`--p-emerald`), dark `#00A343`; **always use dark text
  `#061F12` on brand green** (`text-onbrand`/`--on-emerald`).
- Status: error `#E63946`, warn `#FFB703`, info `#219EBC`, success `#008F41`.
- Light/dark schemes flip the *semantic* vars; dark follows `prefers-color-scheme`.
- Tailwind exposes semantic utilities via `@theme inline`: `bg-canvas`, `bg-surface`,
  `text-ink`, `text-ink-soft`, `bg-brand`, `text-brand-dark`, `text-onbrand`,
  `bg-danger|warn|info|ok`, `bg-brand/15` etc.

**Rules:**
- Never hardcode hex colors in components — use the semantic utilities so dark mode
  keeps working (`canvas`/`surface`/`ink`/`ink-soft`).
- Status colors are centralized in **`app/src/lib/components/ui/tone.ts`**
  (`toneBg`, `toneText`, `toneIcon`, type `Tone`). Reuse it instead of ad-hoc classes.
- Icons come from the central registry **`ui/icons.ts`** (`<Icon name="..." />`,
  stroke-based Lucide-style 24×24). Add new icons only there.
- Signature styles: brand gradient wash at the top of every page
  (`--header-wash` in `+layout.svelte`), floating primary buttons with a green glow
  (`shadow-[0_4px_14px_rgba(0,238,87,0.4)]`), and the `.pop-in` staggered card
  entrance animation.
- Layout rules: pages use `PageHeader` + `PageContent` (+ `Scrollable`/`List`);
  content must reserve bottom padding (`pb-40`/`pb-48`) so the fixed `BottomNav`
  and `FormActionBar` never cover content. Safe-area insets (`env(safe-area-inset-*)`)
  are handled in `+layout.svelte`, `BottomNav`, and `FormActionBar`.

---

## 4. Backend (Go)

### 4.1 Components

| Binary                       | Purpose                                                        | Cadence            | Scaling         |
| ---------------------------- | -------------------------------------------------------------- | ------------------ | --------------- |
| `cmd/api`                    | REST API (auth, plants, uploads, notification config, stats)  | HTTP :8080         | Horizontal (2+) |
| `cmd/notification-worker`    | Scans plants needing care, sends FCM push                     | every 5 min        | Single instance |
| `cmd/cleanup-worker`         | Deletes orphaned S3 uploads (unreferenced, older than 1 h)    | every 30 min       | Single instance |
| `cmd/setup-cors`             | One-off S3 bucket CORS/policy setup                           | manual             | —               |

Shared packages: `types`, `services`, `constants`, `validation`, `util`, `middlewares`.
Docker builds a single image parameterized by `--build-arg COMPONENT=<cmd dir>`.

### 4.2 Request lifecycle (API)

1. `gorilla/handlers.CORS` (allowed origins are the same list passed to the S3
   bucket CORS setup in `cmd/api/main.go`).
2. Routes registered in `routes/routes.go` on a `mux.Router`.
3. **`PlantHandler` installs the rate limiter via `router.Use(...)`** — because
   gorilla/mux applies router-level middlewares to *every* matched route, it effectively
   covers all routes, not only plant routes (`services/ratelimit.go`): 100 req/min per
   `userID|IP` and 1000 req/min per IP. Sequencing caveat: `r.Use` order in `main.go` is
   CORS → rate limiter → auth, so the limiter runs **before** auth injects `userID`; in
   practice the per-user dimension is empty and it behaves as a per-IP limit.
4. `middlewares.AuthMiddleware` verifies `Authorization: Bearer <firebase-id-token>`
   via Firebase Admin, then injects `userID` into context under `constants.UserIdKey`
   (read with `routes.getUserID`). OPTIONS requests and `/api/stats` bypass auth.
5. Handlers read `userID`, validate input (`util.DecodeJSON` is **strict** —
   unknown fields are rejected), call `services/database.go`, respond with
   `util.RespondJSON`. Errors: `util.BadRequest` → `{error, details}`,
   `util.ServerError`, `util.NotFound`.

**Every plant/upload/notification query is scoped by `userId`.** Never write a query
without the user filter — this is the core security invariant (see
`routes/notifications_security_test.go`).

### 4.3 Endpoints (implemented in code)

| Method | Path                                  | Notes                                                |
| ------ | ------------------------------------- | ---------------------------------------------------- |
| GET    | `/api/plants`                         | List; enriches `photoUrls` (presigned, 1 h)          |
| POST   | `/api/plants`                         | Create; enforces `MaxPlantsPerUser` (50), slug       |
| GET    | `/api/plants/{id}`                    | Get one                                              |
| GET    | `/api/plants/slug/{slug}`             | Get by slug                                          |
| PATCH  | `/api/plants/{id}`                    | Partial update; empty sub-objects are `$unset`       |
| DELETE | `/api/plants/{id}`                    | Delete                                               |
| POST   | `/api/plants/water`                   | Bulk set `watering.lastWatered = now`                |
| POST   | `/api/uploads/presign`                | Get presigned S3 PUT URL (validates type/size/quota) |
| POST   | `/api/uploads/register`               | Confirm upload (HeadObject) → `uploads` collection   |
| DELETE | `/api/uploads/{key}`                  | Delete object + record                               |
| GET/PUT/DELETE | `/api/notifications`          | Notification config (GET returns defaults if none)   |
| POST   | `/api/notifications/tokens`           | Register FCM device token                            |
| DELETE | `/api/notifications/tokens/{deviceId}` | Remove device token                                 |
| GET    | `/api/stats`                          | **Public**; users/plants/reminders, cached 5 min in-process |

### 4.4 Data model & MongoDB

Collections (`constants.MongoDBCollections`): **`plants`**, **`uploads`**,
**`notifications`**. IDs: `_id` is a Mongo `ObjectID` (24-char hex) for plants;
notification configs use a `nanoid` string. Field names are camelCase (BSON tags match
JSON). `types/plants.go` is the domain model — it also defines all care enums
(`SunlightRequirement`, `WateringMethod`, `WaterType`, `FertilizerType`, `PestType`,
`PestStatus`, `PlantFlag`, `HealthStatus`) which are mirrored in the app.

- **Plants**: user-scoped, unique-per-user `slug` (generated in `routes/plants.go`
  from name, then location, then a counter). On create, `watering.lastWatered` and
  `fertilizing.lastFertilized` default to now. Empty optional sub-objects are removed
  (`isEmpty*` helpers) so the DB does not store unset fields.
- **Uploads**: `{userId, key, sizeBytes, createdAt}`. S3 keys are
  `users/{userId}/{uuid}_{filename}` and **must** carry the user prefix
  (`services.KeyBelongsToUser`). Limits: 2 MB, JPEG/PNG/WebP, 20 uploads/user.
- **Notifications**: one config doc per user — device tokens, enable flags, per-type
  reminders, muted plant IDs, preferred time, quiet hours, batching. See
  `types/notifications.go`.

### 4.5 Notification worker

`services/notification_worker.go` + `cmd/notification-worker`:

- For each care type (watering/fertilizing/misting/repotting) it runs an aggregation
  in `services/database.go` (`GetPlantsNeeding*`, batch 1000) that computes
  `last* + intervalDays <= now`.
- Groups plants by user, loads the user's notification config, then checks:
  `isEnabled`, the per-type flag, a **4-hour cooldown** (`NotificationCooldown`),
  muted plant IDs, and active device tokens.
- Sends FCM multicast in batches of 500 (`FCMBatchSize`), marks permanently failing
  tokens `isActive: false`, and updates `lastNotificationSentAt`.
- Message strings come from `backend/messages/{type}.json` — `single`/`multiple`
  variants with placeholders `{plantName}`, `{count}`, `{plantNames}`, `{remaining}`.
  A random variant is chosen per send; `buildDefaultMessage` is the fallback.

### 4.6 Running the backend locally

```bash
cd backend
go run ./cmd/api                 # :8080
go run ./cmd/notification-worker # care reminders
go run ./cmd/cleanup-worker      # orphaned-upload cleanup
```

Environment (dotenv auto-loads `.env`; `secret/fb.json` holds the Firebase service
account, both gitignored):

```
DATABASE_URL=mongodb://localhost:27017/plants
MONGODB_USERNAME=…  MONGODB_PASSWORD=…  MONGODB_DATABASE=plants
GOOGLE_APPLICATION_CREDENTIALS=./secret/fb.json
AWS_ACCESS_KEY_ID=…  AWS_SECRET_ACCESS_KEY=…
AWS_REGION=de   S3_BUCKET=plants-app-images
S3_ENDPOINT=https://s3.eu-central-4.ionoscloud.com
PORT=8080
```

S3 is **IONOS** (path-style URLs). Tests: `go test ./routes -v` and
`go test ./services -v` (the S3 test skips without AWS creds; the notification tests
are documentation-style security assertions, not behavioral tests).

### 4.7 Backend conventions

- **Add/change an endpoint:** `types/` → `validation/` → `services/database.go` →
  `routes/*.go` → register in `routes/routes.go` → update `openapi.yaml` → regenerate
  app types (§5.4).
- Keep validation strict and mirrored. The app relies on the shape
  `{error: string, details: ValidationError[] | object}` where
  `ValidationError = {field, message}`.
- Format Go with `gofmt`; keep the shared-package layout (no direct DB access outside
  `services/`).

---

## 5. App (`app/` — SvelteKit + Capacitor)

### 5.1 Runtime model

The app is a **client-only SPA**: `src/routes/+layout.ts` sets `ssr = false` and
`prerender = true`, and `svelte.config.js` uses `adapter-static` with
`fallback: '404.html'`, `strict: false` so dynamic routes aren't prerendered. Every
dynamic route (`/plant/[plant]`, `/manage/**`) sets `export const prerender = false;
export const ssr = false;`.

It ships two ways: as a **web PWA** (Cloudflare, `wrangler.jsonc` serving `www/`) and
as an **Android app** via Capacitor (`capacitor.config.json`, app id `app.waterme.app`,
Firebase Google sign-in through `@capacitor-firebase/authentication`).

### 5.2 Auth & data flow

- `src/lib/auth/Auth.svelte` gates all content. It initializes Firebase, observes
  auth state (web SDK vs `FirebaseAuthentication` on native), shows a Google sign-in
  card when logged out, and triggers `plantsStore.loadPlants()` on login.
- `src/lib/auth/fetch.svelte.ts` exposes `fetchData<P>(path, opts)` — a **typed**
  wrapper that attaches the Firebase ID token as a Bearer header and returns a
  discriminated `{ok: true, data} | {ok: false, error}`. All API calls go through it.
- `src/lib/constants.ts` holds `API_BASE_URL` (**hardcoded to production**
  `https://api.water-me.app`). To develop against a local backend, change it there.

### 5.3 State (Svelte 5 runes singletons)

Apps use `$state`-based singleton stores in `src/lib/stores/*.svelte.ts`:

- **`plants.svelte.ts`** — the source of truth for the plant list
  (`getPlantsStore()`): `loadPlants`, `reloadSilently`, `refresh`, prefetches images.
- **`watering.svelte.ts`** — in-flight watering session (`waterPlant`, `waterPlants`),
  with haptics; invalidates the API cache and silently reloads after success.
- **`imageCache.svelte.ts`** — in-memory object-URL cache keyed by photo id.
- `language.ts` is **legacy** (`writable` store) — see gotchas.

### 5.4 OpenAPI → typed client

`src/lib/types/api.ts` is generated, not hand-written:

```bash
npx openapi-typescript ./backend/openapi.yaml -o ./app/src/lib/types/api.ts \
  --root-types --root-types-no-schema-prefix --generate-path-params --enum --enum-values
```

Run this after any backend API change. Enums (e.g. `SunlightRequirement`) are
consumed by forms and `plantForm.ts`. Note the backend response also includes
`photoUrls`, which is **not** in the spec (the app casts for it) — see gotchas.

### 5.5 Routes / screens

| Route                            | Screen                                                                 |
| -------------------------------- | ---------------------------------------------------------------------- |
| `/`                              | "My Plants" list (search, sort, cards → detail)                        |
| `/water`                         | Water view: overdue/due-soon first, "water all" hold button            |
| `/plant/[plant]`                 | Plant detail: mood face, watering progress, care overview, history     |
| `/manage/new`                    | Create plant (unified `PlantEditor`, then → photos)                    |
| `/manage/[plant]/edit`           | Unified editor; hash deep-links (`#watering`) + scroll-spy sections    |
| `/manage/[plant]/edit/[section]` | 307 redirect → `/edit#<section>` (legacy per-section URLs)             |
| `/manage/[plant]/photos`         | Photo upload/manage (compress → presign → PUT)                         |
| `/manage/[plant]`                | 307 redirect → `/edit`                                                 |
| `/notifications`                 | Push opt-in + (placeholder) notification preferences                   |

Navigation shell: fixed `BottomNav` (Garden, Water, burger → `AppMenu` overlay with
language switch, refresh, create plant, legal links via `openExternalLink`).

### 5.6 The plant editor

`components/plants/PlantEditor.svelte` is the heart of write operations. It holds one
`FormData` (`types/forms.ts`), renders `EDITOR_SECTIONS` (basic, photos, location,
watering, fertilizing, humidity, soil, seasonality, metadata), shows live per-section
summaries, does scroll-spy, and guards unsaved changes (`beforeNavigate` + `beforeunload`).
Payload building lives in **`utils/plantForm.ts`**: `buildCreatePayload` sends only
changed fields; `buildFullPayload` sends every section at once but **deliberately omits
care timestamps** (`lastWatered`, `lastFertilized`, …) so saving settings never rewrites
care history. Keep this invariant when adding editor fields.

### 5.7 Caching (two layers — important)

1. **Service worker** `static/sw.js` (registered in `+layout.svelte`, web only):
   - images → cache-first,
   - `GET /api/*` → cache-first with a **15-minute TTL** (`X-Cache-Time` header),
     network fallback returns stale data if offline.
   - Invalidation: `utils/cache.ts → invalidateApiCache(urls)` posts
     `INVALIDATE_CACHE` to the SW (after mutations like create/update/water). It also
     cross-invalidates `/api/plants` ↔ `/api/plants/{id}`.
2. **Persistent image cache** `utils/imageCache.ts` + `stores/imageCache.svelte.ts`:
   Capacitor Filesystem on native / IndexedDB on web, 7-day TTL, WebP.

Because it is a client SPA, the store is also kept coherent by optimistic in-memory
updates (e.g. `plantEditor` patches the store after save) plus a silent reload.

### 5.8 i18n

Custom loader in **`src/lib/i18n/index.ts`**: locales `en|de|es`, "concerns"
`common|menu|auth|plants|time`, each a JSON file under
`src/lib/i18n/translations/<locale>/`. Keys are `concern.path.to.key`; single-segment
keys resolve to `common`. Interpolation uses **`%s`** placeholders (`t(key, [args])`).
Use `$tStore('key')` in components and `t('key')` in TS. English is the fallback.
Language is persisted via Capacitor Preferences (`stores/language.ts`).

**When adding UI text: add the key to all three locales** (`en`, `de`, `es`), and run
`npm run format`/`npm run lint`.

### 5.9 Photos

Upload pipeline (`utils/uploads.ts` + photos page): `compressPhoto` → transcode to WebP
under 2 MB (progressive quality/downscale) → `POST /api/uploads/presign` → `PUT` blob to
S3 with the exact returned headers → `POST /api/uploads/register` → attach the returned
`key` to `plant.photoIds`.

### 5.10 Running / quality

```bash
cd app
npm install
npm run dev        # vite dev server
npm run build      # → www/
npm run check      # svelte-check (TS/Svelte)
npm run lint       # prettier --check + eslint
npm run format     # prettier --write
npx cap sync android   # after web build, before Android build
```

Style: **tabs, single quotes, no trailing commas, print width 100**
(`app/.prettierrc`), Tailwind class order via `prettier-plugin-tailwindcss`.

---

## 6. Website (`website/`)

Static marketing site (SvelteKit + `adapter-static`, `prerender = true`, output `www/`,
deployed via `wrangler.jsonc` to Cloudflare Pages). Synchronous i18n via
`src/lib/i18n.svelte.ts` with `translations/{en,de,es}.json`; language is derived from
the URL prefix (`/en`, `/de`, `/es`); `/` redirects by browser language.
Landing sections live in `src/lib/components/landing/`. The live stats block
(`StatsSection.svelte`) calls the **public** `GET /api/stats` using
`API_BASE_URL` from `src/lib/index.ts` (`https://api.water-me.app/api`). Legal pages:
`/privacy`, `/terms`, `/imprint`.

```bash
cd website
npm install && npm run dev
npm run build      # → www/
npm run check && npm run lint
```

---

## 7. Build & deploy

```bash
# Backend images (tag from backend/VERSION, pushed to ghcr.io/qreepex/*)
cd backend && ./build.sh            # builds plants-backend-api + plants-notification-worker
docker build --build-arg COMPONENT=api -t ghcr.io/qreepex/plants-backend-api:dev .
```

k3s manifests in `backend/infra/k8s/`: namespace `water-me`; `backend-api` (2 replicas
+ ClusterIP service + nginx ingress/TLS for `api.water-me.app`); `notification-worker`
(1 replica). Both mount a `firebase-secret` and read env from `water-me-secret` /
`mongodb-secret` (gitignored). **There is currently no manifest for `cleanup-worker`.**
Frontends deploy as static assets via Cloudflare (wrangler) — there is no CI
workflow; `.github/` is effectively empty (all agent guidance lives in this file).

---

## 8. Conventions (must follow)

- **Svelte 5 runes only** — `$state`, `$derived`, `$effect`, `$props()`, `$bindable`.
  No `$:`, no `export let`, no `on:click` (use `onclick`), no `createEventDispatcher`,
  no stores for new state. Prefer `$app/state` over `$app/stores`. Compose with
  `{#snippet}`/`{@render}`, not slots. When touching a legacy file, migrate the parts
  you touch (runes mode is per file — don't mix).
- **Use the semantic design tokens** (`canvas/surface/ink/ink-soft/brand/…`) and the
  `tone.ts` helpers; never hardcode colors.
- **Type API calls** through `fetchData` with paths from the generated `api.ts`; add
  new endpoints to `openapi.yaml` first.
- **i18n**: never hardcode user-visible strings — add keys to all three locales.
- **Go**: gofmt, keep HTTP handlers thin, put queries in `services/`, always scope by
  `userId`, add validation for new fields (mirror the app's expectations).
- **IDs**: plants use 24-char ObjectID hex; the detail route matches by `id` *or*
  `slug`.
- **Caching**: cache API responses and images where possible to cut network
  requests and speed up the app (the two-layer cache in §5.7 already exists — reuse
  it rather than adding new mechanisms).
- **File naming**: Svelte components are `*.svelte`; module-context files that use
  runes/state are `*.svelte.ts`.
- **Secrets**: Firebase service account goes in `backend/secret/fb.json` (gitignored),
  referenced via `GOOGLE_APPLICATION_CREDENTIALS`; never commit credentials.
- Formatting is enforced by prettier (tabs, single quotes, width 100) for TS/Svelte.

---

## 9. Known inconsistencies / gotchas

Read these before you "fix" something that looks wrong:

1. **The old `.github/copilot-instructions.md` has been removed** and its unique
   content folded into this file (see [§8 Conventions](#8-conventions-must-follow)
   and [§10 Development workflow](#10-development-workflow)). It was
   partly stale — it said `go run .` at the backend root, referenced
   `backend/validation.ts`, "CORS wildcard `*`", and a different component tree.
   None of that is accurate; use this file and the code.
2. **`openapi.yaml` is incomplete.** It documents only `/api/plants*` and
   `/api/notifications`. Implemented-but-undocumented: `/api/uploads/*`,
   `/api/notifications/tokens*`, `/api/stats`. Add them when you touch the API.
3. **`photoUrls` is not in the generated types.** The backend adds presigned
   `photoUrls` to plant responses; the app accesses it via casts
   (`(plant as { photoUrls?: string[] })`). Prefer fixing the spec + regenerating.
4. **Duplicate rate limiter.** `services/ratelimit.go` is the live implementation
   (used by `routes/plants.go`); `middlewares/ratelimit.go` is a near-identical,
   **unused** copy. `getRealIP` is likewise duplicated. Consolidate before extending.
5. **`PlantHandler` installs the rate limiter on the shared router** (`router.Use`),
   so it applies globally, not only to plant routes. Keep that in mind when adding
   routes or load tests.
6. **Dead code in `util/helpers.go`** (`HashPassword`/`VerifyPassword`/`GenerateJWT`
   stubs) and **`types/app.go`** (`User`, `SignupRequest`, `LoginRequest`, `Claims`)
   are leftovers from a pre-Firebase auth design. Don't build on them.
7. **Notification worker ignores several config fields.** `preferredTime`,
   `quietHours`, `batchingDays`, and `groupByType` exist in the model/validation but
   the worker does not honor them yet (it only uses `isEnabled`, per-type flags,
   cooldown, muted plants, tokens). This is active backlog in `TODO.md`.
8. **`profile.json` translations are unused** — `profile` is not in the i18n
   `concerns` list in `src/lib/i18n/index.ts`.
9. **`stores/language.ts` uses the legacy `writable` store** (also imported by
   `Auth.svelte`/`AppMenu`). Migrate to runes if you touch it.
10. **Push tokens are not sent to the backend yet.** `notifications/push.ts` registers
    FCM but nothing calls `POST /api/notifications/tokens`; the notifications page's
    preference checkboxes are placeholders. Tracked in `TODO.md`
    ("Frontend für Notification Configs").
11. **`API_BASE_URL` is hardcoded to production** in `app/src/lib/constants.ts` — you
    must edit it for local API development.
12. **`/api/stats` is cached 5 min in-process per replica** and `reminders` is a
    hardcoded `0` placeholder.
13. **Legal/website URLs** are centralized in `app/src/lib/constants.ts` and
    `website/src/lib/index.ts`; open external links with `openExternalLink` so Capacitor
    uses the system browser.

---

## 10. Development workflow

1. **Update the API contract first**: edit `backend/openapi.yaml` when adding or
   changing endpoints.
2. **Regenerate types**: run the `openapi-typescript` command (§5.4) so
   `app/src/lib/types/api.ts` stays in sync.
3. **Backend changes**: implement in `types/` → `validation/` → `services/` →
   `routes/`, then register the route in `routes/routes.go`.
4. **Frontend changes**: call the API through the typed client `fetchData`
   (`src/lib/auth/fetch.svelte.ts`) — never raw `fetch` for API calls.
5. **Verify**: backend `go test ./routes -v` / `go vet ./...`; app
   `npm run check && npm run lint`; website `npm run check && npm run lint`. For
   manual API testing see `backend/README.md`.

---

## 11. Recipe: adding a feature end-to-end

**New plant field / care type**
1. `backend/types/plants.go` — add struct/field (+ enum if needed).
2. `backend/validation/plant.go` — validate it.
3. `backend/services/database.go` — persist/query it (and `isEmpty*` if optional).
4. If it triggers reminders: add `GetPlantsNeeding*` + a `process*` in
   `services/notification_worker.go` + `messages/<type>.json`.
5. `backend/openapi.yaml` → regenerate `app/src/lib/types/api.ts`.
6. App: extend `types/forms.ts`, `utils/plantForm.ts` (both create + full payload),
   the matching `components/plants/forms/*.svelte`, `EDITOR_SECTIONS` if it's a new
   section, and show it on the detail page.
7. Add i18n keys to `en`, `de`, `es`; run `npm run check && npm run lint`.

**New API endpoint**
`openapi.yaml` → `types/` → `services/database.go` → `routes/*.go` →
`routes/routes.go` → regenerate types → call it via `fetchData('/api/…')` with cache
invalidation (`invalidateApiCache`) after mutations.

**New screen**
Create `src/routes/<name>/+page.svelte` using `PageHeader` + `PageContent` (+
`Scrollable`/`List` if scrollable), reachable from `BottomNav`/`AppMenu` or a card.
Set `prerender = false; ssr = false;` if it's dynamic.

---

## 12. Backlog

See **`TODO.md`** for the live backlog (Fixes / Features / Future) — it is maintained
and is the best signal of what to build next (notification config UI, delete
plant/account, log out, multi-replica notification worker, species database,
calendar/task view, shared read-only access).
