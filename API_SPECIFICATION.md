# Speadwear API仕様書

## 目次

1. [概要](#概要)
2. [認証](#認証)
3. [共通仕様](#共通仕様)
4. [エラーハンドリング](#エラーハンドリング)
5. [API詳細](#api詳細)
   - [認証API](#認証api)
   - [ユーザーAPI](#ユーザーapi)
   - [アイテムAPI](#アイテムapi)
   - [コーディネートAPI](#コーディネートapi)
   - [ソーシャルAPI](#ソーシャルapi)
6. [定数定義](#定数定義)

## 概要

### ベースURL
```
https://api.speadwear.com/api/v1
```
開発環境: `http://localhost:8080/api/v1`

### 通信プロトコル
- HTTPS (本番環境)
- HTTP (開発環境)

### データ形式
- リクエスト: JSON または multipart/form-data（画像アップロード時）
- レスポンス: JSON

### 文字コード
UTF-8

## 認証

### 認証方式
JWT (JSON Web Token) を使用したBearer認証

### 認証ヘッダー
```
Authorization: Bearer <access_token>
```

### トークン有効期限
- アクセストークン: 24時間
- リフレッシュトークン: 30日間

### 認証不要エンドポイント
- POST /auth/signup
- POST /auth/login
- GET /health

## 共通仕様

### HTTPステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 204 | 成功（レスポンスボディなし） |
| 400 | リクエスト不正 |
| 401 | 認証エラー |
| 403 | 権限エラー |
| 404 | リソース未存在 |
| 409 | 競合エラー |
| 422 | バリデーションエラー |
| 500 | サーバーエラー |

### ページネーション

リスト系APIは以下のクエリパラメータでページネーションをサポート：

| パラメータ | 型 | デフォルト | 説明 |
|------------|-----|------------|------|
| page | integer | 1 | ページ番号 |
| limit | integer | 20 | 1ページあたりの件数（最大100） |

レスポンスヘッダー：
```
X-Total-Count: 総件数
X-Page: 現在のページ
X-Per-Page: 1ページあたりの件数
```

### 日時フォーマット
ISO 8601形式（RFC3339）
例: `2024-01-01T09:00:00Z`

## エラーハンドリング

### エラーレスポンス形式

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ",
    "details": {
      "field": "詳細情報"
    }
  }
}
```

### エラーコード一覧

| コード | 説明 |
|--------|------|
| INVALID_REQUEST | リクエストが不正 |
| VALIDATION_ERROR | バリデーションエラー |
| UNAUTHORIZED | 認証が必要 |
| FORBIDDEN | アクセス権限なし |
| NOT_FOUND | リソースが見つからない |
| CONFLICT | リソースの競合 |
| INTERNAL_ERROR | サーバー内部エラー |

## API詳細

### 認証API

#### ユーザー登録

```
POST /auth/signup
```

**リクエスト**
```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

**レスポンス**
```json
{
  "user": {
    "id": 1,
    "name": "string",
    "email": "string",
    "created_at": "2024-01-01T09:00:00Z",
    "updated_at": "2024-01-01T09:00:00Z"
  },
  "token": "string"
}
```

**バリデーション**
- name: 必須、1-255文字
- email: 必須、有効なメールアドレス、ユニーク
- password: 必須、8文字以上

#### ログイン

```
POST /auth/login
```

**リクエスト**
```json
{
  "email": "string",
  "password": "string"
}
```

**レスポンス**
```json
{
  "user": {
    "id": 1,
    "name": "string",
    "email": "string",
    "avatar": "string",
    "introduction": "string"
  },
  "token": "string"
}
```

#### ログアウト

```
POST /auth/logout
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Successfully logged out"
}
```

#### トークンリフレッシュ

```
POST /auth/refresh
```

**認証**: 必要

**レスポンス**
```json
{
  "token": "string"
}
```

#### 現在のユーザー情報取得

```
GET /auth/me
```

**認証**: 必要

**レスポンス**
```json
{
  "id": 1,
  "name": "string",
  "email": "string",
  "avatar": "string",
  "introduction": "string",
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-01T09:00:00Z"
}
```

### ユーザーAPI

#### ユーザー一覧取得

```
GET /users
```

**認証**: 必要

**クエリパラメータ**
- page: integer
- limit: integer

**レスポンス**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string",
      "email": "string",
      "avatar": "string",
      "introduction": "string",
      "followers_count": 10,
      "following_count": 5,
      "is_following": false,
      "is_blocked": false
    }
  ]
}
```

#### ユーザー詳細取得

```
GET /users/{id}
```

**認証**: 必要

**パスパラメータ**
- id: integer（ユーザーID）

**レスポンス**
```json
{
  "id": 1,
  "name": "string",
  "email": "string",
  "avatar": "string",
  "introduction": "string",
  "followers_count": 10,
  "following_count": 5,
  "items_count": 20,
  "coordinates_count": 15,
  "is_following": false,
  "is_blocked": false,
  "created_at": "2024-01-01T09:00:00Z"
}
```

#### プロフィール更新

```
PUT /users/{id}
```

**認証**: 必要（本人のみ）

**リクエスト**
```json
{
  "name": "string",
  "email": "string",
  "introduction": "string",
  "current_password": "string",
  "new_password": "string"
}
```

**レスポンス**
```json
{
  "id": 1,
  "name": "string",
  "email": "string",
  "avatar": "string",
  "introduction": "string",
  "updated_at": "2024-01-01T09:00:00Z"
}
```

#### アバター画像アップロード

```
POST /users/{id}/avatar
```

**認証**: 必要（本人のみ）

**リクエスト**
- Content-Type: multipart/form-data
- avatar: file（画像ファイル、最大5MB）

**レスポンス**
```json
{
  "avatar": "string"
}
```

### アイテムAPI

#### アイテム一覧取得

```
GET /items
```

**認証**: 必要

**クエリパラメータ**
- user_id: integer（ユーザーIDでフィルタ）
- season: integer（季節でフィルタ）
- tpo: integer（TPOでフィルタ）
- color: integer（色でフィルタ）
- super_item: string（カテゴリでフィルタ）
- min_rating: float（最低評価）
- max_rating: float（最高評価）
- page: integer
- limit: integer

**レスポンス**
```json
{
  "items": [
    {
      "id": 1,
      "user_id": 1,
      "name": "string",
      "image": "string",
      "super_item": "tops",
      "season": 1,
      "tpo": 1,
      "color": 1,
      "content": "string",
      "rating": 5,
      "created_at": "2024-01-01T09:00:00Z"
    }
  ]
}
```

#### アイテム詳細取得

```
GET /items/{id}
```

**認証**: 必要

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "user": {
    "id": 1,
    "name": "string",
    "avatar": "string"
  },
  "name": "string",
  "image": "string",
  "super_item": "tops",
  "season": 1,
  "tpo": 1,
  "color": 1,
  "content": "string",
  "rating": 5,
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-01T09:00:00Z"
}
```

#### アイテム作成

```
POST /items
```

**認証**: 必要

**リクエスト**
- Content-Type: multipart/form-data
- name: string（必須）
- image: file（画像ファイル、必須、最大5MB）
- super_item: string（必須、tops/bottoms/shoes/outer/accessory）
- season: integer（必須、1-4）
- tpo: integer（必須、1-5）
- color: integer（必須、1-13）
- content: string
- rating: integer（必須、1-5）

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "string",
  "image": "string",
  "super_item": "tops",
  "season": 1,
  "tpo": 1,
  "color": 1,
  "content": "string",
  "rating": 5,
  "created_at": "2024-01-01T09:00:00Z"
}
```

#### アイテム更新

```
PUT /items/{id}
```

**認証**: 必要（作成者のみ）

**リクエスト**
- Content-Type: multipart/form-data
- name: string
- image: file（画像ファイル、最大5MB）
- super_item: string
- season: integer
- tpo: integer
- color: integer
- content: string
- rating: integer

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "string",
  "image": "string",
  "super_item": "tops",
  "season": 1,
  "tpo": 1,
  "color": 1,
  "content": "string",
  "rating": 5,
  "updated_at": "2024-01-01T09:00:00Z"
}
```

#### アイテム削除

```
DELETE /items/{id}
```

**認証**: 必要（作成者のみ）

**レスポンス**
- ステータス: 204 No Content

### コーディネートAPI

#### コーディネート一覧取得

```
GET /coordinates
```

**認証**: 必要

**クエリパラメータ**
- user_id: integer
- season: integer
- tpo: integer
- min_rating: float
- max_rating: float
- page: integer
- limit: integer

**レスポンス**
```json
{
  "coordinates": [
    {
      "id": 1,
      "user_id": 1,
      "user": {
        "id": 1,
        "name": "string",
        "avatar": "string"
      },
      "image": "string",
      "season": 1,
      "tpo": 1,
      "rating": 5,
      "content": "string",
      "likes_count": 10,
      "comments_count": 5,
      "is_liked": false,
      "items": {
        "shoes": {
          "id": 1,
          "name": "string",
          "image": "string"
        },
        "bottoms": {
          "id": 2,
          "name": "string",
          "image": "string"
        },
        "tops": {
          "id": 3,
          "name": "string",
          "image": "string"
        }
      },
      "created_at": "2024-01-01T09:00:00Z"
    }
  ]
}
```

#### コーディネート詳細取得

```
GET /coordinates/{id}
```

**認証**: 必要

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "user": {
    "id": 1,
    "name": "string",
    "avatar": "string",
    "is_following": false
  },
  "image": "string",
  "season": 1,
  "tpo": 1,
  "rating": 5,
  "content": "string",
  "likes_count": 10,
  "comments_count": 5,
  "is_liked": false,
  "items": {
    "shoes": {
      "id": 1,
      "name": "string",
      "image": "string",
      "super_item": "shoes"
    },
    "bottoms": {
      "id": 2,
      "name": "string",
      "image": "string",
      "super_item": "bottoms"
    },
    "tops": {
      "id": 3,
      "name": "string",
      "image": "string",
      "super_item": "tops"
    },
    "outer": {
      "id": 4,
      "name": "string",
      "image": "string",
      "super_item": "outer"
    }
  },
  "created_at": "2024-01-01T09:00:00Z",
  "updated_at": "2024-01-01T09:00:00Z"
}
```

#### コーディネート作成

```
POST /coordinates
```

**認証**: 必要

**リクエスト**
```json
{
  "si_shoes": 1,
  "si_bottoms": 2,
  "si_tops": 3,
  "si_outer": 4,
  "season": 1,
  "tpo": 1,
  "rating": 5,
  "content": "string"
}
```
または
- Content-Type: multipart/form-data
- image: file（画像ファイル、最大5MB）
- 上記のJSONフィールドをform-dataとして送信

**バリデーション**
- 最低3つのアイテム（shoes, bottoms, tops）が必須
- 全アイテムは同一ユーザーが所有している必要がある

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "image": "string",
  "season": 1,
  "tpo": 1,
  "rating": 5,
  "content": "string",
  "created_at": "2024-01-01T09:00:00Z"
}
```

#### コーディネート更新

```
PUT /coordinates/{id}
```

**認証**: 必要（作成者のみ）

**リクエスト**
- 作成時と同じ形式

**レスポンス**
- 作成時と同じ形式

#### コーディネート削除

```
DELETE /coordinates/{id}
```

**認証**: 必要（作成者のみ）

**レスポンス**
- ステータス: 204 No Content

#### タイムライン取得

```
GET /coordinates/timeline
```

**認証**: 必要

**説明**: フォローしているユーザーのコーディネートを取得

**クエリパラメータ**
- page: integer
- limit: integer

**レスポンス**
- コーディネート一覧と同じ形式

#### いいね

```
POST /coordinates/{id}/like
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Liked successfully"
}
```

#### いいね取り消し

```
DELETE /coordinates/{id}/like
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Unliked successfully"
}
```

### ソーシャルAPI

#### コメント一覧取得

```
GET /comments
```

**認証**: 必要

**クエリパラメータ**
- coordinate_id: integer（必須）
- page: integer
- limit: integer

**レスポンス**
```json
{
  "comments": [
    {
      "id": 1,
      "user_id": 1,
      "user": {
        "id": 1,
        "name": "string",
        "avatar": "string"
      },
      "coordinate_id": 1,
      "comment": "string",
      "created_at": "2024-01-01T09:00:00Z",
      "updated_at": "2024-01-01T09:00:00Z"
    }
  ]
}
```

#### コメント投稿

```
POST /comments
```

**認証**: 必要

**リクエスト**
```json
{
  "coordinate_id": 1,
  "comment": "string"
}
```

**レスポンス**
```json
{
  "id": 1,
  "user_id": 1,
  "coordinate_id": 1,
  "comment": "string",
  "created_at": "2024-01-01T09:00:00Z"
}
```

#### コメント更新

```
PUT /comments/{id}
```

**認証**: 必要（投稿者のみ）

**リクエスト**
```json
{
  "comment": "string"
}
```

**レスポンス**
```json
{
  "id": 1,
  "comment": "string",
  "updated_at": "2024-01-01T09:00:00Z"
}
```

#### コメント削除

```
DELETE /comments/{id}
```

**認証**: 必要（投稿者のみ）

**レスポンス**
- ステータス: 204 No Content

#### フォロー

```
POST /users/{id}/follow
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Followed successfully"
}
```

#### フォロー解除

```
DELETE /users/{id}/follow
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Unfollowed successfully"
}
```

#### フォロワー一覧

```
GET /users/{id}/followers
```

**認証**: 必要

**クエリパラメータ**
- page: integer
- limit: integer

**レスポンス**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string",
      "avatar": "string",
      "is_following": true
    }
  ]
}
```

#### フォロー中一覧

```
GET /users/{id}/following
```

**認証**: 必要

**クエリパラメータ**
- page: integer
- limit: integer

**レスポンス**
- フォロワー一覧と同じ形式

#### ブロック

```
POST /users/{id}/block
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Blocked successfully"
}
```

#### ブロック解除

```
DELETE /users/{id}/block
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Unblocked successfully"
}
```

#### ブロックリスト

```
GET /blocks
```

**認証**: 必要

**レスポンス**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string",
      "avatar": "string"
    }
  ]
}
```

