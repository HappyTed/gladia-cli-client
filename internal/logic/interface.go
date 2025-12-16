package logic

import (
	"context"
	"os"
	"time"

	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/entities/upload"
)

type (
	IGladiaClient interface {
		AudioUploadFromFile(file *os.File) (*upload.UploadResponce, error)
		InitTranscription(body *prerecorderv2.PreRecorderBody) (*prerecorderv2.PreRecorderInitResponse, error)
		TranscriptionResult(jobId string, timeInterval time.Duration) (*prerecorderv2.PreRecorderResultResponse, error)
		AwaitTranscriptionResult(ctx context.Context, jobId string, timeInterval time.Duration) (*prerecorderv2.PreRecorderResultResponse, error)
	}

	IUsecase interface {
		Upload(filePath string) error
		Transcription() error
		Info() error
		List() error
		Dump() error
	}
)
