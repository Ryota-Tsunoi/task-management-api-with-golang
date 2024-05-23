# タスク管理 API

これは、Go と Echo Web フレームワークを使用して構築されたシンプルなタスク管理 API です。ユーザーはタスクの作成、読み取り、更新、削除を行うことができます。

## 目次

- セットアップ
  - 前提条件
  - インストール
- 開発
  - API の実行
  - API のビルド
  - テスト
  - コード構造
  - Git Hooks
- API エンドポイント
- 貢献
- ライセンス

## セットアップ

### 前提条件

- Go (バージョン 1.16 以降)
- Air (開発中のホットリロード用)

### インストール

1. リポジトリをクローンします:

   ```bash
   git clone https://github.com/yourusername/task-management-api.git
   ```

2. プロジェクトディレクトリに移動します:

   ```bash
   cd task-management-api
   ```

3. 依存関係をインストールします:

   ```bash
   go mod download
   ```

4. ホットリロード用に Air をインストールします:

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

5. Git hooks 用の LeftHook をインストールします:

   ```bash
   go install github.com/evilmartians/lefthook@latest

   lefthook install
   ```

6. セキュリティ向上のため、Gosec をインストールします:

   ```bash
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   ```

## 開発

### API の実行

ホットリロードを使用して API を実行するには、以下のコマンドを使用します:

```bash
make run
```

これにより、<http://localhost:18080> で API サーバーが起動します。

### API のビルド

API をビルドするには、以下のコマンドを使用します:

```bash
make build
```

これにより、bin ディレクトリに app という名前の実行可能ファイルが生成されます。

### テスト

テストを実行するには、以下のコマンドを使用します:

```bash
go test ./...
```

### コード構造

プロジェクトは以下のコード構造に従っています:

```plaintext
task-management-api/
├── cmd/
│   └── server/
│       └── main.go
├── pkg/
│   ├── handlers/
│   ├── models/
│   └── repositories/
├── .gitignore
├── go.mod
├── Makefile
└── README.md
```

- `cmd/server/main.go`: API サーバーのエントリーポイントです。
- `pkg/handlers`: API ハンドラー関数が含まれています。
- `pkg/models`: データモデルが含まれています。
- `pkg/repositories`: データベースリポジトリが含まれています。

## Git Hooks

このプロジェクトでは、[Lefthook](https://github.com/evilmartians/lefthook)を使用して Git hooks を管理しています。以下のフックが設定されています:

### pre-commit

- `lint`: `go vet`コマンドを使用してコードの静的解析を行います。
- `test`: `go test`コマンドを使用してテストを実行します。
- `fmt`: `go fmt`コマンドを使用してコードのフォーマットを整えます。
- `mod`: `go mod tidy`コマンドを使用して、使用されていない依存関係を削除し、`go.mod`ファイルと`go.sum`ファイルを更新します。

### pre-push

- `test`: `go test`コマンドを使用してテストを実行します。
- `security`: `gosec`コマンドを使用して、コードのセキュリティ脆弱性をスキャンします。

## API エンドポイント

以下の API エンドポイントが利用可能です:

- `GET /tasks`: すべてのタスクを取得します。
- `GET /tasks/:id`: 指定した ID のタスクを取得します。
- `POST /tasks`: 新しいタスクを作成します。
- `PUT /tasks/:id`: タスクを更新します。
- `DELETE /tasks/:id`: タスクを削除します。

## 貢献

貢献は大歓迎です！問題を見つけたり、改善のための提案がある場合は、Issue を開くかプルリクエストを送信してください。

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています。