#### 通知一覧

```
GET /notifications
```

**認証**: 必要

**クエリパラメータ**
- unread_only: boolean（未読のみ）
- page: integer
- limit: integer

**レスポンス**
```json
{
  "notifications": [
    {
      "id": 1,
      "sender": {
        "id": 1,
        "name": "string",
        "avatar": "string"
      },
      "action": "follow",
      "coordinate": {
        "id": 1,
        "image": "string"
      },
      "comment": {
        "id": 1,
        "comment": "string"
      },
      "is_read": false,
      "created_at": "2024-01-01T09:00:00Z"
    }
  ],
  "unread_count": 5
}
```

#### 通知既読

```
PUT /notifications/{id}/read
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "Marked as read"
}
```

#### 全通知既読

```
PUT /notifications/read_all
```

**認証**: 必要

**レスポンス**
```json
{
  "message": "All notifications marked as read"
}
```

## 定数定義

### 季節 (Season)
| 値 | 意味 |
|----|------|
| 1 | 春 |
| 2 | 夏 |
| 3 | 秋 |
| 4 | 冬 |

### TPO (Time, Place, Occasion)
| 値 | 意味 |
|----|------|
| 1 | 仕事 |
| 2 | プライベート |
| 3 | スポーツ |
| 4 | デート |
| 5 | リラックス |

### 色 (Color)
| 値 | 意味 |
|----|------|
| 1 | 白 |
| 2 | 黒 |
| 3 | グレー |
| 4 | ベージュ |
| 5 | ブラウン |
| 6 | ネイビー |
| 7 | カーキ |
| 8 | 赤 |
| 9 | ピンク |
| 10 | オレンジ |
| 11 | イエロー |
| 12 | グリーン |
| 13 | ブルー |

### カテゴリ (SuperItem)
| 値 | 意味 |
|------|------|
| tops | トップス |
| bottoms | ボトムス |
| shoes | シューズ |
| outer | アウター |
| accessory | アクセサリー |

### 通知アクション (NotificationAction)
| 値 | 意味 |
|------|------|
| follow | フォロー |
| like | いいね |
| comment | コメント |