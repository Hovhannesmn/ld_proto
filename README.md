# ld_proto

A Go package providing Protocol Buffer definitions and generated code for language detection services.

## Installation

```bash
go get github.com/hovman/ld_proto
```

## Overview

This package contains:
- Protocol Buffer definitions for language detection services
- Generated Go code for gRPC client and server implementations
- Message types for language detection requests and responses

## Usage

### Import the package

```go
import "github.com/hovman/ld_proto/pb"
```

### Example: Creating a gRPC Client

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "github.com/hovman/ld_proto/pb"
)

func main() {
    // Connect to the language detection service
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Create a client
    client := pb.NewLanguageDetectionServiceClient(conn)
    
    // Create a request
    req := &pb.DetectLanguageRequest{
        Text:       "Hello, world!",
        DocumentId: "doc-123",
        Metadata: map[string]string{
            "source": "user_input",
        },
    }
    
    // Call the service
    resp, err := client.DetectLanguage(context.Background(), req)
    if err != nil {
        log.Fatalf("Language detection failed: %v", err)
    }
    
    log.Printf("Detected language: %s (confidence: %.2f)", 
        resp.LanguageCode, resp.Confidence)
}
```

### Example: Creating a gRPC Server

```go
package main

import (
    "context"
    "log"
    "net"
    
    "google.golang.org/grpc"
    "github.com/hovman/ld_proto/pb"
)

type server struct {
    pb.UnimplementedLanguageDetectionServiceServer
}

func (s *server) DetectLanguage(ctx context.Context, req *pb.DetectLanguageRequest) (*pb.DetectLanguageResponse, error) {
    // Implement your language detection logic here
    return &pb.DetectLanguageResponse{
        LanguageCode: "en",
        Confidence:   0.95,
        DocumentId:   req.DocumentId,
        Metadata: &pb.ProcessingMetadata{
            ProcessingTimeMs: 50,
            ServiceVersion:   "1.0.0",
            ModelVersion:     "v1.2",
            Provider:         "custom",
        },
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    pb.RegisterLanguageDetectionServiceServer(s, &server{})
    
    log.Println("Server starting on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

## Message Types

### DetectLanguageRequest
- `text`: The text to analyze
- `document_id`: Optional document identifier
- `metadata`: Additional metadata as key-value pairs

### DetectLanguageResponse
- `language_code`: Detected language code (e.g., "en", "es", "fr")
- `confidence`: Confidence score (0.0 to 1.0)
- `alternatives`: List of alternative language predictions
- `document_id`: Echo of the document ID from request
- `metadata`: Processing metadata including timing and version info

## Dependencies

- `google.golang.org/grpc`: gRPC framework
- `google.golang.org/protobuf`: Protocol Buffer support

## License

This project is licensed under the MIT License.
