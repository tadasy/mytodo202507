module github.com/tadasy/mytodo202507/server/bff

go 1.23.0

toolchain go1.24.5

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/labstack/echo/v4 v4.11.2
	github.com/tadasy/mytodo202507/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.74.2
)

replace github.com/tadasy/mytodo202507/proto => ../../proto

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
