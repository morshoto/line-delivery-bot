# delivery-bot server (Go)

- GET `/health` -> 200 "ok"
- env:
  - `LINE_CHANNEL_SECRET` (later)
  - `LINE_CHANNEL_ACCESS_TOKEN` (later)
  - `SHARED_TOKEN` (optional)

## Local run

```bash
cd server
# First time: resolve modules
go mod tidy
# Run server
go run ./cmd/server
# -> listening on :10000
# Health check
curl -i http://localhost:10000/health
```

> Note: `go mod tidy` needs network to fetch dependencies.

