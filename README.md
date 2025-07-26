# Speadwear-Go

「テンションがアガル服に袖を通して1日を過ごしたい」をコンセプトにしたファッションコーディネートアプリケーション

## 📋 概要

Speadwearは、気分や季節、TPOに応じて簡単にコーディネートを選べるファッションアプリケーションです。このリポジトリは、元のRuby on RailsアプリケーションをGolangで再実装したものです。

### 特徴
- クリーンアーキテクチャによる保守性の高い設計
- RESTful APIによるバックエンド実装
- JWT認証による安全なユーザー管理
- 画像アップロード対応のアイテム・コーディネート管理
- リアルタイム通知システム

## 📚 ドキュメント

- [テーブル定義書](./docs/table_definition.md) - データベース設計の詳細
- [ER図](./docs/er_diagram.md) - エンティティ関係図
- [API仕様書（詳細版）](./docs/api_definition_detailed.md) - OpenAPI 3.0形式の完全なAPI定義
- [画面設計書](./docs/screen_design.md) - 全画面の機能とUI要素
- [システム構成図](./docs/system_architecture.md) - アーキテクチャとインフラ設計

## 🚀 主な機能

### ユーザー管理
- ユーザー登録・ログイン・ログアウト
- プロフィール編集（アバター画像アップロード対応）
- パスワード変更
- フォロー/アンフォロー機能
- ブロック機能

### アイテム管理
- 服アイテムの登録・編集・削除
- 画像アップロード（最大5MB）
- カテゴリー分類（トップス、ボトムス、アウター等）
- 季節・TPO・色による分類
- 評価機能（1-5段階）

### コーディネート機能
- 複数アイテムを組み合わせたコーディネート作成
- コーディネート画像のアップロード
- いいね機能
- コメント機能
- タイムライン表示（フォローユーザーのコーディネート）

### 検索・フィルタリング
- TPO（Time, Place, Occasion）による検索
- 季節（春夏秋冬）による絞り込み
- 色による検索
- 評価による並び替え

### ソーシャル機能
- フォロー/フォロワー管理
- いいね機能
- コメント投稿・編集・削除
- 通知システム（フォロー、いいね、コメント）

## 🛠 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **データベース**: MySQL 8.0
- **認証**: JWT (JSON Web Token)
- **画像保存**: ローカルファイルシステム（アップロード対応）
- **コンテナ**: Docker & Docker Compose

## 📁 プロジェクト構造

```
speadwear-go/
├── cmd/                    # アプリケーションのエントリーポイント
│   ├── server/            # APIサーバー
│   └── migrate/           # マイグレーションツール
├── internal/              # プライベートアプリケーションコード
│   ├── domain/           # ドメインモデルとビジネスロジック
│   │   ├── models.go     # エンティティ定義
│   │   └── constants.go  # 定数定義
│   ├── dto/              # データ転送オブジェクト
│   ├── handler/          # HTTPハンドラー（コントローラー）
│   ├── middleware/       # ミドルウェア（認証、CORS、エラー処理）
│   ├── repository/       # データアクセス層
│   ├── router/           # ルーティング設定
│   └── usecase/          # ユースケース層（ビジネスロジック）
├── pkg/                   # 公開可能なライブラリコード
│   ├── config/           # 設定管理
│   ├── database/         # データベース接続
│   └── utils/            # ユーティリティ関数（JWT、パスワード）
├── migrations/            # SQLマイグレーションファイル
├── uploads/              # アップロードされた画像の保存先
├── scripts/              # セットアップスクリプト
├── docker-compose.yml    # Docker Compose設定
├── Dockerfile           # 本番用Dockerfile
└── Dockerfile.dev       # 開発用Dockerfile
```

## 🚀 セットアップ

### 前提条件

以下のソフトウェアがインストールされている必要があります：

#### 必須
- **Go** 1.21以上
  ```bash
  # バージョン確認
  go version
  ```
- **Git**
  ```bash
  git --version
  ```

#### 推奨（Docker環境を使用する場合）
- **Docker** 20.10以上
  ```bash
  docker --version
  ```
- **Docker Compose** 2.0以上
  ```bash
  docker compose version
  ```

#### ローカル環境で実行する場合
- **MySQL** 8.0以上
  ```bash
  mysql --version
  ```
- **Air**（ホットリロード用）
  ```bash
  # インストール
  go install github.com/cosmtrek/air@latest
  ```

