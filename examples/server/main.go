package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"github.com/hovman/ld_proto/pb"
)

// Simple language detection based on common words
func detectLanguage(text string) (string, float32, []*pb.LanguageAlternative) {
	text = strings.ToLower(text)
	
	// Simple word-based detection
	languageWords := map[string][]string{
		"en": {"the", "and", "is", "in", "to", "of", "a", "that", "it", "with"},
		"es": {"el", "la", "de", "que", "y", "a", "en", "un", "es", "se"},
		"fr": {"le", "de", "et", "à", "un", "il", "être", "et", "en", "avoir"},
		"de": {"der", "die", "und", "in", "den", "von", "zu", "das", "mit", "sich"},
		"it": {"il", "di", "che", "e", "la", "per", "un", "in", "con", "da"},
	}
	
	scores := make(map[string]int)
	totalWords := 0
	
	words := strings.Fields(text)
	for _, word := range words {
		totalWords++
		for lang, langWords := range languageWords {
			for _, langWord := range langWords {
				if word == langWord {
					scores[lang]++
				}
			}
		}
	}
	
	if totalWords == 0 {
		return "unknown", 0.0, nil
	}
	
	// Find best match
	bestLang := "en" // default
	bestScore := 0
	for lang, score := range scores {
		if score > bestScore {
			bestScore = score
			bestLang = lang
		}
	}
	
	confidence := float32(bestScore) / float32(totalWords)
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	// Create alternatives
	var alternatives []*pb.LanguageAlternative
	for lang, score := range scores {
		if lang != bestLang && score > 0 {
			altConfidence := float32(score) / float32(totalWords)
			if altConfidence > 1.0 {
				altConfidence = 1.0
			}
			alternatives = append(alternatives, &pb.LanguageAlternative{
				LanguageCode: lang,
				Confidence:   altConfidence,
			})
		}
	}
	
	return bestLang, confidence, alternatives
}

type server struct {
	pb.UnimplementedLanguageDetectionServiceServer
}

func (s *server) DetectLanguage(ctx context.Context, req *pb.DetectLanguageRequest) (*pb.DetectLanguageResponse, error) {
	start := time.Now()
	
	// Perform language detection
	languageCode, confidence, alternatives := detectLanguage(req.Text)
	
	processingTime := time.Since(start)
	
	return &pb.DetectLanguageResponse{
		LanguageCode: languageCode,
		Confidence:   confidence,
		Alternatives: alternatives,
		DocumentId:   req.DocumentId,
		Metadata: &pb.ProcessingMetadata{
			ProcessingTimeMs: processingTime.Milliseconds(),
			ServiceVersion:   "1.0.0",
			ModelVersion:     "simple-word-based-v1.0",
			Provider:         "ld_proto_example",
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

	log.Println("Language Detection Server starting on :50051")
	log.Println("Use the client example to test the service")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
