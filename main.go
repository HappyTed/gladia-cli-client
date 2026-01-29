// audio транскрибация через gladio.ai
// https://docs.gladia.io/chapters/pre-recorded-stt/quickstart#polling
package main

import (
	"context"
	"fmt"
	"os"

	"go-gladia.io-client/internal/cmd"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/logic/client"
	"go-gladia.io-client/internal/logic/usecase"
	"go-gladia.io-client/internal/repo"
	"go-gladia.io-client/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(
		cfg.LogLevel,
	)

	// repo

	fileRepo, err := repo.NewFilesRepo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// http client
	gaClient, err := client.NewGladiaClient(
		client.WithLogger(log),
		client.WithApiToken(cfg.Token),
		client.WithBaseUrl(cfg.BaseUrl),
		client.WithTimeout(cfg.RequestTimeout),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	uc, err := usecase.New(
		usecase.WithLogger(log),
		usecase.WithFileRepo(*fileRepo),
		usecase.WithHttpClient(gaClient),
	)

	ctx := context.TODO()
	cmd.Execute(ctx, cfg, log, uc)
}
