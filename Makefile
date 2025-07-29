.PHONY: proto build start-user-service start-todo-service start-bff run-bff start-client start-backend start-all start-user-service-bg start-todo-service-bg start-bff-bg stop stop-graceful stop-force stop-aggressive stop-docker-safe status clean help

# Variables
PORTS := 50051 50052 8080 5173
GO_PROCESSES := "go run.*cmd/server" "server/bff/bin/bff" "server/services/user/bin/user-service" "server/services/todo/bin/todo-service"
NODE_PROCESSES := "npm run dev" "vite.*client" "node.*todo.*client"

# Protocol Buffers compilation
proto:
	@echo "Compiling protocol buffers..."
	protoc --go_out=proto --go_opt=module=github.com/tadasy/todo-app/proto --go-grpc_out=proto --go-grpc_opt=module=github.com/tadasy/todo-app/proto proto/*.proto

# Build all services
build:
	@echo "Building all services..."
	cd server/bff && go build -o bin/bff ./cmd/server
	cd server/services/user && go build -o bin/user-service ./cmd/server
	cd server/services/todo && go build -o bin/todo-service ./cmd/server

# Start individual services
start-user-service:
	@echo "Starting User Service..."
	cd server/services/user && go run ./cmd/server

start-todo-service:
	@echo "Starting Todo Service..."
	cd server/services/todo && go run ./cmd/server

start-bff:
	@echo "Starting BFF..."
	cd server/bff && go run ./cmd/server

run-bff: start-bff

# Start individual services in background
start-user-service-bg:
	@echo "Starting User Service in background..."
	@(cd server/services/user && go run ./cmd/server) &

start-todo-service-bg:
	@echo "Starting Todo Service in background..."
	@(cd server/services/todo && go run ./cmd/server) &

start-bff-bg:
	@echo "Starting BFF in background..."
	@(cd server/bff && go run ./cmd/server) &

start-client:
	@echo "Starting Client..."
	cd client && npm run dev

# Start all backend services (user, todo, bff) concurrently
start-backend:
	@echo "Starting all backend services..."
	@echo "Starting User Service in background..."
	@(cd server/services/user && go run ./cmd/server) &
	@echo "Starting Todo Service in background..."
	@(cd server/services/todo && go run ./cmd/server) &
	@echo "Starting BFF Service..."
	@cd server/bff && go run ./cmd/server

# Start all services (backend + frontend) concurrently
start-all:
	@echo "Starting all services (backend + frontend)..."
	@echo "Starting User Service in background..."
	@(cd server/services/user && go run ./cmd/server) &
	@echo "Starting Todo Service in background..."
	@(cd server/services/todo && go run ./cmd/server) &
	@echo "Starting BFF Service in background..."
	@(cd server/bff && go run ./cmd/server) &
	@echo "Starting Frontend..."
	@cd client && npm run dev

# Stop all background processes with port verification (Docker-safe)
stop:
	@echo "Stopping all services (Docker-safe mode)..."
	@echo "Killing project processes only..."
	@for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			pid=$$(lsof -ti:$$port); \
			cmd=$$(ps -p $$pid -o args= 2>/dev/null || echo "unknown"); \
			if echo "$$cmd" | grep -q "$$(pwd)\|todo"; then \
				echo "  Killing project process on port $$port (PID: $$pid)"; \
				kill -9 $$pid 2>/dev/null || true; \
			else \
				echo "  Skipping non-project process on port $$port: $$cmd"; \
			fi; \
		fi; \
	done
	@sleep 1
	$(call verify-ports)
	@echo "All project services stopped safely."

# Graceful stop with SIGTERM first, then SIGKILL if needed
stop-graceful:
	@echo "Gracefully stopping all services..."
	@echo "Sending SIGTERM to processes..."
	$(call kill-processes,$(GO_PROCESSES),Go,)
	$(call kill-processes,$(NODE_PROCESSES),Node.js,)
	@echo "Waiting 3 seconds for graceful shutdown..."
	@sleep 3
	@echo "Force killing any remaining processes..."
	$(call kill-processes,$(GO_PROCESSES),Go,-9)
	$(call kill-processes,$(NODE_PROCESSES),Node.js,-9)
	$(call kill-processes-on-ports,-9)
	$(call verify-ports)
	@echo "All services stopped gracefully."

# Force stop with aggressive cleaning
stop-force:
	@echo "Force stopping all services and cleaning ports..."
	@pkill -9 -f "go run\|npm\|vite\|node" 2>/dev/null || true
	$(call kill-processes-on-ports,-9)
	@sleep 2
	@echo "Final verification..."
	@all_clear=true; \
	for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			echo "❌ Port $$port is still in use"; \
			all_clear=false; \
		else \
			echo "✅ Port $$port is available"; \
		fi; \
	done; \
	if [ "$$all_clear" = "true" ]; then \
		echo "🎉 All ports are clean and ready!"; \
	else \
		echo "⚠️  Some ports may still be in use. Manual intervention may be required."; \
	fi

# Aggressive stop (use with caution in Docker environments)
stop-aggressive:
	@echo "Stopping all services (aggressive mode)..."
	@echo "⚠️  This may affect other applications!"
	$(call kill-processes,$(GO_PROCESSES),Go,-9)
	$(call kill-processes,$(NODE_PROCESSES),Node.js,-9)
	@sleep 2
	$(call kill-processes-on-ports,-9)
	@sleep 1
	$(call verify-ports)
	@echo "All services stopped aggressively."

# Install dependencies
install:
	@echo "Installing dependencies..."
	cd client && npm install
	cd server/bff && go mod tidy
	cd server/services/user && go mod tidy
	cd server/services/todo && go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf server/bff/bin
	rm -rf server/services/user/bin
	rm -rf server/services/todo/bin
	rm -rf client/dist

# Development setup
dev-setup: install proto
	@echo "Development environment setup complete!"

# Check status of all services
status:
	@echo "Service Status Check:"
	@echo "====================="
	@for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			pid=$$(lsof -ti:$$port); \
			process=$$(ps -p $$pid -o comm= 2>/dev/null || echo "unknown"); \
			echo "✅ Port $$port: ACTIVE (PID: $$pid, Process: $$process)"; \
		else \
			echo "❌ Port $$port: INACTIVE"; \
		fi; \
	done
	@echo "====================="

# Docker-aware stop with process inspection
stop-docker-safe:
	@echo "Docker-aware safe stopping..."
	@echo "Inspecting processes before stopping..."
	@for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			pid=$$(lsof -ti:$$port); \
			cmd=$$(ps -p $$pid -o args= 2>/dev/null || echo "unknown"); \
			echo "Port $$port is used by PID $$pid"; \
			echo "  Command: $$cmd"; \
			if echo "$$cmd" | grep -q -E "docker|Docker"; then \
				echo "  🐳 Docker process detected - SKIPPING"; \
			elif echo "$$cmd" | grep -q -E "$$(pwd)|todo|go run|npm|vite"; then \
				echo "  ✅ Project process detected - KILLING"; \
				kill -9 $$pid 2>/dev/null || true; \
			else \
				echo "  ⚠️  Unknown process - SKIPPING for safety"; \
			fi; \
			echo ""; \
		fi; \
	done
	@sleep 1
	$(call verify-ports)
	@echo "Docker-safe stop completed."

# Show help for all available commands
help:
	@echo "📋 利用可能なコマンド:"
	@echo "=================="
	@echo ""
	@echo "🚀 開始コマンド:"
	@echo "  make start-all              - 全サービス起動（バックエンド + フロントエンド）"
	@echo "  make start-backend          - バックエンドサービスのみ起動"
	@echo "  make start-client           - フロントエンドのみ起動"
	@echo "  make start-user-service     - ユーザーサービスのみ起動"
	@echo "  make start-todo-service     - Todoサービスのみ起動"
	@echo "  make start-bff              - BFFのみ起動"
	@echo ""
	@echo "🛑 停止コマンド:"
	@echo "  make stop                   - 🐳 Docker対応安全停止（推奨）"
	@echo "  make stop-graceful          - 段階的停止（SIGTERMから開始）"
	@echo "  make stop-docker-safe       - プロセス詳細検査付き超安全停止"
	@echo "  make stop-aggressive        - ⚠️  積極的停止（他アプリに影響の可能性）"
	@echo "  make stop-force             - 💥 最終手段（緊急時のみ使用）"
	@echo ""
	@echo "🔧 ユーティリティコマンド:"
	@echo "  make status                 - サービス状態確認"
	@echo "  make build                  - 全サービスビルド"
	@echo "  make clean                  - ビルド成果物削除"
	@echo "  make install                - 依存関係インストール"
	@echo "  make proto                  - Protocol Buffersコンパイル"
	@echo "  make dev-setup              - 開発環境完全セットアップ"
	@echo ""
	@echo "💡 Docker環境での使用のコツ:"
	@echo "  - 日常的な開発では 'make stop' を使用（Docker対応済み）"
	@echo "  - プロセス競合が発生した場合は 'make stop-docker-safe' を試す"
	@echo "  - 'make stop-aggressive' は必要時のみ使用してください"

# Helper functions for process management
define kill-processes
	@echo "Killing $(2) processes..."
	@for pattern in $(1); do \
		echo "  Searching for pattern: $$pattern"; \
		pids=$$(pgrep -f "$$pattern" 2>/dev/null || true); \
		if [ ! -z "$$pids" ]; then \
			echo "  Found PIDs: $$pids"; \
			echo "$$pids" | xargs kill $(3) 2>/dev/null || true; \
		else \
			echo "  No processes found for pattern: $$pattern"; \
		fi; \
	done
endef

define kill-processes-on-ports
	@echo "Killing processes on specified ports..."
	@for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			echo "  Killing processes on port $$port..."; \
			lsof -ti:$$port | xargs kill $(1) 2>/dev/null || true; \
		fi; \
	done
endef

define verify-ports
	@echo "Verifying port availability..."
	@for port in $(PORTS); do \
		if lsof -ti:$$port >/dev/null 2>&1; then \
			echo "  ⚠️  Port $$port is still in use"; \
		else \
			echo "  ✅ Port $$port is available"; \
		fi; \
	done
endef
