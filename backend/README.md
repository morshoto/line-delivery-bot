# delivery-bot server (Go)

**Prerequisites**

-   Go 1.22+

**Environment variables**

-   `LINE_CHANNEL_SECRET`: for `/callback` signature verification
-   `LINE_CHANNEL_ACCESS_TOKEN`: Messaging API token to push messages
-   `SHARED_TOKEN` (optional): shared header token required by `/api/scan`
-   `PORT` (optional): default `10000`

Run locally

```bash
# Fetch dependencies (generates go.sum)
go mod tidy
# Start the server
go run ./cmd/server
```

Testing on local

```bash
# Run tests using go
go test ./... -v
# Run tests using gotestfmt
go test -json ./... -count=1 | gotestfmt -hide=empty-packages
```

**Notes**

-   Dedupe: in‑memory TTL (30m) keyed by `carrier+tracking_no`; duplicates append `（再スキャン）` to the message.
-   If `LINE_CHANNEL_ACCESS_TOKEN` is missing, the server skips push and logs `push_skip`.
-   Behind a proxy: `go env -w GOPROXY=https://proxy.golang.org,direct`

**Postman**

-   Collection: `backend/data/postman/line-delivery-bot.postman_collection.json`
-   Environment: `backend/data/postman/line-delivery-bot.postman_environment.json`
