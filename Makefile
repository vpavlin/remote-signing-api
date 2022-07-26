generate:
	go generate ./...
fmt:
	go fmt ./...

serve:
	go run server/server.go config.json

build: fmt
	go build -o remote-signing-api server/server.go 
