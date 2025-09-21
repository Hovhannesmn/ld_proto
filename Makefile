# Proto Package Makefile

.PHONY: proto clean help

# Generate Go code from protobuf definitions
proto:
	mkdir -p pb && protoc \
		-I=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/language_detection.proto

# Clean generated files
clean:
	rm -f proto/*.pb.go

# Install protobuf compiler and Go plugins
install-deps:
	@echo "Installing protobuf compiler and Go plugins..."
	@echo "Please install protoc: https://grpc.io/docs/protoc-installation/"
	@echo "Then run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	@echo "And: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"

# Help
help:
	@echo "Available commands:"
	@echo "  proto         - Generate Go code from protobuf"
	@echo "  clean         - Clean generated files"
	@echo "  install-deps  - Show instructions for installing dependencies"
	@echo "  help          - Show this help message"