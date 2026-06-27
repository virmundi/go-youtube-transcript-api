package lib

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestFetchTranscript_Success(t *testing.T) {
	videoID := "J0lVsnlEtyM"
	languageCode := "en"

	transcript, err := FetchTranscript(videoID, languageCode)
	if err != nil {
		t.Fatalf("Failed to fetch transcript: %v", err)
	}

	if transcript.VideoID != videoID {
		t.Errorf("Expected video ID %s, got %s", videoID, transcript.VideoID)
	}

	if len(transcript.Snippets) == 0 {
		t.Fatal("Expected snippets, got 0")
	}

	// Verify the first few snippets match expected content from Python
	expectedFirst := "I've spent the last 20 years of my life"
	if transcript.Snippets[0].Text != expectedFirst {
		t.Errorf("Expected first snippet '%s', got '%s'", expectedFirst, transcript.Snippets[0].Text)
	}

	expectedSecond := "immersed in money. First through"
	if transcript.Snippets[1].Text != expectedSecond {
		t.Errorf("Expected second snippet '%s', got '%s'", expectedSecond, transcript.Snippets[1].Text)
	}

	// Write to a file for manual inspection or as requested
	outputFile := "test_transcript_success.json"
	data, err := json.MarshalIndent(transcript, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal transcript: %v", err)
	}

	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write transcript to file: %v", err)
	}
	t.Logf("Transcript written to %s", outputFile)
}

func TestFetchTranscript_InvalidVideo(t *testing.T) {
	videoID := "EnSJN9zl-yQ" // Bogus video ID
	languageCode := "en"

	_, err := FetchTranscript(videoID, languageCode)
	if err == nil {
		t.Fatal("Expected error for invalid video ID, but got nil")
	}

	// In the current implementation, a 404 or missing API key results in specific errors
	t.Logf("Received expected error: %v", err)
}

func TestFetchTranscriptWithMetadata_Success(t *testing.T) {
	videoID := "J0lVsnlEtyM"
	languageCode := "en"

	result, err := FetchTranscriptWithMetadata(videoID, languageCode)
	if err != nil {
		t.Fatalf("Failed to fetch transcript with metadata: %v", err)
	}

	if result.Transcript == nil {
		t.Fatal("Expected transcript, got nil")
	}

	if result.Transcript.VideoID != videoID {
		t.Errorf("Expected video ID %s, got %s", videoID, result.Transcript.VideoID)
	}

	metadata := result.Metadata
	expectedChannelName := "Felix & Friends (Goat Academy)"
	if metadata.ChannelName != expectedChannelName {
		t.Errorf("Expected channel name '%s', got '%s'", expectedChannelName, metadata.ChannelName)
	}

	expectedChannelID := "UCJtfma0mE_XrBAD9uakcjfA"
	if metadata.ChannelID != expectedChannelID {
		t.Errorf("Expected channel ID '%s', got '%s'", expectedChannelID, metadata.ChannelID)
	}

	if len(metadata.Keywords) == 0 {
		t.Errorf("Expected keywords, got 0")
	}

	expectedFirstKeyword := "felix & friends"
	if metadata.Keywords[0] != expectedFirstKeyword {
		t.Errorf("Expected first keyword '%s', got '%s'", expectedFirstKeyword, metadata.Keywords[0])
	}

	if !strings.Contains(metadata.ShortDescription, "In this video I explain the most common mistakes investors make") {
		t.Errorf("Short description does not contain expected text. Got: %s", metadata.ShortDescription)
	}
	if metadata.Title != "9 Investing Habits Keeping You Poor [BANKER EXPLAINS]" {
		t.Errorf("Title does not match expected. Got: %s", metadata.Title)
	}
}
