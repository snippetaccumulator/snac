set dotenv-load

default: test build

# Build the CLI and Server
build: build-cli build-server
alias b := build

# Build CLI binary
build-cli:
    @echo "Building CLI"
    @go build -o $CLI_BIN ./cmd/cli/.

# Run tests in all packages
test-all:
    @echo "Testing All"
    @go test ./...

# Run tests from CLI and Server
test: test-cli test-server
alias t := test

# Run tests from CLI
test-cli:
    @echo "Testing CLI"
    @go test ./cmd/cli/...

# Build Server binary
build-server:
    @echo "Building Server"
    @go build -o $SERVER_BIN ./cmd/web/.

# Run tests from Server
test-server:
    @echo "Testing Server"
    @go test ./cmd/web/...

# Run the Server
watch:
    @echo "Watching Web"
    @air -c .air.toml
alias w := watch
