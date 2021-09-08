## Run tests
test:
	@go test -race ./...

## Run tests with coverage
test-cover:
	@go test -coverprofile=cover.out -race ./...

lint:
	@go fmt ./... && go vet ./...

## Fetch dependencies
fetch:
	@echo Download go.mod dependencies
	@go mod download

install-tools: fetch
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
