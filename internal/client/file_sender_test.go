package client

import (
	"errors"
	"google.golang.org/grpc"
	"goph_keeper/internal/services/entities"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockUploadFileClient struct {
	grpc.ClientStream
	sendFunc  func(*v1.UploadFileRequest) error
	recvFunc  func() (*v1.UploadFileResponse, error)
	closeFunc func() error
}

func (m *MockUploadFileClient) Send(req *v1.UploadFileRequest) error {
	return m.sendFunc(req)
}

func (m *MockUploadFileClient) CloseAndRecv() (*v1.UploadFileResponse, error) {
	return m.recvFunc()
}

func (m *MockUploadFileClient) CloseSend() error {
	return m.closeFunc()
}

func TestSendFile(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Записываем данные в временный файл
	data := []byte("test data")
	_, err = tempFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Создаем метаданные файла
	metadata := &entities.FileMetadata{
		FileSize:     int64(len(data)),
		IsCompressed: false,
		FileName:     "testfile",
		MemType:      "file",
	}

	// Создаем клиент
	client := &GrpcClient{
		sizeChunk: 1024,
		logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	// Создаем мок для stream
	mockStream := &MockUploadFileClient{
		sendFunc: func(req *v1.UploadFileRequest) error {
			return nil
		},
		recvFunc: func() (*v1.UploadFileResponse, error) {
			return &v1.UploadFileResponse{}, nil
		},
		closeFunc: func() error {
			return nil
		},
	} // Создаем канал для прогресса
	progressChan := make(chan int, 1)

	// Тест: успешная отправка файла
	t.Run("Successful file send", func(t *testing.T) {
		done := make(chan struct{})
		progressChan := make(chan int, 1)

		go func() {
			defer close(done)
			err := client.sendFile(mockStream, metadata, tempFile.Name(), progressChan)
			close(progressChan)
			assert.NoError(t, err)
		}()

		for progress := range progressChan {
			if progress < 100 {
				continue
			}
			assert.Equal(t, 100, progress)
		}
		<-done
	})

	// Тест: ошибка при открытии файла
	t.Run("Error opening file", func(t *testing.T) {
		err := client.sendFile(mockStream, metadata, "nonexistentfile", progressChan)
		assert.Error(t, err)
	})

	// Тест: ошибка при сжатии файла
	t.Run("Error compressing file", func(t *testing.T) {
		metadata.IsCompressed = true
		err := client.sendFile(mockStream, metadata, tempFile.Name(), progressChan)
		assert.Error(t, err)
	})

	// Тест: ошибка при отправке данных
	t.Run("Error sending data", func(t *testing.T) {
		progressChan := make(chan int, 1)
		mockStream := &MockUploadFileClient{
			sendFunc: func(req *v1.UploadFileRequest) error {
				return errors.New("send error")
			},
			recvFunc: func() (*v1.UploadFileResponse, error) {
				return nil, errors.New("recv error")
			},
			closeFunc: func() error {
				return errors.New("close error")
			},
		} // Создаем канал для прогресса

		done := make(chan struct{})
		go func() {
			defer close(done)
			err := client.sendFile(mockStream, metadata, tempFile.Name(), progressChan)
			assert.Error(t, err)
		}()
		close(progressChan)
		<-done
	})
}

func TestPreparedFile(t *testing.T) {
	// Создаем временный файл для тестирования
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Закрываем файл, чтобы его можно было открыть в тесте
	tempFile.Close()

	// Создаем метаданные файла
	metadata := &entities.FileMetadata{
		FileName:        "testfile",
		FileExtension:   ".txt",
		MemType:         "text/plain",
		IsCompressed:    false,
		CompressionType: "",
		FileSize:        0,
	}

	// Создаем клиент
	client := &GrpcClient{}

	// Тест: успешное открытие файла без сжатия
	t.Run("Successful open without compression", func(t *testing.T) {
		file, err := client.preparedFile(metadata, tempFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.Equal(t, "", metadata.CompressionType)
		file.Close()
	})

	// Тест: успешное открытие файла с сжатием
	t.Run("Successful open with compression", func(t *testing.T) {
		metadata.IsCompressed = true
		file, err := client.preparedFile(metadata, tempFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.Equal(t, "gzip", metadata.CompressionType)
		file.Close()
	})

	// Тест: ошибка при открытии файла
	t.Run("Error opening file", func(t *testing.T) {
		metadata.IsCompressed = false
		_, err := client.preparedFile(metadata, "nonexistentfile")
		assert.Error(t, err)
	})
}
