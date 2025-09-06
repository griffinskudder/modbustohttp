
BUF = docker run --rm -v "$(CURDIR):/workspace" -w /workspace bufbuild/buf

generate:
	echo "Generating code with Buf..."
	echo "Current Working Directory: $(CURDIR)"
	$(BUF) generate

lint:
	echo "Linting with Buf..."
	echo "Current Working Directory: $(CURDIR)"
	$(BUF) lint

build:
	echo "Building the Go application..."
	mkdir dist
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/modbustohttp_windows-amd64.exe .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/modbustohttp_linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/modbustohttp_darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/modbustohttp_darwin-arm64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dist/modbustohttp_linux-arm64 .
