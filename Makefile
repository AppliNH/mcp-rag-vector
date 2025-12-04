download-deps:
	go mod download

install-goa:
	go install goa.design/goa/v3/cmd/goa@latest

install-mockery:
	go install github.com/vektra/mockery/v2@latest

generate-code-design:
	goa gen github.com/applinh/mcp-rag-vector/design


install-lint: 
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-dev-tools: install-goa install-mockery install-lint

tidy:
	@go mod tidy

lint: 
	@echo "Running golangci-lint..."
	@golangci-lint run

lint-fix: 
	@echo "Fixing lint issues..."
	@golangci-lint run --fix

generate-mocks:
	@echo "Generating mocks..."
	@mockery

test: 
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

run-server:
	@go run main.go server

run-mcp:
	@go run main.go mcp

build: 
	@echo "Building server..."
	@mkdir -p bin
	@CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/server .

install-bin:
	@echo "Installing binary to GOPATH/bin..."
	@go install ./...

docker-build:
	@echo "Building Docker image..."
	@docker build -t  .

http-call-mcp:
	@SESSION_ID=$$(curl -s -i -X POST "http://localhost:3000/mcp" \
		-H "Content-Type: application/json" \
		-d '{ \
			"jsonrpc": "2.0", \
			"id": "0", \
			"method": "initialize" \
		}' | grep "mcp-session-" | tr -d '\r'); \
	curl -X POST http://localhost:3000/mcp \
		-H "Content-Type: application/json" \
		-H "$$SESSION_ID" \
		-H "Mcp-Protocol-Version: 2025-06-18" \
		-d '{ \
			"jsonrpc": "2.0", \
			"id": "0", \
			"method": "tools/call", \
			"params": { \
				"name": "greet", \
				"arguments":  { "name": "Thomas" } \
			} \
		}'