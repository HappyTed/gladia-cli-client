package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go-gladia.io-client/internal/clients/http/models/prerecorderv2"
	"go-gladia.io-client/internal/clients/http/models/upload"
	"go-gladia.io-client/internal/config"
	"go-gladia.io-client/pkg/output"
)

type GladiaClient struct {
	l       output.IOutput
	token   string
	client  *http.Client
	timeout time.Duration
	baseURL string
}

func NewGladiaClient(cfg config.HTTPClientConfig, l output.IOutput, apiToken string, urlPath string) (*GladiaClient, error) {
	gc := &GladiaClient{
		l:       l,
		baseURL: urlPath,
		token:   apiToken,
		client:  &http.Client{Timeout: cfg.Timeout},
	}

	return gc, nil
}

/*
# Загрузить аудио файл на платформу для дальнейшего транскрибирования

	curl --request POST \
	  --url https://api.gladia.io/v2/upload \
	  --header 'Content-Type: multipart/form-data' \
	  --header 'x-gladia-key: <api-key>' \
	  --form audio='@example-file'
*/
func (gc *GladiaClient) AudioUploadFromFile(file *os.File) (*upload.UploadResponce, error) {
	path := "/v2/upload"

	body, multipartHeader, err := makeMultipartBody(file, "audio", file.Name())
	if err != nil {
		return nil, err
	}

	gc.l.Print("Try do request:", gc.baseURL+path)

	req, err := http.NewRequest("POST", gc.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", multipartHeader)
	req.Header.Add("x-gladia-key", gc.token)

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseBody upload.UploadResponce
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		gc.l.Print("response parsing error:", err)
	}
	gc.l.Print("Response:", string(jsonResponse))

	err = httpErrorParse(resp, 200)
	if err != nil {
		return nil, err
	}
	return &responseBody, err
}

/*
# Начните предварительно записанную работу по транскрипции. Используйте возвращенное idи GET /v2/pre-записанная/:id конечная точка для получения результатов.

	curl --request POST \
	  --url https://api.gladia.io/v2/pre-recorded \
	  --header 'Content-Type: application/json' \
	  --header 'x-gladia-key: <api-key>' \
	  --data '{...}'
*/
func (gc *GladiaClient) InitTranscription(body *prerecorderv2.PreRecorderBody) (*prerecorderv2.PreRecorderInitResponse, error) {
	path := "/v2/pre-recorded"
	method := "POST"

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	gc.l.Print("Try do request:", gc.baseURL+path)
	gc.l.Print("Body:", string(jsonBody))

	req, err := http.NewRequest(method, gc.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-gladia-key", gc.token)

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = httpErrorParse(resp, 201)
	if err != nil {
		return nil, err
	}

	var responseBody prerecorderv2.PreRecorderInitResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		gc.l.Print("response parsing error:", err)
	}
	gc.l.Print("Response:", string(jsonResponse))

	return &responseBody, err
}

/*
# Получите предварительно записанный статус, параметры и результат транскрипции.

	curl --request GET \
	  --url https://api.gladia.io/v2/pre-recorded/{id} \
	  --header 'x-gladia-key: <api-key>'
*/
func (gc *GladiaClient) GetTranscriptionResult(jobId string) (*prerecorderv2.PreRecorderResultResponse, error) {
	path := fmt.Sprintf("/v2/pre-recorded/%s", jobId)
	method := "GET"

	gc.l.Print("Try do request:", gc.baseURL+path)

	req, err := http.NewRequest(method, gc.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-gladia-key", gc.token)
	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = httpErrorParse(resp, 200)
	if err != nil {
		return nil, err
	}

	var responseBody prerecorderv2.PreRecorderResultResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		gc.l.Print("response parsing error:", err)
	}
	gc.l.Print("Response:", string(jsonResponse))

	return &responseBody, err
}

/*
# Скачать аудиофайл, используемый на предварительно записанной транскрипции.

	curl --request GET \
	  --url https://api.gladia.io/v2/pre-recorded/{id}/file \
	  --header 'x-gladia-key: <api-key>'
*/
func (gc *GladiaClient) DownloadAudioFile(id string) error {
	path := fmt.Sprintf("/v2/pre-recorded/%s/file", id)
	URL := gc.baseURL + path

	gc.l.Print("Try do request:", URL)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-gladia-key", gc.token)

	resp, err := gc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = httpErrorParse(resp, 202)
	if err != nil {
		return err

		var response map[string]any

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			gc.l.Print("response parsing error:", err)
		}
		gc.l.Print("Response:", string(jsonResponse))
	}

	return nil
}

/*
# Удалить предварительно записанную транскрипцию и все ее данные (аудио-файл, транскрипция).

	curl --request DELETE \
	  --url https://api.gladia.io/v2/pre-recorded/{id} \
	  --header 'x-gladia-key: <api-key>'
*/
func (gc *GladiaClient) DeleteTranscription(id string) error {
	path := fmt.Sprintf("/v2/pre-recorded/%s", id)
	URL := gc.baseURL + path

	gc.l.Print("Try do request:", URL)

	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("x-gladia-key", gc.token)

	resp, err := gc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = httpErrorParse(resp, 202)
	if err != nil {
		return err

		var response map[string]any

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			gc.l.Print("response parsing error:", err)
		}
		gc.l.Print("Response:", string(jsonResponse))
	}

	return nil
}

/*
# Получить список отправленных в обработку записей

	curl --request GET \
		--url 'https://api.gladia.io/v2/pre-recorded?limit=20' \
		--header 'x-gladia-key: <api-key>'
*/
func (gc *GladiaClient) List(limit int) (*prerecorderv2.ListResponse, error) {
	path := fmt.Sprintf("/v2/pre-recorded?limit=%d", limit)
	URL := gc.baseURL + path

	gc.l.Print("Try do request:", URL)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-gladia-key", gc.token)

	resp, err := gc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseBody prerecorderv2.ListResponse

	err = httpErrorParse(resp, 200)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		gc.l.Print("response parsing error:", err)
	}
	gc.l.Print("Response:", string(jsonResponse))

	if err != nil {
		return nil, err
	}
	return &responseBody, err
}
