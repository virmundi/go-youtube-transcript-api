# go-youtube-transcript-api

A Go implementation of the popular [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api) Python project. It allows you to fetch YouTube transcripts (captions) without needing the YouTube Data API or any browser automation.

## Features

- Fetches transcripts directly from YouTube's internal API (InnerTube).
- Supports both manually created and auto-generated transcripts.
- Translates between Python idioms and appropriate Go idioms.
- Provides a clean library interface and a powerful CLI.

## Installation

```bash
go get go-youtube-transcript-api
```

## Usage

### As a Command Line Tool

Build the binary and run it with a video ID:

```bash
go build -o yt-transcript .
./yt-transcript <VIDEO_ID>
```

**Options:**
- `--language`: Specify the language code (default: `en`).

Example:
```bash
./yt-transcript J0lVsnlEtyM --language en
```

### As a Library

Import the `lib` package into your Go project:

```go
package main

import (
	"fmt"
	"log"

	"go-youtube-transcript-api/lib"
)

func main() {
	videoID := "J0lVsnlEtyM"
	languageCode := "en"

	// Simplified fetch (Transcript only)
	transcript, err := lib.FetchTranscript(videoID, languageCode)
	if err != nil {
		log.Fatalf("Error fetching transcript: %v", err)
	}

	for _, snippet := range transcript.Snippets {
		fmt.Printf("[%0.2fs] %s\n", snippet.Start, snippet.Text)
	}

	// Fetch with Metadata (Channel details, Keywords, Description)
	result, err := lib.FetchTranscriptWithMetadata(videoID, languageCode)
	if err != nil {
		log.Fatalf("Error fetching with metadata: %v", err)
	}

	fmt.Printf("Channel: %s\n", result.Metadata.ChannelName)
	fmt.Printf("Description: %s\n", result.Metadata.ShortDescription)
}
```

### As a Go Script (One-liner)

If you have Go installed, you can run it directly:

```bash
go run . J0lVsnlEtyM --language en
```

## Running Tests

To verify the implementation against baseline Python data:

```bash
go test -v ./lib/...
```

## Project Structure

- `lib/`: Core library logic for fetching and parsing transcripts.
- `cmd/`: CLI implementation using `cobra`.
- `main.go`: Application entry point.

## License

MIT
