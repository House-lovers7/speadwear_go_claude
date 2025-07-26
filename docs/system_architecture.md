# Speadwear システム構成図

## 概要
このドキュメントでは、Speadwearアプリケーションのシステムアーキテクチャを図解します。
クリーンアーキテクチャの原則に基づいた設計を採用しています。

## 1. システム全体構成

```mermaid
graph TB
    subgraph "Frontend"
        WEB[Web Application<br/>React/Next.js]
        MOBILE[Mobile Apps<br/>React Native]
    end
    
    subgraph "API Gateway"
        NGINX[Nginx<br/>Reverse Proxy]
    end
    
    subgraph "Backend Services"
        API[Speadwear API<br/>Go + Gin]
        AUTH[Auth Service<br/>JWT]
    end
    
    subgraph "Data Layer"
        MYSQL[(MySQL 8.0<br/>Primary DB)]
        REDIS[(Redis<br/>Cache/Session)]
        S3[S3 Compatible<br/>Object Storage]
    end
    
    subgraph "External Services"
        SMTP[SMTP Server<br/>Email]
        PUSH[Push Notification<br/>Service]
    end
    
    WEB --> NGINX
    MOBILE --> NGINX
    NGINX --> API
    API --> AUTH
    API --> MYSQL
    API --> REDIS
    API --> S3
    API --> SMTP
    API --> PUSH
```

## 2. クリーンアーキテクチャ層構造

```mermaid
graph TD
    subgraph "外部層"
        UI[UI/Presentation]
        DB[Database]
        WEB[Web Framework]
        EXT[External Services]
    end
    
    subgraph "インターフェース層"
        HANDLER[Handler/Controller]
        REPO[Repository Implementation]
        PRES[Presenter]
    end
    
    subgraph "アプリケーション層"
        USECASE[Use Cases]
        DTO[Data Transfer Objects]
    end
    
    subgraph "ドメイン層"
        ENTITY[Entities]
        DOMAIN[Domain Services]
        REPO_IF[Repository Interface]
    end
    
    UI --> HANDLER
    HANDLER --> USECASE
    USECASE --> REPO_IF
    REPO_IF --> REPO
    REPO --> DB
    USECASE --> DOMAIN
    DOMAIN --> ENTITY
    HANDLER --> DTO
    PRES --> UI
```

### 各層の責務

#### ドメイン層 (Domain Layer)
- **場所**: `/internal/domain/`
- **責務**: ビジネスロジックとエンティティの定義
- **含まれるもの**:
  - エンティティ（User, Item, Coordinate等）
  - ビジネスルール
  - ドメインサービス
  - リポジトリインターフェース

#### アプリケーション層 (Application Layer)
- **場所**: `/internal/usecase/`
- **責務**: アプリケーション固有のビジネスルール
- **含まれるもの**:
  - ユースケース実装
  - アプリケーションサービス
  - DTOの定義

#### インターフェース層 (Interface Layer)
- **場所**: `/internal/handler/`, `/internal/repository/`
- **責務**: 外部との接続
- **含まれるもの**:
  - HTTPハンドラー
  - リポジトリ実装
  - ミドルウェア

#### インフラストラクチャ層 (Infrastructure Layer)
- **場所**: `/pkg/`
- **責務**: 技術的な実装詳細
- **含まれるもの**:
  - データベース接続
  - 外部サービス連携
  - 設定管理

## 3. 技術スタック詳細

```mermaid
graph LR
    subgraph "Backend Stack"
        GO[Go 1.21+]
        GIN[Gin Web Framework]
        GORM[GORM ORM]
        JWT[JWT-Go]
        VALIDATOR[Go Validator]
    end
    
    subgraph "Database Stack"
        MYSQL[MySQL 8.0]
        REDIS[Redis 7.0]
    end
    
    subgraph "Infrastructure"
        DOCKER[Docker]
        COMPOSE[Docker Compose]
        NGINX[Nginx]
    end
    
    subgraph "Development Tools"
        AIR[Air - Hot Reload]
        MIGRATE[Go-Migrate]
        MOCKGEN[Mockgen]
        TEST[Go Test]
    end
```

### バックエンド技術選定理由

