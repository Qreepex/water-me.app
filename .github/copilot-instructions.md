# Plants - Copilot Instructions

A plant management app with a Go backend API and a SvelteKit mobile app built with Capacitor.

## Project Structure

```
plants/
├── backend/         Go REST API + Notification Service with MongoDB
├── app/            SvelteKit + Capacitor mobile app
```

## Frontend / App

The frontend mobile app is located in the `app/` directory. It is built with SvelteKit and uses Capacitor to run on Android devices. The app communicates with the backend API for data storage and retrieval.

### Caching

If possible, cache API responses in the app to reduce network requests and improve performance.
Cache images locally on the device after downloading.

### Components

- Use reusable Svelte components for plant lists, plant details, and forms.
- Follow Svelte best practices for state management and component communication.

Use Svelte 5 features like `$state`, `$derived`, and `$effect` for reactive state management.

## Backend (Go + MongoDB)

### Running

```bash
cd backend
go run .
```

The server runs on port 8080 by default. It connects to MongoDB on startup and automatically creates collections on first use.

### Environment Variables

Create a `.env` file in `backend/`:

```bash
DATABASE_URL=mongodb://localhost:27017
MONGODB_USERNAME=test2
MONGODB_PASSWORD=test
MONGODB_DATABASE=plants

GOOGLE_APPLICATION_CREDENTIALS=./secret/fb.json

AWS_REGION=us-east-1
AWS_S3_BUCKET=your-bucket-name

PORT=8080
```

Firebase service account key should be placed at `backend/secret/fb.json`.

### Architecture

- **Authentication**: Firebase Authentication (Google sign-in) with ID token verification
- **Database**: MongoDB via official Go driver (`go.mongodb.org/mongo-driver`)
- **File Storage**: AWS S3 for plant photos with presigned URL uploads
- **Background Jobs**: Orphaned upload cleanup worker runs every 30 minutes
- **API Structure**:
  - `main.go` - Server entry point + cleanup worker
  - `routes/` - HTTP endpoint handlers (plants, uploads, notifications)
  - `services/` - MongoDB queries, Firebase auth, S3 operations, upload service
  - `middlewares/` - Auth middleware validates Firebase tokens
  - `validation/` - Input validation (Go implementation)
  - `types/` - Go type definitions with BSON tags
  - `constants/` - Collection names, limits, allowed MIME types

### Key Conventions

- All plant operations are **user-scoped** - MongoDB queries filter by `userId`
- Plant IDs use **MongoDB ObjectID** format (24-char hex strings)
- Plants also have **slugs** (human-readable, auto-generated from species name)
- Database field naming: camelCase with BSON tags (e.g., `userId`, `sunlight`)
- API field naming: camelCase (matches BSON)
- Protected endpoints require `Authorization: Bearer <firebase-id-token>` header
- CORS is configured for wildcard origins (`*`)
- User ID is extracted from Firebase token and stored in request context

### MongoDB Collections

1. **plants** - User plants with ObjectID `_id`, user-scoped by `userId`
2. **notifications** - Push notification config per user
3. **uploads** - Track S3 uploads for quota and orphan cleanup

### Validation

Both backend and frontend share validation logic:

- `backend/validation.ts` - TypeScript validation (reference implementation)
- `backend/validation/` - Go validation (enforced by backend)

When updating validation, **update both** to maintain consistency.

## Frontend (SvelteKit + Capacitor)

### Running

```bash
cd app
npm run dev          # Development server
npm run build        # Production build
npm run format       # Format with Prettier
npm run lint         # Lint with ESLint
```

### Architecture

- **Framework**: SvelteKit 2 with Svelte 5 (uses runes: `$state`, `$derived`, `$effect`)
- **Adapter**: `@sveltejs/adapter-static` - outputs to `www/` directory
- **Mobile**: Capacitor for Android (config: `capacitor.config.json`)
- **Auth**: Firebase Authentication (Google sign-in)
- **Styling**: Tailwind CSS 4
- **i18n**: Multi-language support (de, en, es) in `src/lib/i18n/translations/`

### File Structure

```
app/src/
├── lib/
│   ├── auth/              Authentication logic
│   │   └── fetch.svelte.ts    Typed API client with auth headers
│   ├── components/        Reusable Svelte components
│   │   └── PlantForms/    Form components for plant management
│   ├── i18n/              Translation files
│   ├── stores/            Svelte stores
│   ├── types/
│   │   └── api.ts         Auto-generated from OpenAPI spec
│   └── constants.ts       API_BASE_URL and other constants
└── routes/                SvelteKit pages
    ├── manage/            Plant management pages
    └── notifications/     Notification settings
```

### Key Conventions

- **Svelte 5 Runes**: Use `$state`, `$derived`, `$effect` instead of legacy Svelte 4 patterns
- **Typed API Client**: `src/lib/auth/fetch.svelte.ts` provides type-safe API calls
- **OpenAPI Types**: Auto-generated types in `src/lib/types/api.ts` from `backend/openapi.yaml`
- **API Base URL**: Configured in `src/lib/constants.ts` (`API_BASE_URL`)
- **Authentication**: Firebase ID tokens passed as Bearer tokens to backend
- **Form Components**: Plant editing split into focused form components under `PlantForms/`

### Regenerating API Types

After updating `backend/openapi.yaml`:

```bash
npx openapi-typescript ./backend/openapi.yaml -o ./app/src/lib/types/api.ts --root-types --root-types-no-schema-prefix --generate-path-params --enum --enum-values
```

### Code Style

- **Tabs for indentation** (configured in `.prettierrc`)
- **Single quotes** for strings
- **100 character line width**
- **No trailing commas**
- Component files: `*.svelte`
- Svelte module context files: `*.svelte.ts`

## Development Workflow

1. **Update API**: Edit `backend/openapi.yaml` when adding/changing endpoints
2. **Regenerate Types**: Run the openapi-typescript command above
3. **Backend Changes**: Update routes, services, validation as needed
4. **Frontend Changes**: Use typed API client from `fetch.svelte.ts`
5. **Test**: Use curl examples in `backend/README.md` for API testing

## Firebase Setup

Backend requires Firebase credentials for authentication. The app uses Google sign-in via Firebase.

## S3 Setup

Backend uses AWS S3 for photo uploads. The S3 service generates presigned URLs for direct browser uploads. CORS must be configured on the bucket to allow browser uploads.
