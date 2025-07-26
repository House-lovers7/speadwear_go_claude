# Speadwear テーブル定義書

## 概要
このドキュメントは、Speadwearアプリケーションのデータベーステーブル定義を記載しています。
データベースはMySQL 8.0を使用し、文字セットはUTF8MB4を採用しています。

## 共通カラム
全てのテーブルは以下の共通カラムを持ちます（BaseModel）：

| カラム名 | データ型 | NULL | デフォルト | 説明 |
|---------|---------|------|-----------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | 主キー |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP ON UPDATE | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | 論理削除日時（ソフトデリート） |

## 1. usersテーブル
ユーザー情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| name | VARCHAR(255) | NO | - | - | ユーザー名 |
| email | VARCHAR(255) | NO | - | UNIQUE | メールアドレス |
| picture | VARCHAR(255) | YES | NULL | - | プロフィール画像URL |
| admin | BOOLEAN | NO | FALSE | - | 管理者フラグ |
| password_digest | VARCHAR(255) | NO | - | - | パスワードハッシュ |
| remember_digest | VARCHAR(255) | YES | NULL | - | Remember meトークン |
| activation_digest | VARCHAR(255) | YES | NULL | - | アカウント有効化トークン |
| activated | BOOLEAN | NO | FALSE | - | アカウント有効化フラグ |
| activated_at | TIMESTAMP | YES | NULL | - | アカウント有効化日時 |
| reset_digest | VARCHAR(255) | YES | NULL | - | パスワードリセットトークン |
| reset_sent_at | TIMESTAMP | YES | NULL | - | パスワードリセットメール送信日時 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### リレーション
- items (1:N) - ユーザーが所有するアイテム
- coordinates (1:N) - ユーザーが作成したコーディネート
- comments (1:N) - ユーザーが投稿したコメント
- followers (N:N via relationships) - フォロワー
- following (N:N via relationships) - フォロー中のユーザー
- blocks_as_blocker (1:N) - ブロックしたユーザー
- blocks_as_blocked (1:N) - ブロックされたユーザー
- sent_notifications (1:N) - 送信した通知
- received_notifications (1:N) - 受信した通知

## 2. itemsテーブル
服アイテム情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| user_id | BIGINT UNSIGNED | NO | - | INDEX | ユーザーID（外部キー） |
| coordinate_id | BIGINT UNSIGNED | YES | NULL | - | コーディネートID（外部キー） |
| super_item | VARCHAR(100) | YES | NULL | - | カテゴリー（tops/bottoms/shoes/outer/accessory） |
| season | INT | YES | NULL | - | 季節（1:春, 2:夏, 3:秋, 4:冬, 5:オールシーズン） |
| tpo | INT | YES | NULL | - | TPO（1:仕事, 2:カジュアル, 3:フォーマル, 4:スポーツ, 5:ホーム） |
| color | INT | YES | NULL | - | 色（1-15: 各色定義） |
| content | TEXT | YES | NULL | - | 説明文 |
| memo | TEXT | YES | NULL | - | メモ |
| picture | VARCHAR(255) | YES | NULL | - | アイテム画像URL |
| rating | FLOAT | YES | NULL | - | 評価（0.0-5.0） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- user_id → users.id (CASCADE DELETE)
- coordinate_id → coordinates.id (SET NULL)

## 3. coordinatesテーブル
コーディネート情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| user_id | BIGINT UNSIGNED | NO | - | INDEX | ユーザーID（外部キー） |
| season | INT | YES | NULL | - | 季節（1-5） |
| tpo | INT | YES | NULL | - | TPO（1-5） |
| picture | VARCHAR(255) | YES | NULL | - | コーディネート画像URL |
| si_top_length | INT | YES | NULL | - | トップス丈（1:クロップ, 2:ノーマル, 3:ロング） |
| si_top_sleeve | INT | YES | NULL | - | トップス袖（1:ノースリーブ, 2:キャップ, 3:半袖, 4:七分袖, 5:長袖） |
| si_bottom_length | INT | YES | NULL | - | ボトムス丈（1-6: ミニ〜マキシ） |
| si_bottom_type | INT | YES | NULL | - | ボトムスタイプ（1:スカート, 2:パンツ） |
| si_dress_length | INT | YES | NULL | - | ワンピース丈（1-6） |
| si_dress_sleeve | INT | YES | NULL | - | ワンピース袖（1-5） |
| si_outer_length | INT | YES | NULL | - | アウター丈（1:ショート, 2:ノーマル, 3:ロング） |
| si_outer_sleeve | INT | YES | NULL | - | アウター袖（1:半袖, 2:七分袖, 3:長袖） |
| si_shoe_size | INT | YES | NULL | - | シューズサイズ |
| memo | TEXT | YES | NULL | - | メモ |
| rating | FLOAT | YES | NULL | - | 評価（0.0-5.0） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- user_id → users.id (CASCADE DELETE)

