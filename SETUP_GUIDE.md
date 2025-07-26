# Speadwear-Go セットアップガイド

このガイドでは、Speadwear-Goアプリケーションを最初から動かすまでの詳細な手順を説明します。

## 📋 目次

1. [前提条件](#前提条件)
2. [環境準備](#環境準備)
3. [プロジェクトの取得](#プロジェクトの取得)
4. [データベースセットアップ](#データベースセットアップ)
5. [環境変数の設定](#環境変数の設定)
6. [アプリケーションの起動](#アプリケーションの起動)
7. [動作確認](#動作確認)
8. [API使用例](#api使用例)
9. [トラブルシューティング](#トラブルシューティング)

## 前提条件

以下のソフトウェアがインストールされている必要があります：

- **Go**: 1.21以上
- **MySQL**: 8.0以上
- **Git**: 最新版
- **Docker & Docker Compose**: （オプション、推奨）

## 環境準備

### 1. Goのインストール

#### macOS
```bash
# Homebrewを使用
brew install go

# バージョン確認
go version
# 出力例: go version go1.21.5 darwin/amd64
```

#### Linux
```bash
# 公式サイトからダウンロード
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# PATHに追加
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# バージョン確認
go version
```

#### Windows
[公式サイト](https://go.dev/dl/)からインストーラーをダウンロードして実行

### 2. MySQLのインストール

#### macOS
```bash
# Homebrewを使用
brew install mysql
brew services start mysql

# バージョン確認
mysql --version
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql

# バージョン確認
mysql --version
```

#### Windows
[MySQL公式サイト](https://dev.mysql.com/downloads/installer/)からインストーラーをダウンロード

### 3. Docker（オプション）

#### 全OS共通
[Docker Desktop](https://www.docker.com/products/docker-desktop/)をダウンロードしてインストール

## プロジェクトの取得

```bash
# 1. 作業ディレクトリに移動
cd ~/projects  # または任意のディレクトリ

# 2. リポジトリをクローン
git clone https://github.com/House-lovers7/speadwear-go.git

# 3. プロジェクトディレクトリに移動
cd speadwear-go

# 4. ディレクトリ構造を確認
ls -la
```

## データベースセットアップ

### 方法1: 手動セットアップ

```bash
# 1. MySQLにrootユーザーでログイン
mysql -u root -p

# 2. データベースとユーザーを作成
CREATE DATABASE speadwear_development CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE speadwear_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'speadwear'@'localhost' IDENTIFIED BY 'speadwear_password';
GRANT ALL PRIVILEGES ON speadwear_development.* TO 'speadwear'@'localhost';
GRANT ALL PRIVILEGES ON speadwear_test.* TO 'speadwear'@'localhost';
FLUSH PRIVILEGES;

# 3. MySQLを終了
exit
```

### 方法2: Dockerを使用

```bash
# docker-compose.ymlがあるディレクトリで実行
docker-compose up -d mysql

# データベースが起動するまで待機（約30秒）
sleep 30

# 接続確認
docker-compose exec mysql mysql -u root -p
# パスワード: root_password
```

## 環境変数の設定

### 1. .envファイルの作成

```bash
# サンプルファイルをコピー
cp .env.example .env

# .envファイルが存在しない場合は新規作成
cat > .env << 'EOF'
# サーバー設定
PORT=8080

# データベース設定
DB_HOST=localhost
DB_PORT=3306
DB_USER=speadwear
DB_PASSWORD=speadwear_password
DB_NAME=speadwear_development

# JWT設定
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24

# アップロード設定
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=5242880

# 環境
ENVIRONMENT=development
EOF
```

### 2. 環境変数の詳細説明

| 変数名 | 説明 | デフォルト値 |
|--------|------|------------|
| PORT | サーバーのポート番号 | 8080 |
| DB_HOST | MySQLホスト | localhost |
| DB_PORT | MySQLポート | 3306 |
| DB_USER | MySQLユーザー名 | speadwear |
| DB_PASSWORD | MySQLパスワード | speadwear_password |
| DB_NAME | データベース名 | speadwear_development |
| JWT_SECRET | JWT署名用の秘密鍵 | ランダムな文字列を設定 |
| JWT_EXPIRE_HOURS | トークン有効期限（時間） | 24 |
| UPLOAD_PATH | 画像アップロード先 | ./uploads |
| MAX_UPLOAD_SIZE | 最大アップロードサイズ | 5242880 (5MB) |

## アプリケーションの起動

### 方法1: セットアップスクリプトを使用（推奨）

```bash
# 1. セットアップスクリプトに実行権限を付与
chmod +x scripts/setup.sh

# 2. セットアップスクリプトを実行
./scripts/setup.sh

# 3. アプリケーションを起動
go run cmd/server/main.go
```

### 方法2: 手動セットアップ

```bash
# 1. Goモジュールの依存関係をダウンロード
go mod download

# 2. アップロードディレクトリを作成
mkdir -p uploads

# 3. データベースマイグレーションを実行
go run cmd/migrate/main.go up

# 4. アプリケーションを起動
go run cmd/server/main.go

# 正常に起動した場合の出力例：
# [GIN-debug] Listening and serving HTTP on :8080
```

### 方法3: Dockerを使用

```bash
# 1. Dockerコンテナをビルドして起動
docker-compose up --build

# バックグラウンドで起動する場合
docker-compose up -d --build

# ログを確認
docker-compose logs -f app
```

## 動作確認

### 1. ヘルスチェック

```bash
curl http://localhost:8080/health

# 期待される応答：
# {"status":"ok"}
```

### 2. ユーザー登録

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "テストユーザー",
    "email": "test@example.com",
    "password": "password123"
  }'

# 成功時の応答例：
# {
#   "user": {
#     "id": 1,
#     "name": "テストユーザー",
#     "email": "test@example.com",
#     "created_at": "2024-01-01T00:00:00Z"
#   },
#   "token": "eyJhbGciOiJIUzI1NiIs..."
# }
```

### 3. ログイン

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## API使用例

### 認証付きリクエストの例

```bash
# 1. ログインしてトークンを取得
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq -r '.token')

# 2. トークンを使ってプロフィールを取得
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me
```

### アイテムの作成（画像付き）

```bash
# フォームデータでアイテムを作成
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -F "name=お気に入りのTシャツ" \
  -F "super_item=tops" \
  -F "season=1" \
  -F "tpo=1" \
  -F "color=1" \
  -F "content=とても着心地が良いです" \
  -F "rating=5" \
  -F "image=@/path/to/image.jpg"
```

### コーディネートの作成

```bash
curl -X POST http://localhost:8080/api/v1/coordinates \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "si_shoes": 1,
    "si_bottoms": 2,
    "si_tops": 3,
    "season": 1,
    "tpo": 1,
    "rating": 5,
    "content": "春のカジュアルコーデ"
  }'
```

## トラブルシューティング

### 問題1: データベース接続エラー

```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**解決方法：**
```bash
# MySQLが起動しているか確認
sudo systemctl status mysql

# 起動していない場合
sudo systemctl start mysql

# Dockerの場合
docker-compose ps
docker-compose up -d mysql
```

### 問題2: ポート8080が使用中

```
bind: address already in use
```

**解決方法：**
```bash
# 使用中のプロセスを確認
lsof -i :8080

# 別のポートを使用（.envファイルを編集）
PORT=8081
```

### 問題3: Go依存関係エラー

```
go: missing go.sum entry
```

**解決方法：**
```bash
# go.sumを再生成
go mod tidy
go mod download
```

### 問題4: マイグレーションエラー

```
Error 1045: Access denied for user
```

**解決方法：**
```bash
# .envファイルのDB設定を確認
cat .env | grep DB_

# MySQLユーザー権限を確認
mysql -u root -p
SHOW GRANTS FOR 'speadwear'@'localhost';
```

### 問題5: 画像アップロードエラー

```
failed to create upload directory
```

**解決方法：**
```bash
# uploadsディレクトリを作成
mkdir -p uploads
chmod 755 uploads
```

## 開発用コマンド

### ホットリロード開発

```bash
# airをインストール（初回のみ）
go install github.com/cosmtrek/air@latest

# ホットリロードで開発サーバーを起動
air
```

### データベースリセット

```bash
# マイグレーションをロールバック
go run cmd/migrate/main.go down

# 再度マイグレーション
go run cmd/migrate/main.go up
```

### ログレベルの変更

```bash
# .envに追加
LOG_LEVEL=debug  # debug, info, warn, error
```

## 次のステップ

1. [API_ENDPOINTS.md](./API_ENDPOINTS.md) で全APIエンドポイントを確認
2. Postmanコレクションをインポートして詳細なテスト
3. フロントエンドアプリケーションとの連携

## サポート

問題が解決しない場合は、[GitHub Issues](https://github.com/House-lovers7/speadwear-go/issues)で報告してください。