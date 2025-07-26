# Speadwear 画面遷移図

## 概要

このドキュメントでは、Speadwearアプリケーションの画面遷移をMermaid記法で図示します。

## 全体の画面構成

```mermaid
graph TB
    Start([アプリ起動]) --> Check{ログイン済み?}
    Check -->|Yes| Home[ホーム画面]
    Check -->|No| Landing[ランディングページ]
    
    Landing --> Login[ログイン画面]
    Landing --> Signup[新規登録画面]
    
    Login --> Home
    Signup --> Tutorial[チュートリアル]
    Tutorial --> Home
    
    Home --> Timeline[タイムライン]
    Home --> MyPage[マイページ]
    Home --> Search[検索]
    Home --> Create[作成]
    
    subgraph "メインナビゲーション"
        Timeline
        Search
        Create
        MyPage
    end
```

## 認証フロー

```mermaid
graph TD
    Start([アプリ起動]) --> CheckAuth{認証状態確認}
    CheckAuth -->|未認証| Landing[ランディングページ]
    CheckAuth -->|認証済| Home[ホーム画面]
    
    Landing --> LoginChoice{選択}
    LoginChoice --> Login[ログイン画面]
    LoginChoice --> Signup[新規登録画面]
    
    Login --> LoginForm[ログインフォーム]
    LoginForm --> LoginSubmit{送信}
    LoginSubmit -->|成功| Home
    LoginSubmit -->|失敗| LoginError[エラー表示]
    LoginError --> LoginForm
    
    Signup --> SignupForm[登録フォーム]
    SignupForm --> SignupSubmit{送信}
    SignupSubmit -->|成功| EmailVerify[メール確認画面]
    SignupSubmit -->|失敗| SignupError[エラー表示]
    SignupError --> SignupForm
    EmailVerify --> Tutorial[チュートリアル]
    Tutorial --> Home
    
    Login --> ForgotPassword[パスワードを忘れた]
    ForgotPassword --> ResetEmail[メールアドレス入力]
    ResetEmail --> ResetSent[送信完了画面]
```

## アイテム管理フロー

```mermaid
graph TD
    MyPage[マイページ] --> ItemList[アイテム一覧]
    ItemList --> ItemDetail[アイテム詳細]
    ItemList --> ItemCreate[アイテム作成]
    
    ItemCreate --> ItemForm[アイテム入力フォーム]
    ItemForm --> PhotoSelect[写真選択]
    PhotoSelect --> Camera[カメラ撮影]
    PhotoSelect --> Gallery[ギャラリー選択]
    
    Camera --> PhotoEdit[写真編集]
    Gallery --> PhotoEdit
    PhotoEdit --> ItemForm
    
    ItemForm --> CategorySelect[カテゴリ選択]
    CategorySelect --> SeasonSelect[季節選択]
    SeasonSelect --> TPOSelect[TPO選択]
    TPOSelect --> ColorSelect[色選択]
    ColorSelect --> RatingSelect[評価選択]
    RatingSelect --> ItemConfirm[確認画面]
    
    ItemConfirm --> ItemSubmit{登録}
    ItemSubmit -->|成功| ItemComplete[完了画面]
    ItemSubmit -->|失敗| ItemError[エラー表示]
    ItemComplete --> ItemDetail
    ItemError --> ItemForm
    
    ItemDetail --> ItemEdit[編集]
    ItemDetail --> ItemDelete{削除確認}
    ItemEdit --> ItemForm
    ItemDelete -->|Yes| ItemList
    ItemDelete -->|No| ItemDetail
```

## コーディネート作成フロー

```mermaid
graph TD
    Create[作成ボタン] --> CreateChoice{作成選択}
    CreateChoice --> ItemCreate[アイテム作成]
    CreateChoice --> CoordCreate[コーディネート作成]
    
    CoordCreate --> ItemSelect[アイテム選択画面]
    ItemSelect --> SelectShoes[シューズ選択]
    SelectShoes --> SelectBottoms[ボトムス選択]
    SelectBottoms --> SelectTops[トップス選択]
    SelectTops --> SelectOuter[アウター選択/スキップ]
    
    SelectOuter --> CoordPreview[プレビュー画面]
    CoordPreview --> CoordEdit[編集画面]
    
    CoordEdit --> SeasonSelect[季節選択]
    SeasonSelect --> TPOSelect[TPO選択]
    TPOSelect --> RatingSelect[評価選択]
    RatingSelect --> CommentInput[コメント入力]
    CommentInput --> CoordConfirm[確認画面]
    
    CoordConfirm --> CoordSubmit{投稿}
    CoordSubmit -->|成功| CoordComplete[完了画面]
    CoordSubmit -->|失敗| CoordError[エラー表示]
    CoordComplete --> CoordDetail[コーディネート詳細]
    CoordError --> CoordEdit
    
    CoordPreview --> PhotoOption{写真オプション}
    PhotoOption --> TakePhoto[写真撮影]
    PhotoOption --> UploadPhoto[写真アップロード]
    TakePhoto --> CoordEdit
    UploadPhoto --> CoordEdit
```

## タイムライン・検索フロー

