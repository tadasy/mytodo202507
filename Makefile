.PHONY: proto build start-user-service start-todo-service start-bff run-bff start-client start-backend start-all start-user-service-bg start-todo-service-bg start-bff-bg stop clean

# Protocol Buffers compilation
proto:
	@echo "Compiling protocol buffers..."
	protoc --go_out=. --go-grpc_out=. proto/*.proto

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

# Stop all background processes
stop:
	@echo "Stopping all services..."
	@pkill -f "go run.*cmd/server" || true
	@pkill -f "npm run dev" || true
	@echo "All services stopped."

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
