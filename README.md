# d-cms

Go で構築された柔軟なヘッドレス CMS。コンテンツスキーマ（コンテンツモデル）を自由に定義し、テキスト・数値・日付・Markdown など多様なフィールドを持つエントリを管理できます。管理画面と JSON REST API の両方を提供します。

## 技術スタック

| レイヤー | 採用技術 |
|----------|----------|
| バックエンド | Go, [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io) |
| データベース | PostgreSQL（動的フィールドは **JSONB** で格納） |
| ストレージ | Cloudflare R2 / AWS S3 互換（AWS SDK v2） |
| 管理画面 | [HTMX](https://htmx.org) + [Alpine.js](https://alpinejs.dev) + [Tailwind CSS](https://tailwindcss.com) (CDN) |
| Markdown | [EasyMDE](https://github.com/Ionaru/easy-markdown-editor) |

## セットアップ

### 1. 環境変数の設定

```bash
cp .env.example .env
```

| 変数 | 必須 | 説明 |
|------|------|------|
| `PORT` | | サーバーポート（デフォルト: `8080`） |
| `DB_URL` | ✓ | PostgreSQL 接続 URL |
| `ADMIN_TOKEN` | ✓ | 管理画面ログイン用トークン |
| `R2_BUCKET_NAME` | | Cloudflare R2 バケット名 |
| `R2_ENDPOINT` | | R2 エンドポイント URL |
| `R2_PUBLIC_BASE_URL` | | メディアファイルの公開ベース URL |
| `AWS_ACCESS_KEY_ID` | | R2 アクセスキー ID |
| `AWS_SECRET_ACCESS_KEY` | | R2 シークレットアクセスキー |

### 2. データベースマイグレーション

```bash
go run ./cmd/migrate
```

### 3. サーバー起動

```bash
go run ./cmd/server
```

`http://localhost:8080` にアクセスし、`ADMIN_TOKEN` でログインします。

## コマンド

```bash
# 開発サーバー起動
go run ./cmd/server

# プロダクションビルド
go build -o bin/server ./cmd/server

# DBマイグレーション
go run ./cmd/migrate

# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/service/...

# Lint
golangci-lint run
```

## アーキテクチャ

厳格な依存方向を持つレイヤードアーキテクチャを採用しています。

```
controller → service → repository → domain
```

```
cmd/
├── server/main.go              エントリポイント・依存注入・サーバー起動
└── migrate/main.go             AutoMigrate によるDBスキーマ作成

config/config.go                環境変数の読み込み

internal/
├── domain/                     純粋な Go 型定義（外部依存なし）
│   ├── content_model.go        ContentModel, FieldDefinition, FieldType
│   ├── entry.go                Entry, ContentData（JSONB 対応 map）
│   ├── media.go                Media
│   ├── api_key.go              APIKey
│   └── errors.go               センチネルエラー定義
├── repository/                 インターフェース定義と GORM 実装
│   ├── interface.go
│   ├── content_model_repository.go
│   ├── entry_repository.go
│   ├── media_repository.go
│   └── api_key_repository.go
├── service/                    ビジネスロジック
│   ├── content_model_service.go    スラグ検証・エントリ存在チェック
│   ├── entry_service.go            動的フィールドのバリデーション
│   ├── media_service.go            R2 アップロード / 削除
│   └── api_key_service.go          APIキー生成（dcms_ プレフィックス）・検証
├── controller/                 Gin ハンドラー
│   ├── content_model_controller.go
│   ├── entry_controller.go
│   ├── media_controller.go
│   ├── api_key_controller.go
│   └── auth_controller.go
├── presenter/                  レスポンス形式の振り分け
│   ├── presenter.go            Respond() — Accept ヘッダーで HTML/JSON を選択
│   ├── html_presenter.go       HTMX テンプレートレンダリング
│   └── json_presenter.go       JSON シリアライズ
├── middleware/
│   ├── auth.go                 管理画面認証（Bearer token / Cookie）
│   ├── api_auth.go             API キー認証（Bearer token）
│   └── request_id.go           X-Request-ID ヘッダー付与
└── infrastructure/
    ├── database/postgres.go    GORM / PostgreSQL 接続
    └── storage/r2_client.go    S3 SDK v2 ラッパー

router/router.go                ルーティング登録・DI・テンプレート設定

templates/                          → [UI コンポーネントクラス一覧](templates/COMPONENTS.md)
├── layout/base.html            共通レイアウト（折りたたみサイドバー、Alpine.js）
├── auth/login.html             ログインページ（スタンドアロン）
├── content_models/
│   ├── list.html
│   └── form.html               Alpine.js によるフィールドビルダー
├── entries/
│   ├── list.html
│   └── form.html               フィールド型別パーシャルを動的レンダリング
├── media/list.html
├── api_keys/list.html          APIキー管理 + APIリファレンス
├── error.html
└── partials/fields/
    ├── text.html
    ├── number.html
    ├── date.html
    └── markdown.html           EasyMDE 統合
```

## サポートするフィールド型

| 型 | 説明 | サーバーバリデーション |
|----|------|----------------------|
| `text` | テキスト入力 | 文字列チェック |
| `number` | 数値入力 | float64 型チェック |
| `date` | 日付入力 | ISO 8601 形式チェック |
| `markdown` | Markdown エディタ | 文字列チェック |

## APIキー

管理画面の **「APIキー管理」** ページからキーを発行します。キーは生成時に一度だけ表示されます。

```
Authorization: Bearer dcms_<48桁のランダム文字列>
```

## REST API リファレンス

**全エンドポイントに `Authorization: Bearer <api-key>` ヘッダーが必要です。**

### コンテンツモデル

| Method | Path | 説明 |
|--------|------|------|
| `GET` | `/api/v1/models` | モデル一覧 |
| `GET` | `/api/v1/models/:slug` | モデルのスキーマ取得 |

### エントリ

| Method | Path | 説明 |
|--------|------|------|
| `GET` | `/api/v1/:modelSlug/entries` | エントリ一覧 |
| `GET` | `/api/v1/:modelSlug/entries/:id` | エントリ取得 |
| `POST` | `/api/v1/:modelSlug/entries` | エントリ作成 |
| `PUT` | `/api/v1/:modelSlug/entries/:id` | エントリ更新 |
| `DELETE` | `/api/v1/:modelSlug/entries/:id` | エントリ削除 |

### 使用例

**エントリ一覧の取得:**

```bash
curl https://your-domain/api/v1/blog-post/entries \
  -H "Authorization: Bearer dcms_xxxxxxxxxxxx..."
```

**エントリの作成:**

```bash
curl -X POST https://your-domain/api/v1/blog-post/entries \
  -H "Authorization: Bearer dcms_xxxxxxxxxxxx..." \
  -H "Content-Type: application/json" \
  -d '{
    "content": {
      "title": "Hello World",
      "body": "# Hello\nThis is my first post.",
      "published_at": "2024-01-15T10:00:00Z"
    }
  }'
```

**エントリの更新:**

```bash
curl -X PUT https://your-domain/api/v1/blog-post/entries/1 \
  -H "Authorization: Bearer dcms_xxxxxxxxxxxx..." \
  -H "Content-Type: application/json" \
  -d '{"content": {"title": "Updated Title"}}'
```

**エントリの削除:**

```bash
curl -X DELETE https://your-domain/api/v1/blog-post/entries/1 \
  -H "Authorization: Bearer dcms_xxxxxxxxxxxx..."
```

## 管理画面ルート（ブラウザ用）

管理画面は `ADMIN_TOKEN` によるセッション Cookie 認証を使用します。

| Path | 説明 |
|------|------|
| `GET /admin/login` | ログイン |
| `GET /admin/content-models` | コンテンツモデル一覧 |
| `GET /admin/:modelSlug/entries` | エントリ一覧 |
| `GET /admin/media` | メディアライブラリ |
| `GET /admin/api-keys` | APIキー管理 |
