package lib

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	watchURL = "https://www.youtube.com/watch?v="
)

// FetchTranscript fetches the transcript for a given video ID and language.
func FetchTranscript(videoID, languageCode string) (*Transcript, error) {
	client := NewClient()

	// 1. Fetch the video watch page
	resp, err := client.Get(watchURL + videoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get video page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrVideoUnavailable
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read video page body: %w", err)
	}

	// 2. Extract InnerTube API key
	apiKey, err := extractAPIKey(string(body))
	if err != nil {
		return nil, err
	}

	// 3. Call the InnerTube API to get transcript list
	transcriptList, err := fetchTranscriptList(client, videoID, apiKey)
	if err != nil {
		return nil, err
	}

	// 4. Find the correct transcript URL
	transcriptURL, isGenerated, err := findTranscriptURL(transcriptList, languageCode)
	if err != nil {
		return nil, err
	}

	// 5. Fetch and parse the XML transcript
	return fetchAndParseXML(client, videoID, languageCode, transcriptURL, isGenerated)
}

func extractAPIKey(html string) (string, error) {
	re := regexp.MustCompile(`"INNERTUBE_API_KEY":"([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", ErrTranscriptNotFound
	}
	return matches[1], nil
}

func fetchTranscriptList(client *Client, videoID, apiKey string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://www.youtube.com/youtubei/v1/player?key=%s", apiKey)

	body := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "ANDROID",
				"clientVersion": "20.10.38",
			},
		},
		"videoId": videoID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	captions, ok := data["captions"].(map[string]interface{})
	if !ok {
		return nil, ErrTranscriptNotFound
	}

	renderer, ok := captions["playerCaptionsTracklistRenderer"].(map[string]interface{})
	if !ok {
		return nil, ErrTranscriptNotFound
	}

	return renderer, nil
}

func findTranscriptURL(transcriptList map[string]interface{}, languageCode string) (string, bool, error) {
	tracks, ok := transcriptList["captionTracks"].([]interface{})
	if !ok {
		return "", false, ErrTranscriptNotFound
	}

	for _, t := range tracks {
		track, ok := t.(map[string]interface{})
		if !ok {
			continue
		}

		if track["languageCode"] == languageCode {
			url, ok := track["baseUrl"].(string)
			if !ok {
				continue
			}
			isGenerated := track["kind"] == "asr"
			return url, isGenerated, nil
		}
	}

	return "", false, ErrTranscriptNotFound
}

type xmlSnippet struct {
	Text     string  `xml:",innerxml"`
	Start    float64 `xml:"t,attr"`
	Duration float64 `xml:"d,attr"`
}

type xmlTranscript struct {
	XMLName  xml.Name     `xml:"timedtext"`
	Snippets []xmlSnippet `xml:"body>p"`
}

func fetchAndParseXML(client *Client, videoID, languageCode, url string, isGenerated bool) (*Transcript, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch xml: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch xml, status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read xml body: %w", err)
	}

	var xt xmlTranscript
	if err := xml.Unmarshal(body, &xt); err != nil {
		return nil, fmt.Errorf("failed to decode xml: %w", err)
	}

	snippets := make([]TranscriptSnippet, 0, len(xt.Snippets))
	for _, s := range xt.Snippets {
		// Clean the text from <s> tags and other HTML-like entities
		text := s.Text
		text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")
		text = html.UnescapeString(text)
		text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		snippets = append(snippets, TranscriptSnippet{
			Text:     text,
			Start:    s.Start / 1000.0,
			Duration: s.Duration / 1000.0,
		})
	}

	return &Transcript{
		VideoID:      videoID,
		LanguageCode: languageCode,
		IsGenerated:  isGenerated,
		Snippets:     snippets,
	}, nil
}
