package audio

import (
	"context"
	"time"

	"go-gladia.io-client/internal/clients/http/models/prerecorderv2"
	"go-gladia.io-client/internal/config"
)

type (
	AudioAwait interface {
		// Загрузить файл для транскрибации
		Upload(filePath string) (string, error)
		// Запустить задачу на транскрибацию
		InitTranscription(cfg config.Config, audioURL string) (string, string, error)
		// Информация о статусе задачи
		Info(taskID string) (*prerecorderv2.Result, error)
		// Ожидать результат
		PollingResult(ctx context.Context, taskID string, timeInterval time.Duration, timeout time.Duration) (*prerecorderv2.Result, error)
		// Список загруженных на сервер задач
		List() error
		// Сдампить результат в файл
		Dump() error
	}

	AudioRecorder interface {
	}
)
