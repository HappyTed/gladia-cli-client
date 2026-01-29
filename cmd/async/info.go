package async

import "github.com/spf13/cobra"

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Find out the status of an audio transcription task.",
}

func setListFlags() {
	// пока флагов нет
}
