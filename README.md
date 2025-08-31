# HakoPit

![Go Badge](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=fff&style=for-the-badge)
![chi Badge](https://img.shields.io/badge/chi-4CAF50?style=for-the-badge)
![LINE Messaging API Badge](https://img.shields.io/badge/LINE%20Messaging%20API-00C300?logo=line&logoColor=fff&style=for-the-badge)
![LIFF Badge](https://img.shields.io/badge/LIFF-00C300?logo=line&logoColor=fff&style=for-the-badge)
![pnpm Badge](https://img.shields.io/badge/pnpm-F69220?logo=pnpm&logoColor=fff&style=for-the-badge)
![Postman Badge](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=fff&style=for-the-badge)

## Description

Small Go service that accepts QR scan payloads, parses common courier/waybill strings, de‑duplicates repeat scans, and pushes formatted messages to a LINE group. Includes a callback endpoint that verifies LINE signatures.

## Branch Naming Rules

| Branch Name          | Description            | Supplemental |
| -------------------- | ---------------------- | ------------ |
| main                 | latest release         | CD action    |
| dev/main             | latest for development | CI/CD action |
| dev/{module name}    | development branch     | CI/CD action |
| hotfix/{module name} | hotfix branch          |              |
| sandbox/{anything}   | test code, etc.        |              |

-   Work is branched from each latest branch.
-   Delete working branches after merging.
-   Review as much as possible (have someone do it for you).
-   Build, deploy, etc. are discussed separately.

## Usage

### LIFF

Copy `liff/.env.example` to `liff/.env` (or `.env.local`) and fill values. Run locally

```bash
# Install dependencies
pnpm install
# Start the production server (after build).
pnpm run dev
```

### Backend

-   Dedupe: in‑memory TTL (30m) keyed by `carrier+tracking_no`; duplicates append `（再スキャン）` to the message.
-   If `LINE_CHANNEL_ACCESS_TOKEN` is missing, the server skips push and logs `push_skip`.
-   Behind a proxy: `go env -w GOPROXY=https://proxy.golang.org,direct`

Run locally

```bash
# Fetch dependencies (generates go.sum)
go mod tidy
# Start the server
go run ./cmd/server
```

**Postman**

-   Collection: `backend/data/postman/line-delivery-bot.postman_collection.json`
-   Environment: `backend/data/postman/line-delivery-bot.postman_environment.json`
