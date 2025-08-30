# delivery-bot server (Go)

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
