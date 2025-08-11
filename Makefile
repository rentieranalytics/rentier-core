.PHONY: tidy test cover lint

tidy:
	go mod tidy

test:
	go test ./... -v

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	go vet ./...
