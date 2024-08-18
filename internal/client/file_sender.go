package client

import (
	"encoding/json"
	"fmt"
	"goph_keeper/internal/services/entities"
	"goph_keeper/internal/services/file_helper"
	v1 "goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"os"
)

func (c *GrpcClient) sendFile(stream v1.GophKeeperV1Service_StorePrivateDataClient, filePath string, progressChan chan<- int) error {
	file, fileMetadata, err := c.preparedFile(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	metadataJson, err := json.Marshal(fileMetadata)
	if err != nil {
		return err
	}

	bufSender := make([]byte, 1024*1024)
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
		err = stream.Send(&v1.StorePrivateDataRequest{
			Data:     chunk,
			DataType: v1.DataType_DATA_TYPE_BINARY,
			Metadata: string(metadataJson),
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

func (c *GrpcClient) preparedFile(filePath string) (*os.File, entities.FileMetadata, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, entities.FileMetadata{}, err
	}

	fileMetadata, err := file_helper.GetFileMetadata(file)
	if err != nil {
		return nil, entities.FileMetadata{}, err
	}

	if !fileMetadata.IsCompressed {
		return file, fileMetadata, nil
	}

	file, err = file_helper.CompressGZIP(file)
	if err != nil {
		return nil, entities.FileMetadata{}, err
	}

	fileMetadata.CompressionType = "gzip"

	return file, fileMetadata, err
}
