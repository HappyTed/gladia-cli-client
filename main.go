// audio транскрибация через gladio.ai
// https://docs.gladia.io/chapters/pre-recorded-stt/quickstart#polling
package main

import (
	"context"
	"os"

	"go-gladia.io-client/internal/cmd"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/logic/client"
	"go-gladia.io-client/internal/logic/usecase"
	"go-gladia.io-client/internal/repo"
	"go-gladia.io-client/internal/repo/database"
	"go-gladia.io-client/pkg/logger"
)

func init() {
	/* TODO:
	1) mkdir ~/.config/myapp/ ~/.local/share/myapp/
	2) dump config ~/.config/myapp/config.yml
	3) logs: ~/.local/share/myapp/myapp.log
	4) data: ~/.local/share/myapp/myapp.db
	*/
}

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(
		cfg.LogLevel,
		logger.FILE|logger.STD,
		cfg.LogPath,
	)

	// repo
	database, err := database.New("sqlite3", "/tmp/test.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fileRepo, err := repo.NewFilesRepo(log)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		os.Exit(1)
	}

	uc, err := usecase.New(
		usecase.WithLogger(log),
		usecase.WithDatabase(database),
		usecase.WithFileRepo(*fileRepo),
		usecase.WithHttpClient(gaClient),
	)

	ctx := context.TODO()
	cmd.Execute(ctx, cfg, log, uc)
}
