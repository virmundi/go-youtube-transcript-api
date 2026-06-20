package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"go-youtube-transcript-api/lib"

	"github.com/spf13/cobra"
)

var (
	languageCode string
)

var rootCmd = &cobra.Command{
	Use:   "go-youtube-transcript-api [videoID]",
	Short: "A Go application to fetch YouTube transcripts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		videoID := args[0]
		fmt.Printf("Fetching transcript for video ID: %s in language: %s\n", videoID, languageCode)

		transcript, err := lib.FetchTranscript(videoID, languageCode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(transcript, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&languageCode, "language", "en", "Language of the transcript")
}
