package usecase

import (
	"fmt"
	"path/filepath"

	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/logic/client"
	"go-gladia.io-client/internal/repo"
	"go-gladia.io-client/pkg/logger"
	"go-gladia.io-client/pkg/networking"
)

type FileType string

const (
	JSON = FileType(".json")
	TXT  = FileType(".txt")
)

type Usecase struct {
	gaClient client.GladiaClient
	repo     repo.FilesRepo
	log      logger.ILogger
	result   *prerecorderv2.Result
}

func (uc *Usecase) Upload(filePath string) error {
	// open audio file
	uc.log.DebugF("Try read file from path: %s\n", filePath)

	err, file, close := uc.repo.Open(filePath)
	if err != nil {
		return err
	}
	defer close()

	uc.log.Debug("Open file done")

	if resp, err := uc.gaClient.AudioUploadFromFile(file); err != nil {
		return err
	} else {
		fmt.Println("Audio Url:", resp.AudioUrl)
		metaData, err := networking.JsonDumpS(resp.MetaData)
		fmt.Printf("Meta Data:", metaData)
		return err
	}
}

func (uc *Usecase) Dump(path string, name string, fileType FileType) error {
	output := filepath.Join(path, fmt.Sprint(name, fileType))
	if fileType == JSON {
	} else if fileType == TXT {

	} else {
		return fmt.Errorf("failed to dump data: unknown output type: %s", fileType)
	}
}
