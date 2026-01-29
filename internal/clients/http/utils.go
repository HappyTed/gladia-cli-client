package http_client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func makeMultipartBody(file *os.File, key string, value string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile(key, value)
	if err != nil {
		return nil, "", err
	}

	io.Copy(part, file)

	return body, writer.FormDataContentType(), nil
}

func httpErrorParse(resp *http.Response, expectedStatusCode int) error {

	if resp.StatusCode != expectedStatusCode {
		return errors.New("unexpected status code")
	}

	return nil
}
