# Speadwear API定義書（OpenAPI 3.0）

## 概要
このドキュメントは、Speadwear APIの詳細仕様をOpenAPI 3.0形式で定義しています。

```yaml
openapi: 3.0.3
info:
  title: Speadwear API
  description: ファッションコーディネートアプリケーションのRESTful API
  version: 1.0.0
  contact:
    name: Speadwear Development Team
    email: dev@speadwear.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.speadwear.com/api/v1
    description: Production server
  - url: http://localhost:8080/api/v1
    description: Development server

tags:
  - name: auth
    description: 認証関連のエンドポイント
  - name: users
    description: ユーザー管理
  - name: items
    description: アイテム管理
  - name: coordinates
    description: コーディネート管理
  - name: social
    description: ソーシャル機能（フォロー、いいね、コメント）
  - name: notifications
    description: 通知管理

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT認証トークン

  schemas:
    # 共通スキーマ
    Error:
      type: object
      required:
        - error
      properties:
        error:
          type: object
          required:
            - code
            - message
          properties:
            code:
              type: string
              example: VALIDATION_ERROR
            message:
              type: string
              example: バリデーションエラーが発生しました
            details:
              type: object
              additionalProperties: true

    Pagination:
      type: object
      properties:
        total_count:
          type: integer
          example: 100
        page:
          type: integer
          example: 1
        per_page:
          type: integer
          example: 20

    # ユーザー関連スキーマ
    User:
      type: object
      properties:
        id:
          type: integer
          format: int64
          example: 1
        name:
          type: string
          example: 山田太郎
        email:
          type: string
          format: email
          example: user@example.com
        picture:
          type: string
          format: uri
          example: https://api.speadwear.com/uploads/avatars/1.jpg
        introduction:
          type: string
          example: ファッションが大好きです
        admin:
          type: boolean
          example: false
        created_at:
          type: string
          format: date-time
          example: 2024-01-01T09:00:00Z
        updated_at:
          type: string
          format: date-time
          example: 2024-01-01T09:00:00Z

    UserWithStats:
      allOf:
        - $ref: '#/components/schemas/User'
        - type: object
          properties:
            followers_count:
              type: integer
              example: 100
            following_count:
              type: integer
              example: 50
            items_count:
              type: integer
              example: 30
            coordinates_count:
              type: integer
              example: 20
            is_following:
              type: boolean
              example: false
            is_blocked:
              type: boolean
              example: false

    # アイテム関連スキーマ
    Item:
      type: object
      properties:
        id:
          type: integer
          format: int64
          example: 1
        user_id:
          type: integer
          format: int64
          example: 1
        coordinate_id:
          type: integer
          format: int64
          nullable: true
          example: null
        super_item:
          type: string
          enum: [tops, bottoms, shoes, outer, accessory]
          example: tops
        season:
          type: integer
          minimum: 1
          maximum: 5
          example: 1
          description: 1:春, 2:夏, 3:秋, 4:冬, 5:オールシーズン
        tpo:
          type: integer
          minimum: 1
          maximum: 5
          example: 2
          description: 1:仕事, 2:カジュアル, 3:フォーマル, 4:スポーツ, 5:ホーム
        color:
          type: integer
          minimum: 1
          maximum: 15
          example: 7
          description: 色コード（1-15）
        content:
          type: string
          example: お気に入りのTシャツです
        memo:
          type: string
          example: 洗濯時は裏返しで
        picture:
          type: string
          format: uri
          example: https://api.speadwear.com/uploads/items/1.jpg
        rating:
          type: number
          format: float
          minimum: 0
          maximum: 5
          example: 4.5
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    ItemWithUser:
      allOf:
        - $ref: '#/components/schemas/Item'
        - type: object
          properties:
            user:
              $ref: '#/components/schemas/User'

    # コーディネート関連スキーマ
    Coordinate:
      type: object
      properties:
        id:
          type: integer
          format: int64
        user_id:
          type: integer
          format: int64
        season:
          type: integer
          minimum: 1
          maximum: 5
        tpo:
          type: integer
          minimum: 1
          maximum: 5
        picture:
          type: string
          format: uri
        si_top_length:
          type: integer
          minimum: 1
          maximum: 3
        si_top_sleeve:
          type: integer
          minimum: 1
          maximum: 5
        si_bottom_length:
          type: integer
          minimum: 1
          maximum: 6
        si_bottom_type:
          type: integer
          minimum: 1
          maximum: 2
        si_dress_length:
          type: integer
          minimum: 1
          maximum: 6
        si_dress_sleeve:
          type: integer
          minimum: 1
          maximum: 5
        si_outer_length:
          type: integer
          minimum: 1
          maximum: 3
        si_outer_sleeve:
          type: integer
          minimum: 1
          maximum: 3
        si_shoe_size:
          type: integer
        memo:
          type: string
        rating:
          type: number
          format: float
          minimum: 0
          maximum: 5
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    CoordinateWithDetails:
      allOf:
        - $ref: '#/components/schemas/Coordinate'
        - type: object
          properties:
            user:
              $ref: '#/components/schemas/User'
            items:
              type: array
              items:
                $ref: '#/components/schemas/Item'
            like_count:
              type: integer
              example: 10
            comment_count:
              type: integer
              example: 5
            is_liked:
              type: boolean
              example: false

    # コメント関連スキーマ
    Comment:
      type: object
      properties:
        id:
          type: integer
          format: int64
        user_id:
          type: integer
          format: int64
        coordinate_id:
          type: integer
          format: int64
        comment:
          type: string
        created_at:
          type: string
          format: date-time
        updated_at:
          type: string
          format: date-time

    CommentWithUser:
      allOf:
        - $ref: '#/components/schemas/Comment'
        - type: object
          properties:
            user:
              $ref: '#/components/schemas/User'

    # 通知関連スキーマ
    Notification:
      type: object
      properties:
        id:
          type: integer
          format: int64
        sender_id:
          type: integer
          format: int64
        receiver_id:
          type: integer
          format: int64
        coordinate_id:
          type: integer
          format: int64
          nullable: true
        comment_id:
          type: integer
          format: int64
          nullable: true
        like_coordinate_id:
          type: integer
          format: int64
          nullable: true
        action:
          type: string
          enum: [follow, like, comment]
        checked:
          type: boolean
        created_at:
          type: string
          format: date-time

    NotificationWithDetails:
      allOf:
        - $ref: '#/components/schemas/Notification'
        - type: object
          properties:
            sender:
              $ref: '#/components/schemas/User'
            coordinate:
              $ref: '#/components/schemas/Coordinate'
            comment:
              $ref: '#/components/schemas/Comment'

  responses:
    UnauthorizedError:
      description: 認証エラー
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error:
              code: UNAUTHORIZED
              message: 認証が必要です

    ForbiddenError:
      description: 権限エラー
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error:
              code: FORBIDDEN
              message: この操作を実行する権限がありません

    NotFoundError:
      description: リソースが見つかりません
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error:
              code: NOT_FOUND
              message: 指定されたリソースが見つかりません

    ValidationError:
      description: バリデーションエラー
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
          example:
            error:
              code: VALIDATION_ERROR
              message: 入力値が不正です
              details:
                email: メールアドレスの形式が正しくありません

security:
  - bearerAuth: []

paths:
  # 認証エンドポイント
  /auth/signup:
    post:
      tags:
        - auth
      summary: ユーザー登録
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - name
                - email
                - password
              properties:
                name:
                  type: string
                  minLength: 1
                  maxLength: 255
                  example: 山田太郎
                email:
                  type: string
                  format: email
                  example: user@example.com
                password:
                  type: string
                  format: password
                  minLength: 8
                  example: password123
      responses:
        '201':
          description: 登録成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  user:
                    $ref: '#/components/schemas/User'
                  token:
                    type: string
                    example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
        '422':
          $ref: '#/components/responses/ValidationError'

  /auth/login:
    post:
      tags:
        - auth
      summary: ログイン
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - email
                - password
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
                  format: password
      responses:
        '200':
          description: ログイン成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  user:
                    $ref: '#/components/schemas/User'
                  token:
                    type: string
        '401':
          $ref: '#/components/responses/UnauthorizedError'

  /auth/logout:
    post:
      tags:
        - auth
      summary: ログアウト
      responses:
        '200':
          description: ログアウト成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Successfully logged out

  /auth/refresh:
    post:
      tags:
        - auth
      summary: トークンリフレッシュ
      responses:
        '200':
          description: リフレッシュ成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string

  /auth/me:
    get:
      tags:
        - auth
      summary: 現在のユーザー情報取得
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserWithStats'

  # ユーザー管理エンドポイント
  /users:
    get:
      tags:
        - users
      summary: ユーザー一覧取得
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          headers:
            X-Total-Count:
              schema:
                type: integer
              description: 総件数
            X-Page:
              schema:
                type: integer
              description: 現在のページ
            X-Per-Page:
              schema:
                type: integer
              description: 1ページあたりの件数
          content:
            application/json:
              schema:
                type: object
                properties:
                  users:
                    type: array
                    items:
                      $ref: '#/components/schemas/UserWithStats'

  /users/{id}:
    get:
      tags:
        - users
      summary: ユーザー詳細取得
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserWithStats'
        '404':
          $ref: '#/components/responses/NotFoundError'

    put:
      tags:
        - users
      summary: プロフィール更新
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                email:
                  type: string
                  format: email
                introduction:
                  type: string
                current_password:
                  type: string
                  format: password
                new_password:
                  type: string
                  format: password
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '422':
          $ref: '#/components/responses/ValidationError'

  /users/{id}/avatar:
    post:
      tags:
        - users
      summary: アバター画像アップロード
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              required:
                - avatar
              properties:
                avatar:
                  type: string
                  format: binary
      responses:
        '200':
          description: アップロード成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  avatar:
                    type: string
                    format: uri

  # アイテム管理エンドポイント
  /items:
    get:
      tags:
        - items
      summary: アイテム一覧取得
      parameters:
        - name: user_id
          in: query
          schema:
            type: integer
            format: int64
        - name: season
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 5
        - name: tpo
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 5
        - name: color
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 15
        - name: super_item
          in: query
          schema:
            type: string
            enum: [tops, bottoms, shoes, outer, accessory]
        - name: min_rating
          in: query
          schema:
            type: number
            format: float
            minimum: 0
            maximum: 5
        - name: max_rating
          in: query
          schema:
            type: number
            format: float
            minimum: 0
            maximum: 5
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          headers:
            X-Total-Count:
              schema:
                type: integer
            X-Page:
              schema:
                type: integer
            X-Per-Page:
              schema:
                type: integer
          content:
            application/json:
              schema:
                type: object
                properties:
                  items:
                    type: array
                    items:
                      $ref: '#/components/schemas/ItemWithUser'

    post:
      tags:
        - items
      summary: アイテム作成
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              required:
                - name
                - image
                - super_item
                - season
                - tpo
                - color
                - rating
              properties:
                name:
                  type: string
                image:
                  type: string
                  format: binary
                super_item:
                  type: string
                  enum: [tops, bottoms, shoes, outer, accessory]
                season:
                  type: integer
                  minimum: 1
                  maximum: 5
                tpo:
                  type: integer
                  minimum: 1
                  maximum: 5
                color:
                  type: integer
                  minimum: 1
                  maximum: 15
                content:
                  type: string
                rating:
                  type: number
                  format: float
                  minimum: 0
                  maximum: 5
      responses:
        '201':
          description: 作成成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Item'
        '422':
          $ref: '#/components/responses/ValidationError'

  /items/{id}:
    get:
      tags:
        - items
      summary: アイテム詳細取得
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ItemWithUser'
        '404':
          $ref: '#/components/responses/NotFoundError'

    put:
      tags:
        - items
      summary: アイテム更新
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              properties:
                name:
                  type: string
                image:
                  type: string
                  format: binary
                super_item:
                  type: string
                  enum: [tops, bottoms, shoes, outer, accessory]
                season:
                  type: integer
                  minimum: 1
                  maximum: 5
                tpo:
                  type: integer
                  minimum: 1
                  maximum: 5
                color:
                  type: integer
                  minimum: 1
                  maximum: 15
                content:
                  type: string
                rating:
                  type: number
                  format: float
                  minimum: 0
                  maximum: 5
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Item'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    delete:
      tags:
        - items
      summary: アイテム削除
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '204':
          description: 削除成功
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

  # コーディネート管理エンドポイント
  /coordinates:
    get:
      tags:
        - coordinates
      summary: コーディネート一覧取得
      parameters:
        - name: user_id
          in: query
          schema:
            type: integer
            format: int64
        - name: season
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 5
        - name: tpo
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 5
        - name: min_rating
          in: query
          schema:
            type: number
            format: float
            minimum: 0
            maximum: 5
        - name: max_rating
          in: query
          schema:
            type: number
            format: float
            minimum: 0
            maximum: 5
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          headers:
            X-Total-Count:
              schema:
                type: integer
            X-Page:
              schema:
                type: integer
            X-Per-Page:
              schema:
                type: integer
          content:
            application/json:
              schema:
                type: object
                properties:
                  coordinates:
                    type: array
                    items:
                      $ref: '#/components/schemas/CoordinateWithDetails'

    post:
      tags:
        - coordinates
      summary: コーディネート作成
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - si_shoes
                - si_bottoms
                - si_tops
                - season
                - tpo
                - rating
              properties:
                si_shoes:
                  type: integer
                  format: int64
                si_bottoms:
                  type: integer
                  format: int64
                si_tops:
                  type: integer
                  format: int64
                si_outer:
                  type: integer
                  format: int64
                season:
                  type: integer
                  minimum: 1
                  maximum: 5
                tpo:
                  type: integer
                  minimum: 1
                  maximum: 5
                rating:
                  type: number
                  format: float
                  minimum: 0
                  maximum: 5
                content:
                  type: string
          multipart/form-data:
            schema:
              type: object
              properties:
                image:
                  type: string
                  format: binary
                si_shoes:
                  type: integer
                si_bottoms:
                  type: integer
                si_tops:
                  type: integer
                si_outer:
                  type: integer
                season:
                  type: integer
                tpo:
                  type: integer
                rating:
                  type: number
                content:
                  type: string
      responses:
        '201':
          description: 作成成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Coordinate'
        '422':
          $ref: '#/components/responses/ValidationError'

  /coordinates/{id}:
    get:
      tags:
        - coordinates
      summary: コーディネート詳細取得
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CoordinateWithDetails'
        '404':
          $ref: '#/components/responses/NotFoundError'

    put:
      tags:
        - coordinates
      summary: コーディネート更新
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                si_shoes:
                  type: integer
                si_bottoms:
                  type: integer
                si_tops:
                  type: integer
                si_outer:
                  type: integer
                season:
                  type: integer
                tpo:
                  type: integer
                rating:
                  type: number
                content:
                  type: string
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Coordinate'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    delete:
      tags:
        - coordinates
      summary: コーディネート削除
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '204':
          description: 削除成功
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

  /coordinates/timeline:
    get:
      tags:
        - coordinates
      summary: タイムライン取得
      description: フォローしているユーザーのコーディネートを取得
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  coordinates:
                    type: array
                    items:
                      $ref: '#/components/schemas/CoordinateWithDetails'

  /coordinates/{id}/like:
    post:
      tags:
        - coordinates
      summary: いいね
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: いいね成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Liked successfully

    delete:
      tags:
        - coordinates
      summary: いいね取り消し
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 取り消し成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Unliked successfully

  # ソーシャル機能エンドポイント
  /comments:
    get:
      tags:
        - social
      summary: コメント一覧取得
      parameters:
        - name: coordinate_id
          in: query
          required: true
          schema:
            type: integer
            format: int64
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  comments:
                    type: array
                    items:
                      $ref: '#/components/schemas/CommentWithUser'

    post:
      tags:
        - social
      summary: コメント投稿
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - coordinate_id
                - comment
              properties:
                coordinate_id:
                  type: integer
                  format: int64
                comment:
                  type: string
                  minLength: 1
      responses:
        '201':
          description: 投稿成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Comment'

  /comments/{id}:
    put:
      tags:
        - social
      summary: コメント更新
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - comment
              properties:
                comment:
                  type: string
                  minLength: 1
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Comment'
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

    delete:
      tags:
        - social
      summary: コメント削除
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '204':
          description: 削除成功
        '403':
          $ref: '#/components/responses/ForbiddenError'
        '404':
          $ref: '#/components/responses/NotFoundError'

  /users/{id}/follow:
    post:
      tags:
        - social
      summary: フォロー
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: フォロー成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Followed successfully

    delete:
      tags:
        - social
      summary: フォロー解除
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 解除成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Unfollowed successfully

  /users/{id}/followers:
    get:
      tags:
        - social
      summary: フォロワー一覧
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  users:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: integer
                        name:
                          type: string
                        avatar:
                          type: string
                        is_following:
                          type: boolean

  /users/{id}/following:
    get:
      tags:
        - social
      summary: フォロー中一覧
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  users:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: integer
                        name:
                          type: string
                        avatar:
                          type: string
                        is_following:
                          type: boolean

  /users/{id}/block:
    post:
      tags:
        - social
      summary: ブロック
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: ブロック成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Blocked successfully

    delete:
      tags:
        - social
      summary: ブロック解除
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 解除成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Unblocked successfully

  /blocks:
    get:
      tags:
        - social
      summary: ブロックリスト
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  users:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: integer
                        name:
                          type: string
                        avatar:
                          type: string

  # 通知エンドポイント
  /notifications:
    get:
      tags:
        - notifications
      summary: 通知一覧
      parameters:
        - name: unread_only
          in: query
          schema:
            type: boolean
            default: false
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  notifications:
                    type: array
                    items:
                      $ref: '#/components/schemas/NotificationWithDetails'
                  unread_count:
                    type: integer

  /notifications/{id}/read:
    put:
      tags:
        - notifications
      summary: 通知既読
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: 既読成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Marked as read

  /notifications/read_all:
    put:
      tags:
        - notifications
      summary: 全通知既読
      responses:
        '200':
          description: 既読成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: All notifications marked as read

  # ヘルスチェック
  /health:
    get:
      tags:
        - health
      summary: ヘルスチェック
      security: []
      responses:
        '200':
          description: サービス正常稼働中
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: ok
                  timestamp:
                    type: string
                    format: date-time
```

