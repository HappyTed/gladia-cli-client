package usecase

import (
	"errors"
	"fmt"
	"net/url"

	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/internal/entities"
	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/logic"
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
	log  logger.ILogger
	c    logic.IGladiaClient
	repo repo.FilesRepo
	db   database.IDatabase
}

type Option func(*Usecase) error

func WithLogger(l logger.ILogger) Option {
	return func(u *Usecase) error {
		if l == nil {
			return errors.New("failed to init usecase: logger is empty")
		}
		u.log = l
		return nil
	}
}

func WithHttpClient(c logic.IGladiaClient) Option {
	return func(u *Usecase) error {
		if c == nil {
			return errors.New("failed to init usecase: http client is empty")
		}
		u.c = c
		return nil
	}
}

func WithFileRepo(f repo.FilesRepo) Option {
	return func(u *Usecase) error {
		u.repo = f
		return nil
	}
}

func WithDatabase(d database.IDatabase) Option {
	return func(u *Usecase) error {
		if d == nil {
			return errors.New("failed to init usecase: database is empty")
		}
		u.db = d
		return nil
	}
}

func New(options ...Option) (*Usecase, error) {
	r := &Usecase{}
	for _, opt := range options {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Загрузить аудио файл на сервер gladia
func (uc *Usecase) Upload(filePath string) (int64, error) {
	// открыть audio file
	uc.log.DebugF("Try read file from path: %s\n", filePath)

	err, file, close := uc.repo.Open(filePath)
	if err != nil {
		uc.log.ErrorF("%s: file read error: %w", filePath, err)
		return -1, fmt.Errorf("%s: file read error: %w", filePath, err)
	}
	defer close()

	uc.log.Debug("Open file done")

	// загрузить файл
	dbRecord := database.Task{
		FilePath: filePath,
	}

	resp, err := uc.c.AudioUploadFromFile(file)

	// обработка ответа от сервера
	if err != nil {
		dbRecord.Results = resp.Dump()
		dbRecord.Status = string(entities.ERROR)
		uc.log.Error(err)
		return -1, err
	}

	dbRecord.FileURL = resp.AudioUrl
	uc.log.Debug("Audio Url:", resp.AudioUrl)

	metaData, err := networking.JsonDumpS(resp.MetaData)
	uc.log.DebugF("Meta Data: %s", metaData)

	id, err := uc.db.Insert(dbRecord)
	uc.log.Debug("The entry was successfully uploaded. Internal ID:", id)

	return id, err
}

// Выполнить асинхронный запрос к сервису на транскрибацию и получить task_id
func (uc *Usecase) Transcription(cfg config.Config, id int64) (string, error) {
	var url *url.URL

	dbRecord, err := uc.db.GetByID(id)

	if err != nil {
		uc.log.Error("error: init transcription: file not upload:", err)
		return "", err
	} else if url, err = url.Parse(dbRecord.FileURL); err != nil {
		uc.log.Error("error: init transcription: gladia file url is not valid:", err)
		return "", err
	}

	body := &prerecorderv2.PreRecorderBody{
		AudioUrl:    url.String(),
		Diarization: cfg.Diarization,
		LangConf: &prerecorderv2.LanguageConf{
			Languages:     cfg.InputLanguages,
			CodeSwitching: false,
		},
		Translation: cfg.Translation,
		TranslationConf: &prerecorderv2.TranslationConf{
			TargetLanguages: cfg.TargetLanguages,
		},
		Subtitle:          false,
		SubtitlesConf:     &prerecorderv2.SubtitlesConf{},
		SentimentAnalysis: true,
	}

	resp, err := uc.c.InitTranscription(body)
	if err != nil {
		uc.log.Error("error: init transcription: ", err)
		return "", err
	}

	dbRecord.TaskID = resp.ID
	dbRecord.ResultURL = resp.ResultUrl

	return "", err
}

// Выгрузить расшифровку аудио в файл
func (uc *Usecase) Dump() error {
	return nil
	// output := filepath.Join(path, fmt.Sprint(name, fileType))

	// fmt.Println(output)

	// if fileType == JSON {
	// } else if fileType == TXT {

	// } else {
	// 	return fmt.Errorf("failed to dump data: unknown output type: %s", fileType)
	// }

	// return nil
}

// Список задач, загруженных через клиент
func (uc *Usecase) List() ([]database.Task, error) {
	list, err := uc.db.GetAll()
	if err != nil {
		uc.log.Error("failed get list data:", err)
		return nil, err
	}

	uc.log.DebugF("List:\n%+v", list)
	return nil, err
}

func (uc *Usecase) PollingResult(id int64) error { return nil }

func (uc *Usecase) Info(id int64) (database.Task, error) { return database.Task{}, nil }
