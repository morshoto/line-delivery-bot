Postman assets for line-delivery-bot backend

- Import `backend/data/postman/line-delivery-bot.postman_collection.json` into Postman.
- Import `backend/data/postman/line-delivery-bot.postman_environment.json` and select the environment.
- Ensure the server runs locally on `http://localhost:10000` or change `baseUrl` in the environment.

Included requests
- GET `/health` — expects 200 `ok`.
- POST `/api/scan` — uses `X-Shared-Token` and JSON body.
- POST `/callback` — includes `X-Line-Signature` header (verify-only).

Environment variables
- `baseUrl`: default `http://localhost:10000`
- `sharedToken`: optional shared token for `/api/scan`
- `groupId`, `qrText`, `displayName`: payload fields for `/api/scan`
- `lineSignature`: placeholder value for `/callback`

