package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/logic"
	"go-gladia.io-client/pkg/logger"
)

var rootCmd = &cobra.Command{
	Use:     "app",
	Version: "v1.0.0",
	Short:   "Demo gladia.io api cli-client",
	Long:    ``,
}

func Execute(
	ctx context.Context,
	cfg *config.Config,
	log logger.ILogger,
	uc logic.IUsecase,
) error {

	// flags set
	rootCmd.PersistentFlags().StringVarP(&cfg.AudioFile, "audio", "a", "", "path to audio file")
	// rootCmd.MarkPersistentFlagRequired("audio")
	// rootCmd.PersistentFlags().StringVarP(&cfg.Token, "token", "a", cfg.Token, "gladia api token")
	setUploadFlags()
	setTranscriptionFlags(cfg)

	listCmd.RunE = func(cmd *cobra.Command, args []string) error {
		list, err := uc.List()
		if err != nil {
			return err
		}
		fmt.Println(list)
		return nil
	}

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// err := uc.Upload(cfg.AudioFile)
		// if err != nil {
		// 	return err
		// }
		return nil
	}

	// set usaceses
	uploadCmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := uc.Upload(cfg.AudioFile)
		if err != nil {
			return err
		}
		fmt.Println("Id:", id)
		return nil
	}

	transcriptionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	cobra.CheckErr(rootCmd.Execute())

	return nil
}
