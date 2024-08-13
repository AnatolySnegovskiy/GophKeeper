package file_helper

import (
	"goph_keeper/internal/services/entities"
	"net/http"
	"os"
	"path/filepath"
)

func GetFileMetadata(file *os.File) (entities.FileMetadata, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return entities.FileMetadata{}, err
	}

	fileInfo, err := file.Stat()

	return entities.FileMetadata{
		FileName:        fileInfo.Name(),
		FileExtension:   filepath.Ext(fileInfo.Name()),
		MemType:         http.DetectContentType(buf),
		IsCompressed:    canCompressBytes(buf),
		CompressionType: "",
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
