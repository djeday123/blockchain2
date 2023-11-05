build:
	@go build -o bin/blockchain2

run: build
	@./bin/docker

test:
	@go test -v ./...