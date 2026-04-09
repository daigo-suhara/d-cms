# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリのコードを操作する際のガイドラインを提供します。

## プロジェクト概要

Goで構築された柔軟なヘッドレスCMS。ユーザーがコンテンツスキーマ（コンテンツモデル）を自由に定義し、数字、日付、Markdownエディタなどの多様なフィールドを持つエントリを管理できます。

- **バックエンド**: Go, **Gin** (HTTPルーター), **GORM** (PostgreSQL)
- **データベース**: PostgreSQL (**JSONB** を使用して動的フィールドを格納)
- **ストレージ**: **Cloudflare R2** (S3互換) / AWS SDK v2
- **フロントエンド**: **HTMX** (サーバーサイドレンダリング) + **Alpine.js** (エディタ等のインタラクション)

## コマンド

```bash
# サーバーの起動
go run ./cmd/server

# ビルド
go build -o bin/server ./cmd/server

# データベースマイグレーション
go run ./cmd/migrate

# 全テストの実行
go test ./...

# 特定パッケージのテスト実行
go test ./internal/service/...

# 単一テストの実行
go test ./internal/service -run TestEntryService_Validate

# リンター実行 (golangci-lint)
golangci-lint run
```

## アーキテクチャ

本プロジェクトは厳格な依存方向を持つレイヤードアーキテクチャを採用しています：
**controller → service → repository → domain**

```
cmd/server/main.go          エントリポイント。依存注入とサーバー起動
config/config.go            環境変数の読み込み

internal/domain/            純粋なGo型定義（ContentModel, FieldDefinition, Entry）
internal/repository/        インターフェース定義とGORM/JSONBによる実装
internal/service/           ビジネスロジック。動的フィールドのバリデーション（数字、日付、Markdown）
internal/controller/        Ginハンドラー。入力バリデーションとサービスの呼び出し
internal/presenter/         レスポンスの分岐（JSON または HTML/HTMX を Respond() で集約）
internal/middleware/        認証、リクエストIDなどのミドルウェア
internal/infrastructure/    GORM/PostgreSQL接続およびS3/R2クライアントの設定

router/router.go            ルーティング登録とミドルウェアの適用
templates/                  HTMXテンプレート（レイアウト、フィールド型ごとのパーシャル）
```

### 重要な設計方針

- **動的コンテンツスキーマ**: 固定の構造体ではなく、`ContentModel` でフィールドを定義します。実際のデータは PostgreSQL の `JSONB` カラム（`Entry.Content`）に格納し、可変なフィールドに対応します。
- **プレゼンターパターン**: レスポンスロジックを `internal/presenter/` に集約。コントローラーは `presenter.Respond()` を呼び出し、リクエストの Accept ヘッダーに基づいて JSON（API用）か HTML（管理画面用）を自動的に振り分けます。
- **Markdownの統合**: 管理画面に Markdown エディタ（EasyMDE 等）を内蔵。データは生の Markdown として保存し、必要に応じてサーバーサイドまたはフロントエンドでパースします。
- **リポジトリインターフェース**: サービス層は `internal/repository/interface.go` で定義されたインターフェースに依存し、テスト時のモック作成を容易にします。

## サポートされるフィールド型

`FieldDefinition` で設定可能な型：
- `text`: 標準的な文字列入力
- `number`: 数値バリデーション（float64/int）
- `date`: 日付/時刻（ISO 8601 形式）
- `markdown`: Markdownエディタによる長文入力

## コーディング規約

- **エラーハンドリング**: Service/Repository 層ではコンテキストを付与してエラーをラップし、Controller 層で適切に処理・ログ出力します。
- **バリデーション**: `Entry` の内容が `ContentModel` の定義（必須チェック、型チェック）に適合するかどうかは Service 層が責任を持ちます。
- **HTML/HTMX**: フォームの部品化を維持するため、`templates/partials/fields/` 内に各フィールド型専用のテンプレートを作成します。
- **フォーマット**: 標準の `gofmt` に従い、Goの慣習（Effective Go）を遵守します。

## 設定
`.env.example` を `.env` にコピーして使用してください。
主な変数: `DB_URL`, `R2_BUCKET_NAME`, `R2_ENDPOINT`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
