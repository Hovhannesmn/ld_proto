package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/Hovhannesmn/ld_proto/pb"
)

func main() {
	// Connect to the language detection service
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewLanguageDetectionServiceClient(conn)

	// Test cases
	testCases := []struct {
		text     string
		expected string
	}{
		{"Hello, world!", "en"},
		{"Hola, mundo!", "es"},
		{"Bonjour le monde!", "fr"},
		{"Hallo Welt!", "de"},
		{"Ciao mondo!", "it"},
	}

	for _, tc := range testCases {
		// Create a request
		req := &pb.DetectLanguageRequest{
			Text:       tc.text,
			DocumentId: "example-doc-" + time.Now().Format("20060102-150405"),
			Metadata: map[string]string{
				"source":    "example_client",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}

		// Call the service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := client.DetectLanguage(ctx, req)
		if err != nil {
			log.Printf("Language detection failed for '%s': %v", tc.text, err)
			continue
		}

		log.Printf("Text: '%s'", tc.text)
		log.Printf("Detected language: %s (confidence: %.2f)", resp.LanguageCode, resp.Confidence)
		log.Printf("Document ID: %s", resp.DocumentId)
		
		if resp.Metadata != nil {
			log.Printf("Processing time: %dms", resp.Metadata.ProcessingTimeMs)
			log.Printf("Service version: %s", resp.Metadata.ServiceVersion)
		}
		
		if len(resp.Alternatives) > 0 {
			log.Printf("Alternatives:")
			for i, alt := range resp.Alternatives {
				if i >= 3 { // Show only top 3 alternatives
					break
				}
				log.Printf("  %d. %s (%.2f)", i+1, alt.LanguageCode, alt.Confidence)
			}
		}
		log.Println("---")
	}
}
