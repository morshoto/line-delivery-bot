# LIFF App (React + TypeScript)

## Available Scripts

- `npm run dev` – start the development server.
- `npm test` – run unit tests.
- `npm run build` – build the app for production.
- `npm run preview` – preview the built app.

## Configuration
Use Vite environment variables via `.env` files (do not commit secrets):

1) Copy `liff/.env.example` to `liff/.env` (or `.env.local`) and fill values

Required variables:

- `VITE_LIFF_ID`: Your LIFF ID
- `VITE_API_BASE`: Backend base URL (optional)
- `VITE_USE_SHARED_TOKEN`: `true`/`false`
- `VITE_SHARED_TOKEN`: Shared token value (optional)
- `VITE_OIDC_ENABLED`: `true`/`false`
- `VITE_APP_ENV`: `prod` | `stg` | `dev` (default: `dev`)

Optional for browser-only development:

- `VITE_LIFF_BROWSER_DEV`: `true` to bypass in-app requirement and use stubbed context/profile
- `VITE_DEV_GROUP_ID`: stub groupId (default: `dev-group-id`)
- `VITE_DEV_USER_ID`: stub userId (default: `U-dev-user`)
- `VITE_DEV_DISPLAY_NAME`: stub displayName (default: `Dev User`)

## Development

Install dependencies and start the development server:

```bash
npm install
npm run dev
```
