module github.com/tadasy/todo-app/server/services/todo

go 1.23.0

toolchain go1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/tadasy/todo-app/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
)

replace github.com/tadasy/todo-app/proto => ../../../proto

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