### リレーション
- items (1:N) - コーディネートに含まれるアイテム
- comments (1:N) - コーディネートへのコメント
- like_coordinates (1:N) - いいね

## 4. commentsテーブル
コメント情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| user_id | BIGINT UNSIGNED | NO | - | INDEX | ユーザーID（外部キー） |
| coordinate_id | BIGINT UNSIGNED | NO | - | INDEX | コーディネートID（外部キー） |
| comment | TEXT | NO | - | - | コメント内容 |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- user_id → users.id (CASCADE DELETE)
- coordinate_id → coordinates.id (CASCADE DELETE)

## 5. like_coordinatesテーブル
コーディネートへのいいねを管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| user_id | BIGINT UNSIGNED | NO | - | INDEX | ユーザーID（外部キー） |
| coordinate_id | BIGINT UNSIGNED | NO | - | INDEX | コーディネートID（外部キー） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- user_id → users.id (CASCADE DELETE)
- coordinate_id → coordinates.id (CASCADE DELETE)

### 複合ユニークインデックス
- (user_id, coordinate_id) - 同一ユーザーが同じコーディネートに複数回いいねできない

## 6. relationshipsテーブル
フォロー関係を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| follower_id | BIGINT UNSIGNED | NO | - | INDEX | フォロワーID（外部キー） |
| followed_id | BIGINT UNSIGNED | NO | - | INDEX | フォロー対象ID（外部キー） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- follower_id → users.id (CASCADE DELETE)
- followed_id → users.id (CASCADE DELETE)

### 複合ユニークインデックス
- (follower_id, followed_id) - 同一の関係を重複して作成できない

## 7. blocksテーブル
ブロック関係を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| blocker_id | BIGINT UNSIGNED | NO | - | INDEX | ブロックしたユーザーID（外部キー） |
| blocked_id | BIGINT UNSIGNED | NO | - | INDEX | ブロックされたユーザーID（外部キー） |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- blocker_id → users.id (CASCADE DELETE)
- blocked_id → users.id (CASCADE DELETE)

### 複合ユニークインデックス
- (blocker_id, blocked_id) - 同一の関係を重複して作成できない

## 8. notificationsテーブル
通知情報を管理するテーブル

| カラム名 | データ型 | NULL | デフォルト | インデックス | 説明 |
|---------|---------|------|-----------|-------------|------|
| id | BIGINT UNSIGNED | NO | AUTO_INCREMENT | PRIMARY | 主キー |
| sender_id | BIGINT UNSIGNED | NO | - | INDEX | 送信者ID（外部キー） |
| receiver_id | BIGINT UNSIGNED | NO | - | INDEX | 受信者ID（外部キー） |
| coordinate_id | BIGINT UNSIGNED | YES | NULL | - | 関連コーディネートID（外部キー） |
| comment_id | BIGINT UNSIGNED | YES | NULL | - | 関連コメントID（外部キー） |
| like_coordinate_id | BIGINT UNSIGNED | YES | NULL | - | 関連いいねID（外部キー） |
| action | VARCHAR(50) | NO | - | - | アクション種別（follow/like/comment） |
| checked | BOOLEAN | NO | FALSE | - | 既読フラグ |
| created_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 作成日時 |
| updated_at | TIMESTAMP | NO | CURRENT_TIMESTAMP | - | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | INDEX | 削除日時 |

### 外部キー制約
- sender_id → users.id (CASCADE DELETE)
- receiver_id → users.id (CASCADE DELETE)
- coordinate_id → coordinates.id (SET NULL)
- comment_id → comments.id (SET NULL)
- like_coordinate_id → like_coordinates.id (SET NULL)

## インデックス設計

### パフォーマンス最適化のための追加インデックス
1. items テーブル
   - (user_id, super_item) - ユーザーのカテゴリ別アイテム検索
   - (season, tpo, color) - フィルタリング検索

2. coordinates テーブル
   - (user_id, season, tpo) - ユーザーの条件別コーディネート検索
   - (created_at DESC) - タイムライン表示

3. notifications テーブル
   - (receiver_id, checked, created_at DESC) - 未読通知の取得

## データベース設定
- 文字セット: utf8mb4
- 照合順序: utf8mb4_unicode_ci
- ストレージエンジン: InnoDB
- タイムゾーン: UTC

## マイグレーション
GORMの自動マイグレーション機能を使用して、モデル定義からテーブルを作成・更新します。
詳細は `/internal/domain/models.go` を参照してください。