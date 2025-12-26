package logic

import (
	"os"

	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/entities/upload"
	"go-gladia.io-client/internal/repo/database"
)

type (
	IGladiaClient interface {
		AudioUploadFromFile(file *os.File) (*upload.UploadResponce, error)
		InitTranscription(body *prerecorderv2.PreRecorderBody) (*prerecorderv2.PreRecorderInitResponse, error)
		GetTranscriptionResult(jobId string) (*prerecorderv2.PreRecorderResultResponse, error)
	}

	IUsecase interface {
		// Загрузить файл для транскрибации
		Upload(filePath string) (int64, error)
		// Запустить задачу на транскрибацию
		Transcription(cfg config.Config, id int64) (string, error)
		// Информация о статусе задачи
		Info(id int64) (database.Task, error)
		// Ожидать результат
		PollingResult(id int64) error
		// Список задач, загруженных через клиент
		List() ([]database.Task, error)
		// Сдампить результат в файл
		Dump() error
	}
)