### 環境変数

`.env`ファイルを作成し、以下の環境変数を設定してください：

```env
# サーバー設定
PORT=8080
GIN_MODE=debug  # debug, release, test

# データベース設定
DB_HOST=localhost  # Docker使用時は "db" に変更
DB_PORT=3306
DB_USER=speadwear
DB_PASSWORD=speadwear_password
DB_NAME=speadwear_development

# JWT設定
JWT_SECRET=your-secret-key-here-please-change-in-production
JWT_EXPIRE_HOURS=24

# アップロード設定
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=5242880  # 5MB in bytes

# CORS設定
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# ログ設定
LOG_LEVEL=debug  # debug, info, warn, error
LOG_FORMAT=json  # json, text

# セッション設定
SESSION_SECRET=your-session-secret-please-change-in-production

# メール設定（オプション）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@speadwear.com

# 開発環境設定
ENABLE_SWAGGER=true
ENABLE_PROFILING=false
```

> **注意**: 本番環境では必ず`JWT_SECRET`と`SESSION_SECRET`を安全なランダム文字列に変更してください。

## 🚀 クイックスタート

### 方法1: 自動セットアップスクリプト（最速・推奨）

```bash
# 1. リポジトリのクローン
git clone https://github.com/House-lovers7/speadwear-go.git
cd speadwear-go

# 2. セットアップスクリプトの実行（Docker環境も自動構築）
chmod +x scripts/setup.sh
./scripts/setup.sh

# 3. アプリケーションの起動
go run cmd/server/main.go

# または、ホットリロード付きで起動
air
```

セットアップスクリプトは以下を自動で行います：
- `.env`ファイルの作成
- Dockerコンテナの起動（MySQL、Redis）
- Go依存関係のインストール
- データベースマイグレーション
- アップロードディレクトリの作成

### 方法2: Docker Compose（開発環境）

```bash
# 1. リポジトリのクローン
git clone https://github.com/House-lovers7/speadwear-go.git
cd speadwear-go

# 2. 環境変数ファイルの作成
cat > .env << EOF
PORT=8080
GIN_MODE=debug
DB_HOST=db
DB_PORT=3306
DB_USER=speadwear
DB_PASSWORD=password
DB_NAME=speadwear_dev
JWT_SECRET=dev-secret-key
JWT_EXPIRE_HOURS=24
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=5242880
EOF

# 3. Dockerコンテナの起動（ホットリロード付き）
docker-compose up -d

# 4. ログの確認
docker-compose logs -f app

# アプリケーションは http://localhost:8080 でアクセス可能
```

### 方法3: ローカル環境（手動セットアップ）

```bash
# 1. リポジトリのクローン
git clone https://github.com/House-lovers7/speadwear-go.git
cd speadwear-go

# 2. Go依存関係のインストール
go mod download

# 3. 環境変数の設定
cp .env.example .env  # または上記の環境変数例を参考に作成
vim .env  # DB_HOST=localhost に設定

# 4. MySQLデータベースの作成
mysql -u root -p << EOF
CREATE DATABASE IF NOT EXISTS speadwear_development CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'speadwear'@'localhost' IDENTIFIED BY 'speadwear_password';
GRANT ALL PRIVILEGES ON speadwear_development.* TO 'speadwear'@'localhost';
FLUSH PRIVILEGES;
EOF

# 5. データベースマイグレーション
go run cmd/migrate/main.go up

# 6. アップロードディレクトリの作成
mkdir -p uploads/items uploads/coordinates uploads/avatars

# 7. アプリケーションの起動
go run cmd/server/main.go

# または、ホットリロード付きで起動（Airインストール済みの場合）
air
```

### 方法4: Makefileを使用（便利コマンド）

```bash
# セットアップ
make setup

# Docker環境の起動
make docker-up

# アプリケーションの起動
make run

# テストの実行
make test

# データベースマイグレーション
make migrate

# 全てのコマンドを確認
make help
```

## 🔍 動作確認

### 1. ヘルスチェック

```bash
# アプリケーションの稼働確認
curl http://localhost:8080/health

# 期待されるレスポンス
{
  "status": "ok",
  "timestamp": "2024-01-01T09:00:00Z"
}
```

### 2. ユーザー登録とログイン

```bash
# ユーザー登録
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'

# ログイン（トークンを取得）
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq -r '.token')

echo "JWT Token: $TOKEN"
```

