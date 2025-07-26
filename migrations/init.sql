-- 初期データベース設定
-- このファイルはDockerコンテナ起動時に自動的に実行されます

-- データベースの文字セットをUTF8MB4に設定
ALTER DATABASE speadwear_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ユーザーへの権限付与（開発環境用）
-- 本番環境では適切な権限設定を行ってください
GRANT ALL PRIVILEGES ON speadwear_dev.* TO 'speadwear'@'%';
FLUSH PRIVILEGES;