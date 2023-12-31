build:
	@go build -o bin/blockchain2

run: build
	@./bin/blockchain2

test:
	@go test -v ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: proto