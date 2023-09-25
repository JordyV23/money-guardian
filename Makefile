build:
	@go build -o bin/money-guardian

run: build
	@./bin/money-guardian

test:
	@go test -v ./...