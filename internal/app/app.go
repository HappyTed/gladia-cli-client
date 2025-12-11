package app

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go-gladia.io-client/internal/cmd"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/logic/client"
	"go-gladia.io-client/pkg/logger"
)

func RunWithoutCobra(audio string) {

	log := logger.NewLogger(
		logger.DEBUG,
		logger.FILE|logger.STD,
		"logs/app.log",
	)

	cfg := config.LoadConfig()

	log.Debug("Try init client...")

	gaClient, err := client.NewGladiaClient(
		client.WithLogger(log),
		client.WithApiToken(cfg.Token),
		client.WithBaseUrl(cfg.BaseUrl),
		client.WithTimeout(cfg.RequestTimeout),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Debug("GladiaClient init done!")

	// open audio
	log.DebugF("Try read file from path: %s\n", audio)

	file, err := os.Open(audio)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	log.Debug("Open file done")

	uploadResult, err := gaClient.AudioUploadFromFile(file)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Debug("Upload file done!")

	body := &prerecorderv2.PreRecorderBody{
		AudioUrl:    uploadResult.AudioUrl,
		Diarization: true,
		DiarizationConf: &prerecorderv2.DiarizationConf{
			Enhanced: true,
		},
	}
	initTranscription, err := gaClient.InitTranscription(body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Debug("Init transcription done!")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if cfg.AwaitTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.AwaitTimeout)
	}
	defer cancel()

	fullResults, err := gaClient.AwaitTranscriptionResult(ctx, initTranscription.ID, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Transcription done!")

	jsonBytes, err := json.MarshalIndent(fullResults, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("results.json", jsonBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(
		cfg.LogLevel,
		logger.FILE|logger.STD,
		cfg.LogPath,
	)

	gaClient, err := client.NewGladiaClient(
		client.WithLogger(log),
		client.WithApiToken(cfg.Token),
		client.WithBaseUrl(cfg.BaseUrl),
		client.WithTimeout(cfg.RequestTimeout),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if cfg.AwaitTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.AwaitTimeout)
	}
	defer cancel()

	cmd.Execute(ctx, log, cfg, gaClient)
}
