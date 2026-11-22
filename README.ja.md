# HaNas

[한국어](README.ko.md) | [English](README.md)

HaNasはGoとSQLiteを使用するセルフホスト型の個人ファイルストレージです。ブラウザー上でアカウント管理、ファイル・フォルダー操作、ごみ箱、共有リンク、メディアプレビューを利用できます。

## 主な機能

- ユーザー登録とログイン
- JWTトークンを使用した認証Cookie
- フォルダーとファイルの作成
- ファイルのアップロードと上書き
- コピー、移動、名前変更、削除
- ごみ箱への移動と30日後の完全削除
- ファイル共有リンクの作成と削除
- 画像・動画サムネイルの生成
- 画像、動画、音声、ドキュメントのプレビュー
- アカウントと関連データの削除

## 技術スタック

- Go `1.25.1`
- `net/http`
- GORMとglebarez SQLite
- JWT
- bcrypt
- Reactベースの静的フロントエンド
- オプション: `ffmpeg`

フロントエンドの静的ファイルはGoバイナリに埋め込まれるため、別途フロントエンドのビルドは不要です。

## 必要環境

- Go 1.25以上
- ポート`80`を使用できる権限
- 基本機能に追加のシステム依存関係は不要
- 動画サムネイルを使用する場合は`ffmpeg`が`PATH`に必要

macOSまたはLinuxでポート80を使用するには、管理者権限が必要になる場合があります。

```bash
brew install go ffmpeg
```

上記はmacOSの例です。Goとffmpegは各OSのパッケージマネージャーでインストールできます。

## クイックスタート

プロジェクトルートでビルドします。

```bash
cd /Users/user/Desktop/server
go mod download
go build -o hanas .
```

実行ファイルと同じフォルダーに`config.js`を配置してサーバーを起動します。

```bash
./hanas
```

デフォルト設定は次のとおりです。

```javascript
window.APP_CONFIG = {
  BACKEND_URL: 'http://127.0.0.1:80'
}
```

ブラウザーで`http://127.0.0.1`またはサーバーのホストアドレスにアクセスします。

ポート80を使用できない場合は、一時的に権限を付与して実行できます。

```bash
sudo ./hanas
```

本番環境ではTLSリバースプロキシの後ろに配置し、`config.js`の`BACKEND_URL`を公開URLへ変更することを推奨します。

## 実行設定

`config.js`は実行ファイルがあるディレクトリから読み込まれます。バイナリを別の場所へ配置する場合は、実行ファイルと同じ場所に設定ファイルをコピーしてください。

例:

```text
/opt/hanas/
├── hanas
├── config.js
├── database.db
├── data/
└── thumbnails/
```

サーバーは常に`:80`で待ち受けます。環境変数でポートを変更する機能はありません。

## データ保存先

アプリケーションはカレントディレクトリを基準に次のパスを使用します。

| パス | 用途 |
| --- | --- |
| `database.db` | ユーザー、ノード、共有、ごみ箱のメタデータ |
| `data/` | アップロードされた元ファイル |
| `thumbnails/` | 生成されたサムネイル |

データをバックアップする際は、次の3項目をすべて保存してください。

```bash
mkdir -p backups
cp database.db backups/database.db
cp -R data backups/data
cp -R thumbnails backups/thumbnails
```

稼働中のSQLiteデータベースを直接コピーすると不整合が発生する可能性があるため、バックアップ前にサーバーを停止するか、SQLite対応のバックアップ手順を使用してください。

## プロジェクト構成

```text
.
├── app.go            # サーバー起動とHTTPルート
├── logic.go          # 認証、ファイル、共有API
├── exec.go           # ファイル操作とデータ処理
├── object.go         # データモデル
├── go.mod            # Go依存関係
├── index.html        # Webアプリケーションのエントリーポイント
├── config.js         # フロントエンドのバックエンドアドレス
├── assets/
│   ├── index.js      # ビルド済みフロントエンド
│   ├── index.css     # フロントエンドスタイル
│   ├── icon.ico
│   └── apple-touch-icon.png
├── data/             # アップロードファイル保存先
├── thumbnails/       # サムネイル保存先
└── database.db       # SQLiteデータベース
```

## 動作の仕組み

1. 起動時にSQLiteデータベースを開き、必要なテーブルを自動作成します。
2. 初回起動時にJWT署名用の秘密鍵を生成し、データベースへ保存します。
3. ユーザー登録時にユーザーごとのルートフォルダーを作成します。
4. ログイン成功時にJWTを`HttpOnly`かつ`SameSite=Strict`のCookieへ保存します。
5. アップロードファイルは`data/`に数値のファイルIDで保存され、ノード情報はSQLiteへ保存されます。
6. ごみ箱へ移動した項目はバックグラウンド処理により30日後に完全削除されます。

JWT署名キーはデータベースに保存されるため、既存セッションとの互換性を維持する場合は再デプロイ時に`database.db`も保持してください。

## 主なAPI

ほとんどのユーザーAPIには`token` Cookieが必要です。

| メソッド | パス | 説明 |
| --- | --- | --- |
| `POST` | `/register` | アカウント作成 |
| `POST` | `/login` | ログイン |
| `POST` | `/logout` | ログアウト |
| `GET` | `/me` | 現在のユーザー取得 |
| `GET` | `/node/` | ルートフォルダー取得 |
| `GET` | `/node/{id}` | フォルダー取得 |
| `GET` | `/file/{nodeId}` | ファイルダウンロード |
| `GET` | `/thumbnail/{nodeId}` | サムネイル取得 |
| `POST` | `/upload` | ファイルアップロードまたはフォルダー作成 |
| `POST` | `/copy` | 項目のコピー |
| `POST` | `/move` | 項目の移動 |
| `POST` | `/rename` | 項目の名前変更 |
| `POST` | `/delete` | 項目の完全削除 |
| `POST` | `/trash/in` | ごみ箱へ移動 |
| `GET` | `/trash` | ごみ箱一覧 |
| `GET` | `/trash/info/{id}` | ごみ箱項目の詳細 |
| `POST` | `/trash/restore/{id}` | ごみ箱項目の復元 |
| `POST` | `/share/create` | 共有リンク作成 |
| `POST` | `/share/delete` | 共有リンク削除 |
| `GET` | `/s/{token}` | 共有ファイルへのアクセス |
| `POST` | `/delete-account` | アカウントとデータの削除 |

## アップロード例

JSONリクエストで小さなファイルをアップロードできます。

```bash
curl -b cookie.txt \
  -H 'Content-Type: application/json' \
  -d '{
    "filename": "hello.txt",
    "is_dir": false,
    "data_base64": "aGVsbG8="
  }' \
  http://127.0.0.1/upload
```

multipartアップロードの例:

```bash
curl -b cookie.txt \
  -F 'filename=photo.jpg' \
  -F 'oya_id=1' \
  -F 'file=@photo.jpg' \
  http://127.0.0.1/upload
```

ログインレスポンスで取得したCookieを使用してください。

## 開発と検証

依存関係を取得してビルドします。

```bash
go mod download
go build ./...
```

このプロジェクトには現在、組み込みのGoテストスイートがありません。
