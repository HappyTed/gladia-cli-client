package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of published files with information about them.",
}

func setListFlags() {
	// пока флагов нет
}
