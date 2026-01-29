package audio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	http_client "go-gladia.io-client/internal/clients/http"
	"go-gladia.io-client/internal/clients/http/models/prerecorderv2"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/pkg/output"
)

type FileType string

const (
	JSON = FileType(".json")
	TXT  = FileType(".txt")
)

type AudoUploader struct {
	l          output.IOutput
	httpClient http_client.IHttpClient
}

func New(l output.IOutput, client http_client.IHttpClient) (*AudoUploader, error) {
	r := &AudoUploader{
		l:          l,
		httpClient: client,
	}

	return r, nil
}

// Загрузить аудио файл на сервер gladia
func (uc *AudoUploader) Upload(filePath string) error {
	// открыть audio file
	uc.l.Printf("Try read file from path: %s\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		uc.l.Printf("file read error %s: %w\n", filePath, err)
		return fmt.Errorf("%s: file read error %w", filePath, err)
	}
	defer file.Close()

	uc.l.Print("Open file done")

	// загрузить файл
	resp, err := uc.httpClient.AudioUploadFromFile(file)

	// обработка ответа от сервера
	if err != nil {
		uc.l.Print("upload error:", err)
		return err
	}

	audioURL := resp.AudioUrl

	uc.l.Print("File upload done!")
	uc.l.Print("Audo Url:", audioURL)

	metaData, err := json.Marshal(resp.MetaData)
	if err != nil {
		return err
	}
	uc.l.Printf("Meta Data: %s\n", metaData)

	return nil
}

// Выполнить асинхронный запрос к сервису на транскрибацию и получить task_id
func (uc *AudoUploader) InitTranscription(cfg config.Config, audioURL string) (string, string, error) {
	var url *url.URL
	var err error

	if url, err = url.Parse(audioURL); err != nil {
		uc.l.Print("error: init transcription: gladia file url is not valid:", err)
		return "", "", err
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

	resp, err := uc.httpClient.InitTranscription(body)
	if err != nil {
		uc.l.Print("Failed init transcription: ", err)
		return "", "", err
	}

	return resp.ResultUrl, resp.ID, err
}

func (uc *AudoUploader) PollingResult(ctx context.Context, taskID string, timeInterval time.Duration, timeout time.Duration) (*prerecorderv2.Result, error) {

	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	var resp *prerecorderv2.PreRecorderResultResponse
	var err error

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout waiting result")
		case <-ticker.C:
			resp, err = uc.httpClient.GetTranscriptionResult(taskID)
			if err != nil {
				return nil, err
			} else if resp.Status == "error" {
				return nil, fmt.Errorf("error: %s", resp.ErrorCode)
			} else if resp.Status == "done" {
				return resp.Result, nil
			}
		case <-time.After(timeout):
			return nil, fmt.Errorf("the result was not obtained within: %d", timeout)
		}
	}
}

func (uc *AudoUploader) Info(taskID string) (*prerecorderv2.Result, error) {
	resp, err := uc.httpClient.GetTranscriptionResult(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed get task result: %w", err)
	}

	uc.l.Print("Task status: ", resp.Status)

	if resp.Status == "error" {
		uc.l.Print("errors:")
		uc.l.Print("\terror code:", resp.ErrorCode)
	} else if resp.Status == "done" {
		return resp.Result, nil
	}

	return nil, nil
}

func (uc *AudoUploader) List(limit int) (string, error) {
	resp, err := uc.httpClient.List(limit)
	if err != nil {
		return "", fmt.Errorf("failed get tasks listt: %w", err)
	}

	type result struct {
		id              int
		uid             string
		status          string
		transactionType string
		fileName        string
		Duration        string
		Date            string
	}

	var results []result

	for idx, item := range resp.Items {
		res := result{
			id:              idx,
			uid:             item.ID,
			status:          item.Status,
			transactionType: "ХЗ ГДЕ БРАТЬ",
			fileName:        item.File.Filename,
			Duration:        "Посчитать разницу: completed_at - created_at",
			Date:            item.CompletedAT,
		}

		results = append(results, res)
	}

	return "", nil
}
