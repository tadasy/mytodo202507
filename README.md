# Todo Management Application

マイクロサービスアーキテクチャで構築されたTodo管理アプリケーション

## アーキテクチャ

### フロントエンド
- **Client**: TypeScript + React (React Router v7)
- **UI**: shadcn/ui デザインシステム

### バックエンド
- **BFF**: Go + Echo (REST API)
- **User Service**: Go + gRPC (ユーザー管理マイクロサービス)
- **Todo Service**: Go + gRPC (Todo管理マイクロサービス)
- **Database**: SQLite

### 通信方式
- Client ↔ BFF: REST API (OpenAPI)
- BFF ↔ Microservices: gRPC

## プロジェクト構造

```
.
├── client/                 # React frontend
├── server/
│   ├── bff/               # Backend for Frontend
│   └── services/
│       ├── user/          # User management service
│       └── todo/          # Todo management service
├── proto/                 # gRPC protocol definitions
└── docs/                  # Documentation
```

## 機能要件

### ユーザー管理
- ユーザー登録
- ログイン/ログアウト (email, password)

### Todo管理
- Todo作成、更新、削除
- Todo完了マーク
- 完了済みTodo一覧表示
- ユーザーごとのTodo管理

## 開発環境のセットアップ

### 前提条件
- Node.js 18+
- Go 1.21+
- Protocol Buffers compiler (protoc)

### セットアップ手順

1. 依存関係のインストール
```bash
# フロントエンド
cd client && npm install

# バックエンド各サービス
cd server/bff && go mod tidy
cd server/services/user && go mod tidy
cd server/services/todo && go mod tidy
```

2. プロトコルバッファのコンパイル
```bash
make proto
```

3. アプリケーションの起動

## 方法1: すべてのサービスを一度に起動
```bash
# バックエンド（User Service + Todo Service + BFF）を起動
make start-backend

# すべて（バックエンド + フロントエンド）を一度に起動  
make start-all

# サービスを停止
make stop
```

## 方法2: 各サービスを個別に起動
```bash
# 各サービスを別々のターミナルで起動
make start-user-service
make start-todo-service
make start-bff
make start-client
```

## 方法3: バックグラウンドで起動
```bash
# 各サービスをバックグラウンドで起動
make start-user-service-bg
make start-todo-service-bg
make start-bff-bg

# フロントエンドのみフォアグラウンドで起動
make start-client
```
