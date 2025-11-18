package cmd

import (
	"github.com/spf13/cobra"
	"go-gladia.io-client/internal/config"
)

var transcriptionCmd = &cobra.Command{
	Use:   "transcription",
	Short: "Transcription audio",
}

func setTranscriptionFlags(cfg *config.Config) {
	transcriptionCmd.Flags().BoolVarP(&cfg.AwaitResults, "await", "a", false, "wait for the transcription to finish")
	transcriptionCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "./result.txt", "name and path of the file for recording the transcription")
}
