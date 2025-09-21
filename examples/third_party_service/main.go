package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/Hovhannesmn/ld_proto/pb"
)

// ThirdPartyService demonstrates how to integrate ld_proto in your own service
type ThirdPartyService struct {
	languageClient pb.LanguageDetectionServiceClient
}

// NewThirdPartyService creates a new service instance
func NewThirdPartyService(grpcConn *grpc.ClientConn) *ThirdPartyService {
	return &ThirdPartyService{
		languageClient: pb.NewLanguageDetectionServiceClient(grpcConn),
	}
}

// ProcessDocument processes a document and detects its language
func (s *ThirdPartyService) ProcessDocument(ctx context.Context, content, docID string) (*DocumentInfo, error) {
	// Create language detection request
	req := &pb.DetectLanguageRequest{
		Text:       content,
		DocumentId: docID,
		Metadata: map[string]string{
			"service":     "third_party_service",
			"version":     "1.0.0",
			"timestamp":   time.Now().Format(time.RFC3339),
			"content_len": string(rune(len(content))),
		},
	}

	// Call the language detection service
	resp, err := s.languageClient.DetectLanguage(ctx, req)
	if err != nil {
		return nil, err
	}

	// Process the response
	docInfo := &DocumentInfo{
		ID:           docID,
		Content:      content,
		Language:     resp.LanguageCode,
		Confidence:   resp.Confidence,
		ProcessedAt:  time.Now(),
		Alternatives: make([]LanguageAlternative, len(resp.Alternatives)),
	}

	// Convert alternatives
	for i, alt := range resp.Alternatives {
		docInfo.Alternatives[i] = LanguageAlternative{
			Language:   alt.LanguageCode,
			Confidence: alt.Confidence,
		}
	}

	// Add processing metadata if available
	if resp.Metadata != nil {
		docInfo.ProcessingTime = time.Duration(resp.Metadata.ProcessingTimeMs) * time.Millisecond
		docInfo.ServiceVersion = resp.Metadata.ServiceVersion
		docInfo.ModelVersion = resp.Metadata.ModelVersion
		docInfo.Provider = resp.Metadata.Provider
	}

	return docInfo, nil
}

// DocumentInfo represents a processed document
type DocumentInfo struct {
	ID             string
	Content        string
	Language       string
	Confidence     float32
	ProcessedAt    time.Time
	Alternatives   []LanguageAlternative
	ProcessingTime time.Duration
	ServiceVersion string
	ModelVersion   string
	Provider       string
}

// LanguageAlternative represents an alternative language detection
type LanguageAlternative struct {
	Language   string
	Confidence float32
}

// BatchProcessDocuments processes multiple documents
func (s *ThirdPartyService) BatchProcessDocuments(ctx context.Context, documents []Document) ([]*DocumentInfo, error) {
	results := make([]*DocumentInfo, 0, len(documents))
	
	for _, doc := range documents {
		docInfo, err := s.ProcessDocument(ctx, doc.Content, doc.ID)
		if err != nil {
			log.Printf("Failed to process document %s: %v", doc.ID, err)
			continue
		}
		results = append(results, docInfo)
	}
	
	return results, nil
}

// Document represents an input document
type Document struct {
	ID      string
	Content string
}

func main() {
	// 1. Connect to the language detection service
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to language detection service: %v", err)
	}
	defer conn.Close()

	// 2. Create your service instance
	service := NewThirdPartyService(conn)

	// 3. Process single document
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docInfo, err := service.ProcessDocument(ctx, "Hello, this is a test document in English.", "doc-001")
	if err != nil {
		log.Fatalf("Failed to process document: %v", err)
	}

	log.Printf("Processed Document:")
	log.Printf("  ID: %s", docInfo.ID)
	log.Printf("  Language: %s (%.2f confidence)", docInfo.Language, docInfo.Confidence)
	log.Printf("  Processing Time: %v", docInfo.ProcessingTime)
	log.Printf("  Service Version: %s", docInfo.ServiceVersion)
	
	if len(docInfo.Alternatives) > 0 {
		log.Printf("  Alternatives:")
		for i, alt := range docInfo.Alternatives {
			if i >= 3 { // Show only top 3
				break
			}
			log.Printf("    %d. %s (%.2f)", i+1, alt.Language, alt.Confidence)
		}
	}

	// 4. Process multiple documents
	documents := []Document{
		{ID: "doc-002", Content: "Hola, este es un documento en español."},
		{ID: "doc-003", Content: "Bonjour, ceci est un document en français."},
		{ID: "doc-004", Content: "Hallo, dies ist ein Dokument auf Deutsch."},
		{ID: "doc-005", Content: "Ciao, questo è un documento in italiano."},
	}

	log.Println("\nProcessing batch of documents...")
	batchResults, err := service.BatchProcessDocuments(ctx, documents)
	if err != nil {
		log.Printf("Batch processing failed: %v", err)
	}

	for _, result := range batchResults {
		log.Printf("  %s: %s (%.2f)", result.ID, result.Language, result.Confidence)
	}
}
