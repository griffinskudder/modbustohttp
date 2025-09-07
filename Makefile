
BUF = docker run --rm -v "$(CURDIR):/workspace" -w /workspace bufbuild/buf
GOLANGCI_LINT = docker run -t --rm -v "$(CURDIR):/app" -w /app golangci/golangci-lint:v2.4.0 golangci-lint

generate:
	echo "Generating code with Buf..."
	echo "Current Working Directory: $(CURDIR)"
	$(BUF) generate

lint:
	echo "Linting with Buf..."
	echo "Current Working Directory: $(CURDIR)"
	$(BUF) lint
	echo "Linting with GolangCI-Lint..."
	$(GOLANGCI_LINT) run

clean:
	echo "Cleaning up..."
	rm -rf "$(CURDIR)/dist"

fmt:
	echo "Formatting code..."
	go fmt ./...

build:
	echo "Building the Go application..."
	mkdir -p dist
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/modbustohttp_windows-amd64.exe .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/modbustohttp_linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/modbustohttp_darwin-arm64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/modbustohttp_linux-arm64 .

test:
	echo "Running tests..."
	go test -v ./...