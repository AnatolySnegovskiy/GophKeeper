package client

import (
	"fmt"
	"goph_keeper/internal/services/entities"
	"goph_keeper/internal/services/file_helper"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"os"
)

func (c *GrpcClient) sendFile(stream v1.GophKeeperV1Service_UploadFileClient, fileMetadata *entities.FileMetadata, filePath string, progressChan chan<- int) error {
	file, err := c.preparedFile(fileMetadata, filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	bufSender := make([]byte, c.sizeChunk)
	resendCounter := 0
	var sentSize int64 = 0
	lastProgress := 0

	_, err = file.Seek(0, io.SeekStart)
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
		err = stream.Send(&v1.UploadFileRequest{
			Data: chunk,
		})

		if err != nil && resendCounter > 3 {
			return err
		} else if err != nil {
			resendCounter += 1
			continue
		}

		resendCounter = 0
		sentSize += int64(num)
		progress := int(float64(sentSize) / float64(fileMetadata.FileSize) * 100)

		if progress > lastProgress {
			lastProgress = progress
			progressChan <- progress
		}
	}

	progressChan <- 100

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