### 3. 認証が必要なAPIのテスト

```bash
# 現在のユーザー情報取得
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me

# ユーザー一覧取得
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/users
```

### 4. Postmanでのテスト

APIの詳細なテストには、付属の[Postmanコレクション](./docs/speadwear-postman-collection.json)を使用することを推奨します：

1. Postmanを開く
2. Import → Upload Files → `docs/speadwear-postman-collection.json`を選択
3. 環境変数で`baseUrl`を`http://localhost:8080/api/v1`に設定
4. ログインしてトークンを取得し、環境変数`token`に設定

## 📚 API仕様

詳細なAPIエンドポイントの一覧は [API_ENDPOINTS.md](./API_ENDPOINTS.md) を参照してください。

### Postmanコレクション

APIのテストを簡単に行うために、[Postmanコレクション](./docs/speadwear-postman-collection.json)を用意しています。Postmanにインポートしてご利用ください。

### 主要なエンドポイント

- **認証**: `/api/v1/auth/*`
  - POST `/auth/signup` - ユーザー登録
  - POST `/auth/login` - ログイン
  - POST `/auth/logout` - ログアウト
  - POST `/auth/refresh` - トークンリフレッシュ

- **ユーザー管理**: `/api/v1/users/*`
  - GET `/users` - ユーザー一覧
  - GET `/users/:id` - ユーザー詳細
  - PUT `/users/:id` - プロフィール更新
  - POST `/users/:id/avatar` - アバター画像アップロード

- **アイテム管理**: `/api/v1/items/*`
  - GET `/items` - アイテム一覧（フィルタリング対応）
  - POST `/items` - アイテム作成
  - PUT `/items/:id` - アイテム更新
  - DELETE `/items/:id` - アイテム削除

- **コーディネート**: `/api/v1/coordinates/*`
  - GET `/coordinates` - コーディネート一覧
  - POST `/coordinates` - コーディネート作成
  - POST `/coordinates/:id/like` - いいね
  - GET `/coordinates/timeline` - タイムライン

- **ソーシャル機能**: 
  - POST `/comments` - コメント投稿
  - POST `/users/:id/follow` - フォロー
  - GET `/notifications` - 通知一覧

## 🧪 開発

### アーキテクチャ

本プロジェクトはクリーンアーキテクチャの原則に従って設計されています：

1. **Domain層**: ビジネスロジックとエンティティを定義
2. **Usecase層**: アプリケーション固有のビジネスルール
3. **Repository層**: データの永続化を抽象化
4. **Handler層**: HTTPリクエスト/レスポンスの処理

詳細は[システム構成図](./docs/system_architecture.md)を参照してください。

### 開発用コマンド

```bash
# ホットリロードで開発
air

# 特定のパッケージのテスト
make test-pkg PKG=internal/usecase

# 特定の関数のテスト
make test-func FUNC=TestCreateUser

# カバレッジレポートの生成
make test-coverage
open coverage/coverage.html

# ベンチマークテスト
make bench

# コードフォーマット
make fmt

# Lintチェック
make lint
```

### サンプルデータの投入

開発時のテスト用にサンプルデータを投入できます：

```bash
# サンプルデータの投入スクリプト（作成予定）
go run cmd/seed/main.go

# または、SQLファイルから直接投入
mysql -h localhost -u speadwear -p speadwear_development < test/fixtures/sample_data.sql
```

### 環境別の設定

#### 開発環境
```env
GIN_MODE=debug
LOG_LEVEL=debug
ENABLE_SWAGGER=true
```

#### 本番環境
```env
GIN_MODE=release
LOG_LEVEL=info
ENABLE_SWAGGER=false
```

### ビルド

```bash
# 開発用ビルド
go build -o bin/speadwear cmd/server/main.go

# 本番用ビルド（最適化）
go build -ldflags="-s -w" -o bin/speadwear cmd/server/main.go

# クロスコンパイル（Linux用）
GOOS=linux GOARCH=amd64 go build -o bin/speadwear-linux cmd/server/main.go
```

### テスト

```bash
# 全てのテストを実行
go test ./...

# カバレッジレポート付きでテスト実行
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### コード品質

```bash
# フォーマット
go fmt ./...

# 静的解析
go vet ./...

