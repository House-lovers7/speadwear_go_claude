.PHONY: help build run test test-verbose test-unit test-integration test-coverage test-pkg test-func bench ci-test clean migrate migrate-down migrate-reset docker-up docker-down setup test-db-create test-db-drop test-db-reset lint fmt

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプ
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application"
	@echo "  make test            - Run all tests"
	@echo "  make test-verbose    - Run tests in verbose mode"
	@echo "  make test-unit       - Run unit tests only"
	@echo "  make test-integration- Run integration tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make test-pkg PKG=path - Run tests for specific package"
	@echo "  make test-func FUNC=name - Run specific test function"
	@echo "  make bench           - Run benchmark tests"
	@echo "  make ci-test         - Run CI/CD tests"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make migrate         - Run database migrations"
	@echo "  make migrate-down    - Drop all tables"
	@echo "  make migrate-reset   - Reset database"
	@echo "  make docker-up       - Start Docker containers"
	@echo "  make docker-down     - Stop Docker containers"

# ビルド
build:
	go build -o bin/speadwear cmd/server/main.go

# 実行
run:
	go run cmd/server/main.go

# テスト
test:
	go test -v -race ./...

# テスト（詳細モード）
test-verbose:
	go test -v -race -count=1 ./...

# ユニットテストのみ
test-unit:
	go test -v -race -short ./...

# 統合テスト
test-integration:
	go test -v -race -run Integration ./...

# カバレッジレポート付きテスト
test-coverage:
	@mkdir -p coverage
	go test -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report: coverage/coverage.html"

# 特定のパッケージのテスト
test-pkg:
	go test -v -race ./$(PKG)/...

# 特定のテスト関数の実行
test-func:
	go test -v -race -run $(FUNC) ./...

# ベンチマークテスト
bench:
	go test -bench=. -benchmem ./...

# CI/CD用テスト
ci-test:
	@mkdir -p coverage
	go test -v -race -coverprofile=coverage/coverage.out -covermode=atomic ./...

# クリーン
clean:
	rm -rf bin/
	rm -rf coverage/
	go clean -cache -testcache

# マイグレーション
migrate:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-reset:
	go run cmd/migrate/main.go reset

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# 開発環境のセットアップ
setup:
	cp .env.example .env
	go mod download
	@echo "Setup complete. Please edit .env file with your database credentials."

# テスト用データベースの作成
test-db-create:
	docker exec -it speadwear-mysql mysql -uroot -proot_password -e "CREATE DATABASE IF NOT EXISTS speadwear_test;"

# テスト用データベースの削除
test-db-drop:
	docker exec -it speadwear-mysql mysql -uroot -proot_password -e "DROP DATABASE IF EXISTS speadwear_test;"

# テスト用データベースのリセット
test-db-reset: test-db-drop test-db-create

# リンター実行
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed."; \
		echo "Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# フォーマット
fmt:
	gofmt -w .
	go fmt ./...