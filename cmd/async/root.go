package async

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/logic"
	"go-gladia.io-client/pkg/logger"
)

var rootCmd = &cobra.Command{
	Use:     "app",
	Version: "v1.0.0",
	Short:   "Demo gladia.io api cli-client",
	Args:    cobra.ExactArgs(1),
	Long:    ``,
}

func Execute(
	ctx context.Context,
	cfg *config.Config,
	log logger.ILogger,
	a_uc logic.IAsyncUsecase,
	s_uc logic.ISyncUsecase,
) error {

	// flags set
	// rootCmd.PersistentFlags().StringVarP(&cfg.Token, "token", "a", cfg.Token, "gladia api token")
	setUploadFlags()
	setTranscriptionFlags(cfg)

	// set usaceses

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// filePath := args[0]
		// _, err := os.Stat(filePath)
		// if err != nil {
		// 	return errors.New("file does not exist:" + filePath)
		// }

		// audioURL, err := uc.Upload(filePath)
		// if err != nil {
		// 	return err
		// }

		// _, _, err := uc.InitTranscription(*cfg, audioURL)

		return nil
	}

	uploadCmd.RunE = func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		_, err := os.Stat(filePath)
		if err != nil {
			return errors.New("file does not exist:" + filePath)
		}

		audioURL, err := uc.Upload(filePath)
		if err != nil {
			return err
		}
		log.Info("Audio Url: ", audioURL)
		return nil
	}

	transcriptionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		audioURL := args[0]
		resultURL, taskID, err := uc.InitTranscription(*cfg, audioURL)
		if err != nil {
			return err
		}

		fmt.Println("Result Url:", resultURL)
		fmt.Println("Task ID:", taskID)

		return nil
	}

	infoCmd.RunE = func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		result, err := uc.Info(taskID)
		if err != nil {
			return err
		}
		if result != nil {
			fmt.Println(*result)
		}
		return nil
	}

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(transcriptionCmd)
	rootCmd.AddCommand(infoCmd)

	cobra.CheckErr(rootCmd.Execute())

	return nil
}