## APIクライアント実装例

### JavaScript (Axios)

```javascript
import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'https://api.speadwear.com/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエストインターセプター
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // トークン期限切れ処理
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 使用例
async function getCoordinates(params) {
  try {
    const response = await apiClient.get('/coordinates', { params });
    return response.data;
  } catch (error) {
    console.error('Error fetching coordinates:', error.response?.data);
    throw error;
  }
}
```

### Go (HTTPクライアント)

```go
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
    Token      string
}

func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        Token: token,
    }
}

func (c *APIClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    var reqBody []byte
    var err error
    
    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }
    
    req, err := http.NewRequest(method, c.BaseURL+endpoint, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Authorization", "Bearer "+c.Token)
    }
    
    return c.HTTPClient.Do(req)
}

// 使用例
func (c *APIClient) GetCoordinates(params map[string]string) (*CoordinateListResponse, error) {
    endpoint := "/coordinates"
    if len(params) > 0 {
        endpoint += "?"
        for k, v := range params {
            endpoint += fmt.Sprintf("%s=%s&", k, v)
        }
    }
    
    resp, err := c.doRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }
    
    var result CoordinateListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

## レート制限

APIには以下のレート制限が適用されます：

- 認証済みユーザー: 1時間あたり5,000リクエスト
- 未認証ユーザー: 1時間あたり60リクエスト

レート制限の状態は以下のレスポンスヘッダーで確認できます：

- `X-RateLimit-Limit`: 制限値
- `X-RateLimit-Remaining`: 残りリクエスト数
- `X-RateLimit-Reset`: リセット時刻（Unix timestamp）

## バージョニング

APIのバージョンはURLパスに含まれます（例: `/api/v1/`）。
新しいバージョンがリリースされても、既存のバージョンは最低1年間はサポートされます。

## 変更履歴

### v1.0.0 (2024-01-01)
- 初回リリース
- 基本的なCRUD操作
- 認証機能
- ソーシャル機能