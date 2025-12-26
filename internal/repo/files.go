package repo

import (
	"encoding/json"
	"errors"
	"os"

	"go-gladia.io-client/pkg/logger"
)

type FilesRepo struct {
	log logger.ILogger
}

func NewFilesRepo(l logger.ILogger) (*FilesRepo, error) {
	if l == nil {
		return nil, errors.New("failed init files repo: logger is empty")
	}

	r := &FilesRepo{log: l}
	return r, nil
}

func (r *FilesRepo) Open(filePath string) (error, *os.File, func() error) {
	file, err := os.Open(filePath)
	if err != nil {
		return err, nil, nil
	}
	return err, file, file.Close
}

func (r *FilesRepo) Write(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

func (r *FilesRepo) JsonDump(filePath string, data any) error {

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	r.Write(filePath, jsonBytes)

	return nil
}

func (r *FilesRepo) TextDump(filePath string, data []byte) error {
	return r.Write(filePath, data)
}
