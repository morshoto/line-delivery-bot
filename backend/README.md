# delivery-bot server (Go)

- GET `/health` -> 200 `ok`
- POST `/api/scan` -> push message to group
- POST `/callback` -> verify LINE signature only

Env vars
- `LINE_CHANNEL_SECRET`: for `/callback` signature verify (Phase 1: verify-only)
- `LINE_CHANNEL_ACCESS_TOKEN`: Messaging API token to push messages
- `SHARED_TOKEN` (optional): Shared header token to accept `/api/scan`
- `PORT` (optional): default `10000`

Run locally
```
cd server
go run ./cmd/server
```

Test
```
curl -sS localhost:10000/health

curl -sS -X POST localhost:10000/api/scan \
  -H 'Content-Type: application/json' \
  -H "X-Shared-Token: $SHARED_TOKEN" \
  -d '{"group_id":"YOUR_GROUP_ID","qr_text":"https://example.com/kuroneko/123456789012","display_name":"Tester"}'
```

Notes
- `/api/scan` logs structured JSON like: `{event:scan, group_id, carrier, tracking_no, dedupe}`
- Dedupe: in-memory TTL (30m) by key `carrier+tracking_no`; duplicates append `（再スキャン）` to message.
- If `LINE_CHANNEL_ACCESS_TOKEN` is missing, server skips push and logs `push_skip`.