| 技術 | 選定理由 |
|------|----------|
| Go | 高いパフォーマンス、並行処理、型安全性 |
| Gin | 軽量で高速、豊富なミドルウェア |
| GORM | Go標準のORM、マイグレーション機能 |
| MySQL | 信頼性、ACID準拠、豊富な機能 |
| Redis | 高速キャッシュ、セッション管理 |
| JWT | ステートレス認証、スケーラビリティ |

## 4. デプロイメント構成

### 開発環境

```mermaid
graph TB
    subgraph "Developer Machine"
        CODE[Source Code]
        DOCKER_DEV[Docker Desktop]
    end
    
    subgraph "Docker Compose"
        APP_DEV[App Container<br/>Hot Reload]
        DB_DEV[MySQL Container]
        REDIS_DEV[Redis Container]
    end
    
    CODE --> DOCKER_DEV
    DOCKER_DEV --> APP_DEV
    APP_DEV --> DB_DEV
    APP_DEV --> REDIS_DEV
```

### 本番環境

```mermaid
graph TB
    subgraph "Load Balancer"
        ALB[Application Load Balancer]
    end
    
    subgraph "App Servers"
        APP1[App Instance 1]
        APP2[App Instance 2]
        APP3[App Instance 3]
    end
    
    subgraph "Data Layer"
        RDS[(RDS MySQL<br/>Multi-AZ)]
        ELASTICACHE[(ElastiCache<br/>Redis Cluster)]
        S3_PROD[S3 Bucket<br/>Static Files]
    end
    
    subgraph "Monitoring"
        CW[CloudWatch]
        XRAY[X-Ray]
    end
    
    ALB --> APP1
    ALB --> APP2
    ALB --> APP3
    APP1 --> RDS
    APP2 --> RDS
    APP3 --> RDS
    APP1 --> ELASTICACHE
    APP2 --> ELASTICACHE
    APP3 --> ELASTICACHE
    APP1 --> S3_PROD
    APP1 --> CW
    APP1 --> XRAY
```

## 5. データフロー図

### 認証フロー

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant U as UseCase
    participant R as Repository
    participant DB as Database
    participant J as JWT Service
    
    C->>H: POST /auth/login
    H->>U: Login(email, password)
    U->>R: GetUserByEmail(email)
    R->>DB: SELECT * FROM users
    DB-->>R: User data
    R-->>U: User entity
    U->>U: Verify password
    U->>J: GenerateToken(userID)
    J-->>U: JWT token
    U-->>H: LoginResult
    H-->>C: 200 OK + token
```

### コーディネート作成フロー

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant M as Middleware
    participant U as UseCase
    participant R as Repository
    participant S as Storage
    participant DB as Database
    
    C->>H: POST /coordinates
    H->>M: Authenticate
    M-->>H: UserID
    H->>U: CreateCoordinate(data)
    U->>R: GetItemsByIDs(itemIDs)
    R->>DB: SELECT items
    DB-->>R: Items
    R-->>U: Items
    U->>U: Validate ownership
    U->>S: UploadImage(file)
    S-->>U: Image URL
    U->>R: CreateCoordinate(coordinate)
    R->>DB: INSERT coordinate
    DB-->>R: Created ID
    R-->>U: Coordinate
    U-->>H: Result
    H-->>C: 201 Created
```

## 6. セキュリティアーキテクチャ

```mermaid
graph TD
    subgraph "Security Layers"
        WAF[Web Application Firewall]
        HTTPS[HTTPS/TLS 1.3]
        AUTH[Authentication<br/>JWT]
        AUTHZ[Authorization<br/>Role-Based]
        VALID[Input Validation]
        SANITIZE[Output Sanitization]
        AUDIT[Audit Logging]
    end
    
    subgraph "Security Measures"
        RATE[Rate Limiting]
        CORS[CORS Policy]
        CSRF[CSRF Protection]
        XSS[XSS Prevention]
        SQL[SQL Injection Prevention]
        CRYPT[Encryption at Rest]
    end
    
    WAF --> HTTPS
    HTTPS --> AUTH
    AUTH --> AUTHZ
    AUTHZ --> VALID
    VALID --> SANITIZE
    SANITIZE --> AUDIT
```

### セキュリティ実装

