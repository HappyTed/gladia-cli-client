package usecase

import (
	"fmt"
	"path/filepath"

	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/entities"
	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/logic/client"
	"go-gladia.io-client/internal/repo"
	"go-gladia.io-client/internal/repo/database"
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
	db       database.IDatabase
	log      logger.ILogger
	result   *prerecorderv2.Result
}

// Загрузить аудио файл на сервер gladia
func (uc *Usecase) Upload(filePath string) (int64, error) {
	// open audio file
	uc.log.DebugF("Try read file from path: %s\n", filePath)

	err, file, close := uc.repo.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer close()

	uc.log.Debug("Open file done")

	dbRecord := database.Request{
		FilePath: filePath,
	}

	resp, err := uc.gaClient.AudioUploadFromFile(file)
	if err != nil {
		dbRecord.Results = resp.Dump()
		dbRecord.Status = string(entities.ERROR)

		id, err := uc.db.Insert(dbRecord)
		return id, err
	}

	dbRecord.FileUrl = resp.AudioUrl
	fmt.Println("Audio Url:", resp.AudioUrl)

	metaData, err := networking.JsonDumpS(resp.MetaData)
	fmt.Printf("Meta Data: %s", metaData)

	id, err := uc.db.Insert(dbRecord)
	return id, err
}

// Выполнить асинхронный запрос к сервису на перевод и получить task_id
func (uc *Usecase) Transcription(cfg config.Config, fileUrl string) (string, error) {
	body := &prerecorderv2.PreRecorderBody{
		AudioUrl:    fileUrl,
		Diarization: cfg.Diarization,
		LangConf: &prerecorderv2.LanguageConf{
			Languages:     cfg.Languages,
			CodeSwitching: false,
		},
		Translation: cfg.Translation,
		TranslationConf: ,
	}

	resp, err := uc.gaClient.InitTranscription(body)
	return "", nil
}

// Выгрузить расшифровку аудио в файл
func (uc *Usecase) Dump(path string, name string, fileType FileType) error {
	output := filepath.Join(path, fmt.Sprint(name, fileType))

	fmt.Println(output)

	if fileType == JSON {
	} else if fileType == TXT {

	} else {
		return fmt.Errorf("failed to dump data: unknown output type: %s", fileType)
	}

	return nil
}
