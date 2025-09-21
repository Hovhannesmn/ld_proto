# Third-Party Service Integration Example

This example shows how to integrate the `ld_proto` package into your own Go service.

## Setup

1. **Install the package:**
   ```bash
   go get github.com/Hovhannesmn/ld_proto
   ```

2. **Add to your go.mod:**
   ```go
   require (
       github.com/Hovhannesmn/ld_proto v1.0.1
       google.golang.org/grpc v1.75.1
       google.golang.org/protobuf v1.36.9
   )
   ```

## Usage Patterns

### 1. Basic Integration

```go
import (
    "context"
    "google.golang.org/grpc"
    "github.com/Hovhannesmn/ld_proto/pb"
)

// Connect to language detection service
conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Create client
client := pb.NewLanguageDetectionServiceClient(conn)

// Use the client
req := &pb.DetectLanguageRequest{
    Text:       "Your text here",
    DocumentId: "doc-123",
    Metadata: map[string]string{
        "source": "your_service",
    },
}

resp, err := client.DetectLanguage(context.Background(), req)
```

### 2. Service Wrapper Pattern

Wrap the language detection client in your own service for better abstraction:

```go
type MyService struct {
    languageClient pb.LanguageDetectionServiceClient
}

func (s *MyService) ProcessText(text string) (*LanguageResult, error) {
    // Your business logic here
    // Call language detection
    // Process results
}
```

### 3. Error Handling

```go
resp, err := client.DetectLanguage(ctx, req)
if err != nil {
    // Handle gRPC errors
    if status.Code(err) == codes.Unavailable {
        // Service unavailable
    } else if status.Code(err) == codes.DeadlineExceeded {
        // Timeout
    }
    return err
}
```

### 4. Context with Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.DetectLanguage(ctx, req)
```

### 5. Batch Processing

```go
func (s *MyService) ProcessBatch(texts []string) ([]*LanguageResult, error) {
    var results []*LanguageResult
    
    for _, text := range texts {
        result, err := s.ProcessText(text)
        if err != nil {
            continue // or handle error
        }
        results = append(results, result)
    }
    
    return results, nil
}
```

## Running the Example

1. **Start the language detection server:**
   ```bash
   cd ../server
   go run main.go
   ```

2. **Run this third-party service:**
   ```bash
   go run main.go
   ```

## Integration Tips

1. **Connection Management:** Reuse gRPC connections instead of creating new ones for each request
2. **Error Handling:** Always handle gRPC errors appropriately
3. **Timeouts:** Use context with timeouts to prevent hanging requests
4. **Metadata:** Use metadata to pass additional context to the language detection service
5. **Batch Processing:** Consider batching requests for better performance
6. **Monitoring:** Add logging and metrics for language detection calls