| 脅威 | 対策 |
|------|------|
| 認証回避 | JWT + リフレッシュトークン |
| SQLインジェクション | パラメータバインディング |
| XSS | 出力エスケープ、CSP |
| CSRF | CSRFトークン |
| DDoS | レート制限、WAF |
| データ漏洩 | 暗号化、アクセス制御 |

## 7. パフォーマンス最適化

### キャッシュ戦略

```mermaid
graph LR
    subgraph "Cache Layers"
        CDN[CDN Cache<br/>Static Assets]
        REDIS[Redis Cache<br/>API Response]
        APP[Application Cache<br/>In-Memory]
        DB[Database Cache<br/>Query Cache]
    end
    
    CDN --> REDIS
    REDIS --> APP
    APP --> DB
```

### 最適化技術

| 領域 | 技術 |
|------|------|
| データベース | インデックス最適化、N+1問題対策 |
| API | ページネーション、遅延読み込み |
| キャッシュ | Redis、HTTP キャッシュヘッダー |
| 画像 | 圧縮、WebP、CDN配信 |
| 並行処理 | Goroutineによる非同期処理 |

## 8. 監視・運用

### 監視スタック

```mermaid
graph TB
    subgraph "Application"
        APP[Speadwear API]
        METRICS[Metrics Endpoint]
        HEALTH[Health Check]
    end
    
    subgraph "Monitoring"
        PROM[Prometheus]
        GRAF[Grafana]
        ALERT[AlertManager]
    end
    
    subgraph "Logging"
        LOG[Application Logs]
        ELK[ELK Stack]
        S3_LOG[S3 Log Archive]
    end
    
    APP --> METRICS
    APP --> HEALTH
    APP --> LOG
    METRICS --> PROM
    PROM --> GRAF
    PROM --> ALERT
    LOG --> ELK
    ELK --> S3_LOG
```

### 監視項目

| カテゴリ | メトリクス |
|----------|------------|
| アプリケーション | レスポンスタイム、エラー率、スループット |
| インフラ | CPU使用率、メモリ使用率、ディスクI/O |
| ビジネス | アクティブユーザー数、投稿数、エンゲージメント率 |

## 9. 災害復旧計画

### バックアップ戦略

```mermaid
graph LR
    subgraph "Backup Schedule"
        DAILY[Daily Backup<br/>Incremental]
        WEEKLY[Weekly Backup<br/>Full]
        MONTHLY[Monthly Backup<br/>Archive]
    end
    
    subgraph "Storage"
        S3_BACKUP[S3 Backup<br/>Cross-Region]
        GLACIER[Glacier<br/>Long-term]
    end
    
    DAILY --> S3_BACKUP
    WEEKLY --> S3_BACKUP
    MONTHLY --> GLACIER
```

### RTO/RPO目標

| 指標 | 目標値 |
|------|--------|
| RTO (Recovery Time Objective) | 4時間 |
| RPO (Recovery Point Objective) | 1時間 |

## 10. スケーリング戦略

### 水平スケーリング

```mermaid
graph TB
    subgraph "Auto Scaling"
        AS[Auto Scaling Group]
        METRIC[CloudWatch Metrics]
        POLICY[Scaling Policy]
    end
    
    subgraph "Scaling Triggers"
        CPU[CPU > 70%]
        MEM[Memory > 80%]
        REQ[Request Count]
        LAT[Latency > 500ms]
    end
    
    CPU --> METRIC
    MEM --> METRIC
    REQ --> METRIC
    LAT --> METRIC
    METRIC --> POLICY
    POLICY --> AS
```

### スケーリング計画

| フェーズ | ユーザー数 | インフラ構成 |
|----------|-----------|--------------|
| 初期 | 〜10,000 | 2 instances + RDS |
| 成長期 | 〜100,000 | 4 instances + RDS Read Replica |
| 拡大期 | 〜1,000,000 | Auto Scaling + Aurora |

## まとめ

Speadwearのシステムアーキテクチャは、以下の原則に基づいて設計されています：

1. **クリーンアーキテクチャ**: ビジネスロジックの独立性
2. **マイクロサービス指向**: 将来的な分割を考慮
3. **スケーラビリティ**: 水平スケーリング対応
4. **高可用性**: 冗長性とフェイルオーバー
5. **セキュリティ**: 多層防御
6. **監視可能性**: 包括的なメトリクス収集