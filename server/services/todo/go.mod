module github.com/tadasy/todo-app/server/services/todo

go 1.21

require (
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/tadasy/todo-app/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.58.3
)

replace github.com/tadasy/todo-app/proto => ../../../proto

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