```mermaid
graph TD
    Timeline[タイムライン] --> TimelineList[コーディネート一覧]
    TimelineList --> CoordDetail[コーディネート詳細]
    TimelineList --> LoadMore[さらに読み込む]
    
    Search[検索] --> SearchFilter[フィルター設定]
    SearchFilter --> SeasonFilter[季節フィルター]
    SearchFilter --> TPOFilter[TPOフィルター]
    SearchFilter --> ColorFilter[色フィルター]
    SearchFilter --> RatingFilter[評価フィルター]
    
    SearchFilter --> SearchResult[検索結果]
    SearchResult --> CoordDetail
    SearchResult --> UserProfile[ユーザープロフィール]
    
    CoordDetail --> LikeAction{いいね}
    CoordDetail --> CommentList[コメント一覧]
    CoordDetail --> ShareMenu[共有メニュー]
    
    CommentList --> CommentInput[コメント入力]
    CommentInput --> CommentSubmit{送信}
    CommentSubmit --> CommentList
    
    CoordDetail --> ItemLink[アイテム詳細へ]
    ItemLink --> ItemDetail[アイテム詳細]
```

## ソーシャル機能フロー

```mermaid
graph TD
    UserProfile[ユーザープロフィール] --> ProfileInfo[プロフィール情報]
    ProfileInfo --> FollowButton{フォローボタン}
    FollowButton -->|未フォロー| Follow[フォローする]
    FollowButton -->|フォロー中| Unfollow[フォロー解除]
    
    UserProfile --> UserItems[ユーザーのアイテム]
    UserProfile --> UserCoords[ユーザーのコーディネート]
    UserProfile --> UserFollowers[フォロワー一覧]
    UserProfile --> UserFollowing[フォロー中一覧]
    
    UserProfile --> MenuButton[メニューボタン]
    MenuButton --> BlockUser[ブロックする]
    MenuButton --> ReportUser[通報する]
    
    MyPage[マイページ] --> EditProfile[プロフィール編集]
    EditProfile --> NameEdit[名前編集]
    EditProfile --> IntroEdit[自己紹介編集]
    EditProfile --> AvatarEdit[アバター編集]
    
    AvatarEdit --> AvatarSelect{選択}
    AvatarSelect --> Camera[カメラ]
    AvatarSelect --> Gallery[ギャラリー]
    
    MyPage --> Settings[設定]
    Settings --> AccountSettings[アカウント設定]
    Settings --> NotificationSettings[通知設定]
    Settings --> PrivacySettings[プライバシー設定]
    Settings --> BlockList[ブロックリスト]
```

## 通知フロー

```mermaid
graph TD
    App[アプリ] --> NotificationBadge[通知バッジ]
    NotificationBadge --> NotificationList[通知一覧]
    
    NotificationList --> NotificationItem{通知タイプ}
    NotificationItem -->|フォロー| FollowNotif[フォロー通知]
    NotificationItem -->|いいね| LikeNotif[いいね通知]
    NotificationItem -->|コメント| CommentNotif[コメント通知]
    
    FollowNotif --> UserProfile[ユーザープロフィール]
    LikeNotif --> CoordDetail[コーディネート詳細]
    CommentNotif --> CoordDetail
    
    NotificationList --> MarkAsRead[既読にする]
    NotificationList --> MarkAllAsRead[すべて既読にする]
```

## エラー・例外処理フロー

```mermaid
graph TD
    Action[ユーザーアクション] --> Process{処理}
    Process -->|成功| Success[成功画面/遷移]
    Process -->|失敗| ErrorType{エラータイプ}
    
    ErrorType -->|ネットワーク| NetworkError[ネットワークエラー]
    ErrorType -->|認証| AuthError[認証エラー]
    ErrorType -->|バリデーション| ValidationError[入力エラー]
    ErrorType -->|サーバー| ServerError[サーバーエラー]
    
    NetworkError --> Retry[再試行ボタン]
    AuthError --> Login[ログイン画面へ]
    ValidationError --> FormCorrection[フォーム修正]
    ServerError --> ErrorMessage[エラーメッセージ表示]
    
    Retry --> Action
    FormCorrection --> Action
```

## 画面一覧

### 認証関連
1. ランディングページ
2. ログイン画面
3. 新規登録画面
4. パスワードリセット画面
5. チュートリアル画面

### メイン画面
1. ホーム画面（タイムライン）
2. 検索画面
3. 作成選択画面
4. マイページ

### アイテム関連
1. アイテム一覧
2. アイテム詳細
3. アイテム作成/編集
4. カテゴリ選択
5. 写真撮影/選択

### コーディネート関連
1. コーディネート一覧
2. コーディネート詳細
3. コーディネート作成
4. アイテム選択（各カテゴリ）
5. コメント一覧

### ユーザー関連
1. ユーザープロフィール
2. プロフィール編集
3. フォロワー/フォロー中一覧
4. 設定画面
5. ブロックリスト

### その他
1. 通知一覧
2. エラー画面
3. 読み込み中画面
4. 空状態画面

## 主要なユーザーストーリー

### 1. 新規ユーザーの初回利用
```
ランディング → 新規登録 → チュートリアル → アイテム登録 → コーディネート作成 → タイムライン
```

### 2. 日常的な利用
```
ホーム → タイムライン閲覧 → いいね/コメント → フォロー → 通知確認
```

### 3. コーディネート投稿
```
作成ボタン → アイテム選択 → 詳細設定 → 投稿 → シェア
```

### 4. アイテム管理
```
マイページ → アイテム一覧 → アイテム追加 → カテゴリ設定 → 完了
```

### 5. 他ユーザーとの交流
```
検索 → ユーザー発見 → プロフィール閲覧 → フォロー → コメント投稿
```