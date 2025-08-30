# HakoPit

![Go Badge](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=fff&style=for-the-badge)
![chi Badge](https://img.shields.io/badge/chi-4CAF50?style=for-the-badge)
![LINE Messaging API Badge](https://img.shields.io/badge/LINE%20Messaging%20API-00C300?logo=line&logoColor=fff&style=for-the-badge)
![Postman Badge](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=fff&style=for-the-badge)

小さな Go サービスです。QR で読み取った荷物情報を解析し、LINE グループへ分かりやすく通知します。重複スキャンは 30 分間抑止し、LINE のコールバックも署名を検証して受け付けます。

## 特徴
- 主要な宅配会社の伝票番号を解析
- 同じ伝票番号の再スキャンを検知して通知を抑制
- LINE Messaging API を使ってグループへメッセージを送信
- Postman コレクションを同梱

## 動作環境
- Go 1.22 以上

### 必要な環境変数
| 変数 | 説明 |
| ---- | ---- |
| `LINE_CHANNEL_SECRET` | `/callback` の署名検証に使用 |
| `LINE_CHANNEL_ACCESS_TOKEN` | メッセージ送信に使用するトークン |
| `SHARED_TOKEN` (任意) | `/api/scan` を呼び出す際の共有ヘッダトークン |
| `PORT` (任意) | サーバーのポート番号。デフォルト `10000` |

## ローカルでの実行手順
```bash
# リポジトリのルートから
cd backend
# 依存関係の取得
go mod tidy
# サーバーを起動
go run ./cmd/server
```

## ロギングと挙動
- `/api/scan` は `{event, group_id, carrier, tracking_no, dedupe}` を含む JSON を出力します。
- 重複検出は `carrier+tracking_no` をキーにメモリ上で 30 分保持し、再スキャン時はメッセージに `（再スキャン）` を付加します。
- `LINE_CHANNEL_ACCESS_TOKEN` が未設定の場合は送信をスキップし `push_skip` をログに記録します。
- プロキシ環境では `go env -w GOPROXY=https://proxy.golang.org,direct` を設定してください。

## Postman
- コレクション: `backend/data/postman/line-delivery-bot.postman_collection.json`
- 環境: `backend/data/postman/line-delivery-bot.postman_environment.json`

## ブランチ運用ルール
| ブランチ名 | 役割 | 補足 |
| ---------- | ---- | ---- |
| `main` | 最新リリース | CD |
| `dev/main` | 開発の最新 | CI/CD |
| `dev/{module}` | 機能開発用 | CI/CD |
| `hotfix/{module}` | ホットフィックス用 | |
| `sandbox/{anything}` | 試験コードなど | |

- 作業は各ブランチの最新から切り出します。
- マージ後は作業ブランチを削除します。
- 可能な限りレビューを行います。
- ビルドやデプロイに関しては別途相談します。
