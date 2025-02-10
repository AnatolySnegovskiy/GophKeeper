package fileHelper

import (
	"goph_keeper/internal/services/entities"
	"net/http"
	"os"
	"path/filepath"
)

func GetFileMetadata(filePath string) (*entities.FileMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &entities.FileMetadata{
		FileName:        fileInfo.Name(),
		FileExtension:   filepath.Ext(fileInfo.Name()),
		MemType:         http.DetectContentType(buf),
		IsCompressed:    canCompressBytes(buf),
		CompressionType: "",
		FileSize:        fileInfo.Size(),
	}, nil
}

func canCompressBytes(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 { // ASCII printable characters
			return false // Бинарный файл
		}
	}

	return true
}
