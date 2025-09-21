# Examples

This directory contains example implementations showing how to use the `ld_proto` package.

## Running the Examples

### 1. Start the Server

```bash
cd examples
go run server/main.go
```

The server will start on port 50051 and provide a simple language detection service.

### 2. Run the Client

In a separate terminal:

```bash
cd examples
go run client/main.go
```

The client will send several test phrases to the server and display the detected languages.

## Example Output

```
Text: 'Hello, world!'
Detected language: en (confidence: 0.20)
Document ID: example-doc-20241220-143022
Processing time: 0ms
Service version: 1.0.0
---
Text: 'Hola, mundo!'
Detected language: es (confidence: 0.20)
Document ID: example-doc-20241220-143022
Processing time: 0ms
Service version: 1.0.0
---
```

## Custom Implementation

You can modify the server implementation in `server/main.go` to integrate with your preferred language detection library or API service.

The client example in `client/main.go` shows how to:
- Connect to a gRPC service
- Create requests with metadata
- Handle responses and alternatives
- Use context for timeouts
