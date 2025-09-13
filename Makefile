NAME=narrator

.PHONY: all build run clean test cli gui

all: build

build:
	go build -o $(NAME)-go.exe .

run:
	./$(NAME)-go.exe

clean:
	rm -f $(NAME)-go.exe $(NAME)-cli.exe $(NAME)-gui.exe

test:
	go test ./...

cli: build-cli
	./$(NAME)-cli.exe

build-cli:
	go build -o $(NAME)-cli.exe ./cmd/cli

gui: build-gui
	./$(NAME)-gui.exe

build-gui:
	go build -o $(NAME)-gui.exe ./cmd/gui

# Build for Windows
build-windows:
	SET CGO_ENABLED=0
	SET GOOS=windows
	SET GOARCH=amd64
	go build -o bin/narrator-cli.exe ./cmd/cli
	go build -o bin/narrator-gui.exe ./cmd/gui

# Build for Linux
build-linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o bin/narrator-cli ./cmd/cli
	go build -o bin/narrator-gui ./cmd/gui

# Build for macOS
build-macos:
	SET CGO_ENABLED=0
	SET GOOS=darwin
	SET GOARCH=amd64
	go build -o bin/narrator-cli-macos ./cmd/cli
	go build -o bin/narrator-gui-macos ./cmd/gui

# Default build (for current OS)
default-build:
	go build -o bin/narrator-cli-$(GOOS)-$(GOARCH).exe ./cmd/cli
	go build -o bin/narrator-gui-$(GOOS)-$(GOARCH).exe ./cmd/gui

# Install dependencies
install-deps:
	go mod tidy

# Run all tests
test-all:
	go test -v ./...

# Run specific test file
test-file:
	go test -v $(TEST_FILE)

# Run specific test function
test-func:
	go test -v -run $(TEST_FUNC) $(TEST_FILE)

# Clean build artifacts
clean-build:
	rm -rf bin/*

# Clean all generated files
clean-all:
	$(MAKE) clean-build
	rm -f $(NAME)-go.exe $(NAME)-cli.exe $(NAME)-gui.exe
	rm -f assets/demo/*.wav
	rm -f mixed.wav

# Format Go code
fmt:
	go fmt ./...

# Lint Go code
lint:
	go vet ./...

# Generate mocks (if using mockery or similar)
# generate-mocks:
#	mockery --all --output internal/mocks

# Run the GUI application (requires build-gui first)
run-gui:
	./bin/narrator-gui.exe

# Run the CLI application (requires build-cli first)
run-cli:
	./bin/narrator-cli.exe

# Example of a CI script (local)
ci-local:
	$(MAKE) clean-all
	$(MAKE) install-deps
	$(MAKE) build-windows
	$(MAKE) test-all
	# Add CLI smoke test here
	# Add artifact checks here

