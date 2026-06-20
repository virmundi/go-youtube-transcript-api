package lib

type TranscriptSnippet struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}

type Transcript struct {
	VideoID      string              `json:"video_id"`
	Language     string              `json:"language"`
	LanguageCode string              `json:"language_code"`
	IsGenerated  bool                `json:"is_generated"`
	Snippets     []TranscriptSnippet `json:"snippets"`
}
