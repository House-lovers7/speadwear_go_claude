#!/bin/bash

# Speadwear-Go 開発環境セットアップスクリプト

echo "=== Speadwear-Go 開発環境セットアップ ==="

# 環境変数ファイルのチェック
if [ ! -f .env ]; then
    echo "✅ .envファイルを作成します..."
    cp .env.example .env
    echo "⚠️  .envファイルを編集してデータベース設定を行ってください"
else
    echo "✅ .envファイルが既に存在します"
fi

# Dockerコンテナの起動
echo ""
echo "📦 Dockerコンテナを起動します..."
docker-compose up -d

# コンテナの起動を待つ
echo "⏳ MySQLの起動を待っています..."
sleep 10

# 依存関係のインストール
echo ""
echo "📚 Go依存関係をインストールします..."
go mod download

# データベースマイグレーション
echo ""
echo "🗄️  データベースマイグレーションを実行します..."
go run cmd/migrate/main.go up

# アップロードディレクトリの作成
echo ""
echo "📁 アップロードディレクトリを作成します..."
mkdir -p uploads/items uploads/coordinates

echo ""
echo "✅ セットアップが完了しました！"
echo ""
echo "🚀 アプリケーションを起動するには以下のコマンドを実行してください:"
echo "   go run cmd/server/main.go"
echo ""
echo "📖 実装状況の詳細は IMPLEMENTATION_STATUS.md を参照してください"