package entities

import (
	"encoding/json"
	"os"
)

type FileEntity interface {
	ToFile() (*os.File, error)
	GetName() string
	FromFile(file *os.File) error
}

type BaseFile struct{}

func (b *BaseFile) toFile(data FileEntity) (*os.File, error) {
	fileJson, err := json.Marshal(data)
	if err != nil {
		return &os.File{}, err
	}

	tmpFile, err := os.CreateTemp("", data.GetName())
	if err != nil {
		return &os.File{}, err
	}

	if _, err := tmpFile.Write(fileJson); err != nil {
		return &os.File{}, err
	}
	if err := tmpFile.Close(); err != nil {
		return &os.File{}, err
	}

	return tmpFile, nil
}
