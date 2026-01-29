package http_client

import (
	"os"

	"go-gladia.io-client/internal/clients/http/models/prerecorderv2"
	"go-gladia.io-client/internal/clients/http/models/upload"
)

type IHttpClient interface {
	AudioUploadFromFile(file *os.File) (*upload.UploadResponce, error)
	InitTranscription(body *prerecorderv2.PreRecorderBody) (*prerecorderv2.PreRecorderInitResponse, error)
	GetTranscriptionResult(jobId string) (*prerecorderv2.PreRecorderResultResponse, error)
	DownloadAudioFileid(id string) error
	DeleteTranscription(id string) error
	List(limit int) (*prerecorderv2.ListResponse, error)
}
