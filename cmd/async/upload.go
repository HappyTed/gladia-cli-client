package async

import "github.com/spf13/cobra"

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload audio file on gladia serv",
}

func setUploadFlags() {
	// пока флагов нет
}
