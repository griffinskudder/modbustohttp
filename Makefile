
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
	go build -o app .