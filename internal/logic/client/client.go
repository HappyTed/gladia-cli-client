package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"go-gladia.io-client/internal/entities/prerecorderv2"
	"go-gladia.io-client/internal/entities/upload"
	"go-gladia.io-client/pkg/logger"
	"go-gladia.io-client/pkg/networking"
)

type GladiaClient struct {
	log         logger.ILogger
	token       string
	client      *http.Client
	timeout     time.Duration
	baseUrlPath string
	result      *prerecorderv2.Result
}

type Setter func(*GladiaClient) error

func WithLogger(log logger.ILogger) Setter {
	return func(gc *GladiaClient) error {
		if log == nil {
			return errors.New("failed to set nil logger")
		}
		gc.log = log
		return nil
	}
}

func WithBaseUrl(b string) Setter {
	return func(gc *GladiaClient) error {
		url, err := url.Parse(b)
		if err != nil {
			return err
		}
		gc.baseUrlPath = url.String()
		return nil
	}
}

func WithApiToken(t string) Setter {
	return func(gc *GladiaClient) error {
		gc.token = t
		return nil
	}
}

func WithTimeout(t time.Duration) Setter {
	return func(gc *GladiaClient) error {
		gc.timeout = t
		return nil
	}
}

func NewGladiaClient(opts ...Setter) (*GladiaClient, error) {
	gc := &GladiaClient{
		log:         logger.DefaultLogger(),
		timeout:     0,
		baseUrlPath: "https://api.gladia.io",
	}

	for _, set := range opts {
		if e := set(gc); e != nil {
			return nil, e
		}
	}

	gc.client = &http.Client{Timeout: gc.timeout}
	return gc, nil
}

func (gc *GladiaClient) AudioUploadFromFile(file *os.File) (*upload.UploadResponce, error) {
	path := "/v2/upload"

	body, multipartHeader, err := networking.MakeMultipartBody(file, "audio", file.Name())
	if err != nil {
		return nil, err
	}

	gc.log.Debug("Try do request:", gc.baseUrlPath+path)

	req, err := http.NewRequest("POST", gc.baseUrlPath+path, body)
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

	networking.PrintJson(responseBody)

	err = networking.HttpErrorParse(resp, 200)
	if err != nil {
		return nil, err
	}
	return &responseBody, err
}

func (gc *GladiaClient) InitTranscription(body *prerecorderv2.PreRecorderBody) (*prerecorderv2.PreRecorderInitResponse, error) {
	path := "/v2/pre-recorded"
	method := "POST"

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	gc.log.Debug("Try do request:", gc.baseUrlPath+path)

	networking.PrintJson(body)

	req, err := http.NewRequest(method, gc.baseUrlPath+path, bytes.NewBuffer(jsonBody))
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

	err = networking.HttpErrorParse(resp, 201)
	if err != nil {
		return nil, err
	}

	var responseBody prerecorderv2.PreRecorderInitResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	networking.PrintJson(responseBody)

	return &responseBody, err
}

func (gc *GladiaClient) TranscriptionResult(jobId string, timeInterval time.Duration) (*prerecorderv2.PreRecorderResultResponse, error) {
	path := fmt.Sprintf("/v2/pre-recorded/%s", jobId)
	method := "GET"

	gc.log.Debug("Try do request:", gc.baseUrlPath+path)

	req, err := http.NewRequest(method, gc.baseUrlPath+path, nil)
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

	err = networking.HttpErrorParse(resp, 200)
	if err != nil {
		return nil, err
	}

	var responseBody prerecorderv2.PreRecorderResultResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	networking.PrintJson(responseBody)

	return &responseBody, err
}

func (gc *GladiaClient) AwaitTranscriptionResult(ctx context.Context, jobId string, timeInterval time.Duration) (*prerecorderv2.PreRecorderResultResponse, error) {

	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout waiting result")
		case <-ticker.C:
			resp, err := gc.TranscriptionResult(jobId, timeInterval)
			if err != nil {
				return nil, err
			} else if resp.Status == "error" {
				return nil, fmt.Errorf("error: %s", resp.ErrorCode)
			} else if resp.Status == "done" {
				gc.result = resp.Result
				return resp, nil
			}
		}
	}
}