# Linter（golangci-lintのインストールが必要）
golangci-lint run
```

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### 1. データベース接続エラー

```bash
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**解決方法：**
```bash
# MySQLサービスの状態確認
docker ps | grep mysql
# または
systemctl status mysql

# Docker環境の場合、コンテナを再起動
docker-compose restart db

# .envファイルのDB_HOSTを確認
# Docker環境: DB_HOST=db
# ローカル環境: DB_HOST=localhost
```

#### 2. ポート使用中エラー

```bash
Error: listen tcp :8080: bind: address already in use
```

**解決方法：**
```bash
# 使用中のポートを確認
lsof -i :8080

# プロセスを終了
kill -9 <PID>

# または別のポートを使用
PORT=8081 go run cmd/server/main.go
```

#### 3. マイグレーションエラー

```bash
Error: Table 'speadwear_development.users' doesn't exist
```

**解決方法：**
```bash
# マイグレーションの実行
go run cmd/migrate/main.go up

# Docker環境の場合
docker-compose exec app go run cmd/migrate/main.go up

# データベースをリセットする場合
go run cmd/migrate/main.go reset
```

#### 4. 画像アップロードエラー

```bash
Error: open ./uploads/items/xxx.jpg: permission denied
```

**解決方法：**
```bash
# ディレクトリの作成と権限設定
mkdir -p uploads/items uploads/coordinates uploads/avatars
chmod -R 755 uploads

# Docker環境の場合
docker-compose exec app mkdir -p uploads/items uploads/coordinates uploads/avatars
docker-compose exec app chmod -R 755 uploads
```

#### 5. 環境変数が読み込まれない

**解決方法：**
```bash
# .envファイルの存在確認
ls -la .env

# .envファイルの作成
cp .env.example .env

# 環境変数の確認
cat .env

# 環境変数を直接指定して起動
DB_HOST=localhost DB_PORT=3306 go run cmd/server/main.go
```

#### 6. Docker Composeが起動しない

```bash
Error: Cannot connect to the Docker daemon
```

**解決方法：**
```bash
# Docker Desktopが起動しているか確認
docker --version

# Dockerサービスの再起動
# macOS: Docker Desktopを再起動
# Linux:
sudo systemctl restart docker

# Docker Composeのバージョン確認
docker compose version
```

#### 7. Airが見つからない（ホットリロード）

```bash
air: command not found
```

**解決方法：**
```bash
# Airのインストール
go install github.com/cosmtrek/air@latest

# PATHの確認
echo $PATH

# GOPATHのbinをPATHに追加
export PATH=$PATH:$(go env GOPATH)/bin

# または直接実行
go run cmd/server/main.go
```

### デバッグモード

詳細なログを出力するには：

```bash
# 環境変数でデバッグモードを有効化
GIN_MODE=debug LOG_LEVEL=debug go run cmd/server/main.go

# Docker環境
docker-compose logs -f app
```

### ログの確認

```bash
# アプリケーションログ
tail -f logs/app.log

# Dockerログ
docker-compose logs -f

# MySQLログ
docker-compose logs -f db
```

## 🛑 環境の停止・クリーンアップ

### Docker環境の停止

```bash
# コンテナの停止
docker-compose down

# コンテナとボリュームの削除（データも削除）
docker-compose down -v

# イメージも含めて完全に削除
docker-compose down -v --rmi all
```

### ローカル環境のクリーンアップ

```bash
# ビルド成果物の削除
make clean

# データベースのリセット
make migrate-reset

# 一時ファイルの削除
rm -rf uploads/*
rm -rf logs/*
rm -rf coverage/*
```

## 📊 プロジェクト状態

- **開発状況**: アクティブ
- **最新バージョン**: v1.0.0
- **Go バージョン**: 1.21+
- **テストカバレッジ**: 約80%

## 📝 ライセンス

MIT License

## 🤝 コントリビューション

1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. Pull Requestを作成

### コーディング規約

- [Effective Go](https://golang.org/doc/effective_go.html)に従う
- コミットメッセージは[Conventional Commits](https://www.conventionalcommits.org/)形式
- テストカバレッジ80%以上を維持

## 📞 お問い合わせ

質問や提案がある場合は、[Issues](https://github.com/House-lovers7/speadwear-go/issues)を作成してください。

## 🔗 関連リンク

- [API仕様書](./API_SPECIFICATION.md)
- [画面フロー図](./SCREEN_FLOW.md)
- [セットアップガイド](./SETUP_GUIDE.md)
- [実装状況](./IMPLEMENTATION_STATUS.md)
