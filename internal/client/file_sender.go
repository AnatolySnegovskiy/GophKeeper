package client

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"goph_keeper/internal/services/entities"
	"goph_keeper/internal/services/file_helper"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"os"
)

func (c *GrpcClient) sendFile(stream v1.GophKeeperV1Service_UploadFileClient, fileMetadata *entities.FileMetadata, filePath string, progressChan chan<- int) error {
	file, err := c.preparedFile(fileMetadata, filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	bufSender := make([]byte, c.sizeChunk)
	resendCounter := 0
	var sentSize int64 = 0
	lastProgress := 0

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// Генерация ключа и nonce для ChaCha20
	key := []byte("example key 1234example key 1234") // 32 байта
	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return err
	}

	// Отправка nonce перед началом передачи данных
	err = stream.Send(&v1.UploadFileRequest{
		Data: nonce,
	})
	if err != nil {
		return err
	}

	for {
		num, err := file.Read(bufSender)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := bufSender[:num]
		encryptedChunk := make([]byte, num)
		cipher.XORKeyStream(encryptedChunk, chunk)

		err = stream.Send(&v1.UploadFileRequest{
			Data: encryptedChunk,
		})

		if err != nil && resendCounter > 3 {
			return err
		} else if err != nil {
			resendCounter += 1
			continue
		}

		resendCounter = 0
		sentSize += int64(num)
		if progressChan != nil {
			progress := int(float64(sentSize) / float64(fileMetadata.FileSize) * 100)
			if progress > lastProgress {
				lastProgress = progress
				progressChan <- progress
			}
		}
	}

	if progressChan != nil {
		progressChan <- 100
	}

	c.logger.Info(fmt.Sprintf("File size: %d bytes; sent: %d bytes", fileMetadata.FileSize, sentSize))

	return nil
}

func (c *GrpcClient) preparedFile(metadata *entities.FileMetadata, filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	if !metadata.IsCompressed {
		return file, nil
	}

	file, err = file_helper.CompressGZIP(file)
	if err != nil {
		return nil, err
	}

	metadata.CompressionType = "gzip"

	return file, err
}
