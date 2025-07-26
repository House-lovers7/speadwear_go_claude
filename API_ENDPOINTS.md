# Speadwear API エンドポイント一覧

## Base URL
```
http://localhost:8080/api/v1
```

## 認証
ほとんどのエンドポイントは認証が必要です。Authorizationヘッダーに`Bearer <token>`形式でJWTトークンを含める必要があります。

## エンドポイント

### 認証 (Authentication)

#### ユーザー登録
```
POST /auth/signup
Content-Type: application/json

{
  "name": "ユーザー名",
  "email": "user@example.com",
  "password": "password123"
}
```

#### ログイン
```
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### トークンリフレッシュ
```
POST /auth/refresh
Authorization: Bearer <token>
```

#### ログアウト
```
POST /auth/logout
Authorization: Bearer <token>
```

#### 現在のユーザー情報取得
```
GET /auth/me
Authorization: Bearer <token>
```

### ユーザー管理 (Users)

#### ユーザー一覧取得
```
GET /users?page=1&per_page=20
Authorization: Bearer <token>
```

#### ユーザー詳細取得
```
GET /users/:id
```

#### プロフィール更新
```
PUT /users/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "新しい名前",
  "picture": "画像URL"
}
```

#### パスワード変更
```
PUT /users/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "現在のパスワード",
  "new_password": "新しいパスワード"
}
```

### アイテム管理 (Items)

#### アイテム作成
```
POST /items
Authorization: Bearer <token>
Content-Type: multipart/form-data

super_item: "カテゴリー"
season: 1-5 (1:春, 2:夏, 3:秋, 4:冬, 5:オールシーズン)
tpo: 1-5 (1:仕事, 2:カジュアル, 3:フォーマル, 4:スポーツ, 5:ホーム)
color: 1-15 (色ID)
content: "説明"
memo: "メモ"
rating: 0-5
picture: (画像ファイル)
```

#### 自分のアイテム一覧取得
```
GET /items?page=1&per_page=20
Authorization: Bearer <token>
```

#### 特定ユーザーのアイテム一覧取得
```
GET /users/:user_id/items?page=1&per_page=20
```

#### アイテム詳細取得
```
GET /items/:id
```

#### アイテム検索
```
GET /items/search?season=1&tpo=2&color=3&super_item=トップス&min_rating=3&max_rating=5
```

#### アイテム更新
```
PUT /items/:id
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

#### アイテム削除
```
DELETE /items/:id
Authorization: Bearer <token>
```

#### アイテム一括削除
```
DELETE /items
Authorization: Bearer <token>
Content-Type: application/json

{
  "item_ids": [1, 2, 3]
}
```

#### アイテム統計情報取得
```
GET /items/statistics
Authorization: Bearer <token>
```

### コーディネート管理 (Coordinates)

#### コーディネート作成
```
POST /coordinates
Authorization: Bearer <token>
Content-Type: multipart/form-data

season: 1-5
tpo: 1-5
item_ids: [1, 2, 3]
memo: "メモ"
rating: 0-5
picture: (画像ファイル)
si_top_length: 1-3
si_top_sleeve: 1-5
si_bottom_length: 1-6
si_bottom_type: 1-2
si_dress_length: 1-6
si_dress_sleeve: 1-5
si_outer_length: 1-3
si_outer_sleeve: 1-3
si_shoe_size: サイズ
```

#### 自分のコーディネート一覧取得
```
GET /coordinates?page=1&per_page=20
Authorization: Bearer <token>
```

#### タイムライン取得（フォローしているユーザーのコーディネート）
```
GET /coordinates/timeline?page=1&per_page=20
Authorization: Bearer <token>
```

#### コーディネート詳細取得
```
GET /coordinates/:id
```

#### コーディネート検索
```
GET /coordinates/search?season=1&tpo=2&min_rating=3&max_rating=5
```

#### コーディネート更新
```
PUT /coordinates/:id
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

#### コーディネート削除
```
DELETE /coordinates/:id
Authorization: Bearer <token>
```

#### コーディネート統計情報取得
```
GET /coordinates/statistics
Authorization: Bearer <token>
```

### いいね機能 (Likes)

#### いいねする
```
POST /coordinates/:id/like
Authorization: Bearer <token>
```

#### いいね解除
```
DELETE /coordinates/:id/like
Authorization: Bearer <token>
```

### コメント機能 (Comments)

#### コメント投稿
```
POST /comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "coordinate_id": 1,
  "comment": "コメント内容"
}
```

#### コーディネートのコメント一覧取得
```
GET /coordinates/:id/comments?page=1&per_page=20
```

#### コメント更新
```
PUT /comments/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "comment": "更新後のコメント"
}
```

#### コメント削除
```
DELETE /comments/:id
Authorization: Bearer <token>
```

### フォロー機能 (Follow)

#### フォローする
```
POST /follow/:user_id
Authorization: Bearer <token>
```

#### フォロー解除
```
DELETE /follow/:user_id
Authorization: Bearer <token>
```

#### フォロワー一覧取得
```
GET /follow/followers?page=1&per_page=20
Authorization: Bearer <token>
```

#### フォロー中一覧取得
```
GET /follow/following?page=1&per_page=20
Authorization: Bearer <token>
```

#### フォロー状態確認
```
GET /follow/status/:user_id
Authorization: Bearer <token>
```

### ブロック機能 (Block)

#### ブロックする
```
POST /blocks/:user_id
Authorization: Bearer <token>
```

#### ブロック解除
```
DELETE /blocks/:user_id
Authorization: Bearer <token>
```

#### ブロック中のユーザー一覧取得
```
GET /blocks
Authorization: Bearer <token>
```

#### ブロック状態確認
```
GET /blocks/status/:user_id
Authorization: Bearer <token>
```

### 通知機能 (Notifications)

#### 通知一覧取得
```
GET /notifications?page=1&per_page=20
Authorization: Bearer <token>
```

#### 未読通知取得
```
GET /notifications/unread
Authorization: Bearer <token>
```

#### 未読通知数取得
```
GET /notifications/unread/count
Authorization: Bearer <token>
```

#### 通知を既読にする
```
PUT /notifications/:id/read
Authorization: Bearer <token>
```

#### 全通知を既読にする
```
PUT /notifications/read_all
Authorization: Bearer <token>
```

## レスポンス形式

### 成功レスポンス
```json
{
  "data": {
    // レスポンスデータ
  }
}
```

### エラーレスポンス
```json
{
  "error": "エラーメッセージ"
}
```

### ページネーションレスポンス
```json
{
  "items": [...],
  "total_count": 100,
  "page": 1,
  "per_page": 20
}
```