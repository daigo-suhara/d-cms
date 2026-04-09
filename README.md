# d-cms

Goで構築された柔軟なヘッドレスCMS。コンテンツスキーマ（コンテンツモデル）を自由に定義し、数字・日付・Markdownなど多様なフィールドを持つエントリを管理できます。

## 技術スタック

- **バックエンド**: Go, Gin, GORM
- **データベース**: PostgreSQL (動的フィールドはJSONBで格納)
- **ストレージ**: Cloudflare R2 (S3互換) / AWS SDK v2
- **フロントエンド**: HTMX + Alpine.js + EasyMDE (Markdownエディタ)

## セットアップ

### 1. 環境変数の設定

```bash
cp .env.example .env
```

`.env` を編集して各値を設定します。

| 変数 | 説明 |
|------|------|
| `PORT` | サーバーポート (デフォルト: `8080`) |
| `DB_URL` | PostgreSQL接続URL |
| `R2_BUCKET_NAME` | Cloudflare R2バケット名 |
| `R2_ENDPOINT` | R2エンドポイントURL |
| `R2_PUBLIC_BASE_URL` | メディアファイルの公開URL base |
| `AWS_ACCESS_KEY_ID` | R2アクセスキーID |
| `AWS_SECRET_ACCESS_KEY` | R2シークレットアクセスキー |
| `ADMIN_TOKEN` | 管理画面の認証トークン |

### 2. データベースマイグレーション

```bash
go run ./cmd/migrate
```

### 3. サーバー起動

```bash
go run ./cmd/server
```

ブラウザで `http://localhost:8080` にアクセスし、`.env` で設定した `ADMIN_TOKEN` でログインします。

## コマンド

```bash
# サーバー起動
go run ./cmd/server

# ビルド
go build -o bin/server ./cmd/server

# データベースマイグレーション
go run ./cmd/migrate

# テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/service/...
```

## アーキテクチャ

厳格な依存方向を持つレイヤードアーキテクチャを採用しています。

```
controller → service → repository → domain
```

```
cmd/
├── server/main.go          エントリポイント。依存注入とサーバー起動
└── migrate/main.go         AutoMigrateによるDBスキーマ作成

config/config.go            環境変数の読み込み

internal/
├── domain/                 純粋なGo型定義
│   ├── content_model.go    ContentModel, FieldDefinition, FieldType
│   ├── entry.go            Entry, ContentData (JSONB対応)
│   ├── media.go            Media
│   └── errors.go           センチネルエラー定義
├── repository/             インターフェース定義とGORM実装
│   ├── interface.go
│   ├── content_model_repository.go
│   ├── entry_repository.go
│   └── media_repository.go
├── service/                ビジネスロジック
│   ├── content_model_service.go  スラグ検証、エントリ存在チェック
│   ├── entry_service.go          動的フィールドのバリデーション
│   └── media_service.go          R2アップロード/削除
├── controller/             Ginハンドラー
│   ├── content_model_controller.go
│   ├── entry_controller.go
│   ├── media_controller.go
│   └── auth_controller.go
├── presenter/              レスポンス形式の振り分け
│   ├── presenter.go        Respond() — Accept headerでHTML/JSONを選択
│   ├── html_presenter.go   HTMXテンプレートレンダリング
│   └── json_presenter.go   JSONシリアライズ
├── middleware/
│   ├── auth.go             Bearer token / Cookie認証
│   └── request_id.go       X-Request-IDヘッダー付与
└── infrastructure/
    ├── database/postgres.go    GORM/PostgreSQL接続
    └── storage/r2_client.go    S3 SDK v2 ラッパー

router/router.go            ルーティング登録、DI、テンプレート設定

templates/
├── layout/base.html        共通レイアウト (HTMX, Alpine.js, EasyMDE)
├── auth/login.html
├── content_models/
│   ├── list.html
│   └── form.html           Alpine.jsによるフィールドビルダー
├── entries/
│   ├── list.html
│   └── form.html           フィールド型別パーシャルを動的レンダリング
├── media/list.html
├── error.html
└── partials/fields/
    ├── text.html
    ├── number.html
    ├── date.html
    └── markdown.html       EasyMDE統合
```

## サポートするフィールド型

| 型 | 説明 | バリデーション |
|----|------|---------------|
| `text` | テキスト入力 | 文字列チェック |
| `number` | 数値入力 | float64型チェック |
| `date` | 日付入力 | ISO 8601形式チェック |
| `markdown` | Markdownエディタ | 文字列チェック |

## API

### 管理画面 (要認証)

| Method | Path | 説明 |
|--------|------|------|
| GET | `/admin/content-models` | モデル一覧 |
| GET | `/admin/content-models/new` | モデル作成フォーム |
| POST | `/admin/content-models` | モデル作成 |
| GET | `/admin/content-models/:id/edit` | モデル編集フォーム |
| POST | `/admin/content-models/:id` | モデル更新 |
| DELETE | `/admin/content-models/:id` | モデル削除 |
| GET | `/admin/:modelSlug/entries` | エントリ一覧 |
| GET | `/admin/:modelSlug/entries/new` | エントリ作成フォーム |
| POST | `/admin/:modelSlug/entries` | エントリ作成 |
| GET | `/admin/:modelSlug/entries/:id/edit` | エントリ編集フォーム |
| POST | `/admin/:modelSlug/entries/:id` | エントリ更新 |
| DELETE | `/admin/:modelSlug/entries/:id` | エントリ削除 |
| GET | `/admin/media` | メディア一覧 |
| POST | `/admin/media/upload` | ファイルアップロード |
| DELETE | `/admin/media/:id` | ファイル削除 |

### 公開JSON API (認証不要)

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/v1/:modelSlug/entries` | エントリ一覧取得 |
| GET | `/api/v1/:modelSlug/entries/:id` | エントリ1件取得 |

### コンテンツモデル作成例 (JSON API)

```bash
curl -X POST http://localhost:8080/admin/content-models \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Blog Post",
    "slug": "blog-post",
    "fields": [
      {"name": "title", "type": "text", "required": true},
      {"name": "body", "type": "markdown", "required": true},
      {"name": "published_at", "type": "date", "required": false}
    ]
  }'
```

### エントリ作成例 (JSON API)

```bash
curl -X POST http://localhost:8080/admin/blog-post/entries \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": {
      "title": "Hello World",
      "body": "# Hello\nThis is my first post.",
      "published_at": "2024-01-15T10:00:00Z"
    }
  }'
```
