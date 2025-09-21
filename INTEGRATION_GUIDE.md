# Integration Guide for Third-Party Services

This guide shows how to integrate the `ld_proto` package into your Go services.

## Quick Start

### 1. Install the Package

```bash
go get github.com/Hovhannesmn/ld_proto
```

### 2. Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "github.com/Hovhannesmn/ld_proto/pb"
)

func main() {
    // Connect to language detection service
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Create client
    client := pb.NewLanguageDetectionServiceClient(conn)
    
    // Detect language
    req := &pb.DetectLanguageRequest{
        Text:       "Hello, world!",
        DocumentId: "doc-123",
    }
    
    resp, err := client.DetectLanguage(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Language: %s (%.2f confidence)", resp.LanguageCode, resp.Confidence)
}
```

## Common Integration Patterns

### Pattern 1: Service Wrapper

```go
type DocumentProcessor struct {
    languageClient pb.LanguageDetectionServiceClient
}

func NewDocumentProcessor(conn *grpc.ClientConn) *DocumentProcessor {
    return &DocumentProcessor{
        languageClient: pb.NewLanguageDetectionServiceClient(conn),
    }
}

func (dp *DocumentProcessor) ProcessDocument(content string) (*pb.DetectLanguageResponse, error) {
    req := &pb.DetectLanguageRequest{
        Text:       content,
        DocumentId: generateID(),
        Metadata: map[string]string{
            "processor": "document_processor",
        },
    }
    
    return dp.languageClient.DetectLanguage(context.Background(), req)
}
```

### Pattern 2: Middleware Integration

```go
func LanguageDetectionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract text from request
        text := extractTextFromRequest(r)
        
        // Detect language
        req := &pb.DetectLanguageRequest{Text: text}
        resp, err := languageClient.DetectLanguage(r.Context(), req)
        if err == nil {
            // Add language info to request context
            ctx := context.WithValue(r.Context(), "language", resp.LanguageCode)
            r = r.WithContext(ctx)
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### Pattern 3: Background Processing

```go
type LanguageDetectionWorker struct {
    client    pb.LanguageDetectionServiceClient
    jobQueue  chan DocumentJob
}

func (w *LanguageDetectionWorker) Start() {
    for job := range w.jobQueue {
        go w.processJob(job)
    }
}

func (w *LanguageDetectionWorker) processJob(job DocumentJob) {
    req := &pb.DetectLanguageRequest{
        Text:       job.Content,
        DocumentId: job.ID,
    }
    
    resp, err := w.client.DetectLanguage(context.Background(), req)
    if err != nil {
        job.ErrorCallback(err)
        return
    }
    
    job.SuccessCallback(resp)
}
```

## Error Handling

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

resp, err := client.DetectLanguage(ctx, req)
if err != nil {
    switch status.Code(err) {
    case codes.Unavailable:
        // Service is down, implement retry logic
        return retryWithBackoff()
    case codes.DeadlineExceeded:
        // Request timed out
        return errors.New("language detection timeout")
    case codes.InvalidArgument:
        // Invalid request
        return errors.New("invalid request parameters")
    default:
        return fmt.Errorf("language detection failed: %w", err)
    }
}
```

## Configuration

### Environment Variables

```go
var (
    LanguageServiceAddr = os.Getenv("LANGUAGE_SERVICE_ADDR")
    if LanguageServiceAddr == "" {
        LanguageServiceAddr = "localhost:50051"
    }
)
```

### Connection Options

```go
conn, err := grpc.Dial(LanguageServiceAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             3 * time.Second,
        PermitWithoutStream: true,
    }),
)
```

## Testing

### Mock Implementation

```go
type MockLanguageDetectionClient struct {
    responses map[string]*pb.DetectLanguageResponse
}

func (m *MockLanguageDetectionClient) DetectLanguage(ctx context.Context, req *pb.DetectLanguageRequest, opts ...grpc.CallOption) (*pb.DetectLanguageResponse, error) {
    if resp, exists := m.responses[req.Text]; exists {
        return resp, nil
    }
    return &pb.DetectLanguageResponse{
        LanguageCode: "en",
        Confidence:   0.95,
    }, nil
}
```

## Production Considerations

1. **Connection Pooling:** Reuse gRPC connections
2. **Circuit Breaker:** Implement circuit breaker pattern for resilience
3. **Retries:** Add exponential backoff for retries
4. **Monitoring:** Add metrics and logging
5. **Timeouts:** Always use context with timeouts
6. **Graceful Shutdown:** Properly close connections on shutdown

## Example Projects

See the `examples/` directory for complete working examples:
- `client/` - Basic client usage
- `server/` - Server implementation
- `third_party_service/` - Integration example
