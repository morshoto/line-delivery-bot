# LIFF App (Next.js + TypeScript)

## Available Scripts

- `npm run dev` – start the development server.
- `npm test` – run unit tests.
- `npm run build` – build the app for production.
- `npm run start` – start the production server (after build).

## Configuration

Use Next.js environment variables via `.env` files (do not commit secrets). In client-side code, only variables prefixed with `NEXT_PUBLIC_` are exposed.

1. Copy `liff/.env.example` to `liff/.env` (or `.env.local`) and fill values

Required variables:

- `NEXT_PUBLIC_LIFF_ID`: Your LIFF ID
- `NEXT_PUBLIC_API_BASE`: Backend base URL (optional)
- `NEXT_PUBLIC_USE_SHARED_TOKEN`: `true`/`false`
- `NEXT_PUBLIC_SHARED_TOKEN`: Shared token value (optional)
- `NEXT_PUBLIC_OIDC_ENABLED`: `true`/`false`
- `NEXT_PUBLIC_APP_ENV`: `prod` | `stg` | `dev` (default: `dev`)

Optional for browser-only development:

- `NEXT_PUBLIC_LIFF_BROWSER_DEV`: `true` to bypass in-app requirement and use stubbed context/profile
- `NEXT_PUBLIC_DEV_GROUP_ID`: stub groupId (default: `dev-group-id`)
- `NEXT_PUBLIC_DEV_USER_ID`: stub userId (default: `U-dev-user`)
- `NEXT_PUBLIC_DEV_DISPLAY_NAME`: stub displayName (default: `Dev User`)

## Development

Install dependencies and start the development server:

```bash
npm install
npm run dev
```
