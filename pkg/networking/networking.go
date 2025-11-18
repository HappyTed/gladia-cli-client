package networking

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func MakeMultipartBody(file *os.File, key string, value string) (*bytes.Buffer, string, error) {
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

func HttpErrorParse(resp *http.Response, expectedStatusCode int) error {

	fmt.Println("Request Status:", resp.Status)
	fmt.Println("Status code:", resp.StatusCode)
	if resp.StatusCode != expectedStatusCode {
		return errors.New("unexpected status code")
	}

	return nil
}

func JsonDumpS(data any) (string, error) {
	fmt.Println("Try print JSON...")
	s, err := json.MarshalIndent(data, "", "  ")
	return string(s), err
}

func PrintJson(data any) {
	fmt.Println("Try print JSON...")
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("JSON conversion error:", err)
		return
	}
	fmt.Println(string(b))
}
