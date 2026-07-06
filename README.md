# go-youtube-transcript-api

A Go implementation of the popular [youtube-transcript-api](https://github.com/jdepoix/youtube-transcript-api) Python project. It allows you to fetch YouTube transcripts (captions) without needing the YouTube Data API or any browser automation.

## Features

- Fetches transcripts directly from YouTube's internal API (InnerTube).
- Supports both manually created and auto-generated transcripts.
- Translates between Python idioms and appropriate Go idioms.
- Provides a clean library interface and a powerful CLI.

## Installation

```bash
go get github.com/virmundi/go-youtube-transcript-api
```

## Usage

### As a Command Line Tool

Install the CLI:

```bash
go install github.com/virmundi/go-youtube-transcript-api/cmd/yt-transcript@latest
```

Or build the binary and run it:

```bash
go build -o yt-transcript ./cmd/yt-transcript
./yt-transcript <VIDEO_ID>
```

**Options:**
- `--language`: Specify the language code (default: `en`).

Example:
```bash
./yt-transcript J0lVsnlEtyM --language en
```

### As a Library

Import the package into your Go project:

```go
package main

import (
	"fmt"
	"log"

	"github.com/virmundi/go-youtube-transcript-api"
)

func main() {
	videoID := "J0lVsnlEtyM"
	languageCode := "en"

	// Simplified fetch (Transcript only)
	transcript, err := youtubetranscript.FetchTranscript(videoID, languageCode)
	if err != nil {
		log.Fatalf("Error fetching transcript: %v", err)
	}

	for _, snippet := range transcript.Snippets {
		fmt.Printf("[%0.2fs] %s\n", snippet.Start, snippet.Text)
	}

	// Fetch with Metadata (Channel details, Keywords, Description)
	result, err := youtubetranscript.FetchTranscriptWithMetadata(videoID, languageCode)
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
go run ./cmd/yt-transcript J0lVsnlEtyM --language en
```

## Running Tests

To verify the implementation:

```bash
go test -v ./...
```

## Project Structure

- `cmd/`: CLI implementation.
- Root directory (`client.go`, `transcript.go`, etc.): Core library logic.

## License

MIT
