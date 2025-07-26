# Speadwear-Go テストガイド

このドキュメントは、speadwear-goプロジェクトのテスト実行方法について説明します。

## テスト構成

プロジェクトのテストは以下のカテゴリに分類されています：

- **ユニットテスト**: 個々の関数やメソッドのテスト
- **リポジトリテスト**: データベースアクセス層のテスト
- **ユースケーステスト**: ビジネスロジック層のテスト
- **統合テスト**: ハンドラー層のテスト（HTTPリクエスト/レスポンス）

## 必要な環境

- Go 1.23以上
- Docker（MySQLコンテナ用）
- Make

## セットアップ

1. 依存関係のインストール:
   ```bash
   make deps
   ```

2. テスト用データベースの作成:
   ```bash
   make docker-up
   make test-db-create
   ```

## テストの実行

### 全てのテストを実行
```bash
make test
```

### ユニットテストのみ実行
```bash
make test-unit
```

### 統合テストのみ実行
```bash
make test-integration
```

### カバレッジレポート付きでテスト実行
```bash
make test-coverage
```
レポートは `coverage/coverage.html` に生成されます。

### 特定のパッケージのテスト実行
```bash
make test-pkg PKG=internal/repository
```

### 特定のテスト関数の実行
```bash
make test-func FUNC=TestUserRepository_Create
```

### ベンチマークテストの実行
```bash
make bench
```

## テストの書き方

### ユニットテストの例

```go
func TestHashPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {
            name:     "valid password",
            password: "password123",
            wantErr:  false,
        },
        {
            name:     "empty password",
            password: "",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hash, err := HashPassword(tt.password)
            if (err != nil) != tt.wantErr {
                t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Additional assertions...
        })
    }
}
```

### リポジトリテストの例

```go
func TestUserRepository_Create(t *testing.T) {
    db := testutil.TestDB(t) // 自動的にクリーンアップされる
    repo := NewUserRepository(db)
    
    user := &domain.User{
        Name:  "Test User",
        Email: "test@example.com",
        // ...
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### 統合テストの例

```go
func TestAuthHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockUsecase := new(mockUserUsecase)
    mockUsecase.On("Login", mock.Anything, "test@example.com", "password").
        Return(&dto.AuthResponse{Token: "test-token"}, nil)
    
    handler := NewAuthHandler(config, mockUsecase)
    
    // HTTPリクエストの作成とテスト
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    handler.Login(c)
    
    assert.Equal(t, http.StatusOK, w.Code)
    mockUsecase.AssertExpectations(t)
}
```

## テストヘルパー

### TestDB
テスト用のデータベース接続を作成し、テスト終了時に自動的にクリーンアップします。

```go
db := testutil.TestDB(t)
```

### Fixtures
テストデータの作成を簡単にするヘルパー関数を提供します。

```go
fixtures := testutil.NewFixtures(t, db)
user := fixtures.CreateUser()
item := fixtures.CreateItem(user.ID)
```

## CI/CD

GitHub Actionsを使用して、プッシュおよびプルリクエスト時に自動的にテストが実行されます。

設定ファイル: `.github/workflows/test.yml`

## トラブルシューティング

### テストデータベースに接続できない
```bash
# Dockerコンテナが起動しているか確認
docker ps

# コンテナを再起動
make docker-down
make docker-up
```

### テストが遅い
並列実行を有効にする:
```go
t.Parallel()
```

### カバレッジが低い
カバレッジレポートを確認して、テストされていない部分を特定:
```bash
make test-coverage
open coverage/coverage.html
```