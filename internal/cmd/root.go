package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/logic"
	"go-gladia.io-client/pkg/logger"
	"go-gladia.io-client/pkg/networking"
)

var (
	token  string
	audio  string
	await  bool
	output string
)

var rootCmd = &cobra.Command{
	Use:     "app",
	Version: "v1.0.0",
	Short:   "Demo gladia.io api cli-client",
	Long:    ``,
}

func Execute(
	ctx context.Context,
	cfg config.Config,
	log logger.ILogger,
	uc logic.IUsecase,
	client logic.IGladiaClient,
) error {

	// flags set
	rootCmd.PersistentFlags().StringVarP(&cfg.AudioFile, "audio", "a", "", "path to audio file")
	rootCmd.MarkPersistentFlagRequired("audio")
	rootCmd.PersistentFlags().StringVarP(&cfg.Token, "token", "a", cfg.Token, "gladia api token")
	setUploadFlags()
	setTranscriptionFlags(cfg)

	// set usaceses
	uploadCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if resp, err := gaClient.AudioUploadFromFile(file); err != nil {
			return err
		} else {
			fmt.Println("Audio Url:", resp.AudioUrl)
			metaData, err := networking.JsonDumpS(resp.MetaData)
			fmt.Printf("Meta Data:", metaData)
			return err
		}
	}

	transcriptionCmd.RunE = func(cmd *cobra.Command, args []string) error {
		uploadResp, err := gaClient.AudioUploadFromFile(file)
		if err != nil {
			return err
		}

		body := &prerecorderv2.PreRecorderBody{
			AudioUrl:    uploadResp.AudioUrl,
			Diarization: cfg.Diarization,
			DiarizationConf: &prerecorderv2.DiarizationConf{
				Enhanced:      cfg.Enhanced,
				MinSpeakers:   cfg.MinSpeakers,
				MaxSpeakers:   cfg.MaxSpeakers,
				NumOfSpeakers: cfg.Speakers,
			},
			Translation: cfg.Translation,
			TranslationConf: &prerecorderv2.TranslationConf{
				TargetLanguages: cfg.TranslationLanguages,
			},
			SentimentAnalysis: cfg.SentimentAnalysis,
		}

		initResp, err := gaClient.InitTranscription(body)
		if err != nil {
			return err
		}

		if cfg.AwaitResults {
			resp, err := gaClient.AwaitTranscriptionResult(ctx, initResp.ID, cfg.AwaitInterval)
		} else {
			resp, err := gaClient.TranscriptionResult(initResp.ID, cfg.AwaitInterval)
		}
	}

	cobra.CheckErr(rootCmd.Execute())
}
