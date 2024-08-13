package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"goph_keeper/internal/services/grpc/goph_keeper/v1"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func handleError(logger *slog.Logger, err error) {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	grpcClient := v1.NewGophKeeperV1ServiceClient(conn)
	token, err := grpcClient.AuthenticateUser(context.Background(), &v1.AuthenticateUserRequest{
		Username: "User1",
		Password: "test",
	})
	if err != nil {
		log.Fatal(err)
	}

	md := metadata.New(map[string]string{"authorization": token.JwtToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := grpcClient.StorePrivateData(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join(cwd, "cmd/client/storage/Kingdom.of.the.Planet.of.the.Apes.2024.D.WEBRip.1O8Op.mkv")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Определяем тип файла
	fileType := http.DetectContentType(buf)
	fmt.Println("File type:", fileType)

	// Определяем возможность сжатия файла
	canCompress := canCompressBytes(buf)
	fmt.Println("Can compress:", canCompress)

	// Создаем структуру для мета данных файла
	fileMetadata := FileMetadata{
		FileName:      filepath.Base(filePath),
		FileExtension: filepath.Ext(filePath),
		IsCompressed:  canCompress,
	}

	// Сжимаем файл, если это возможно
	if canCompress {
		gzipFile, err := os.CreateTemp("", "compressed_*.gz")
		if err != nil {
			log.Fatal(err)
		}
		defer gzipFile.Close()

		gzipWriter := gzip.NewWriter(gzipFile)
		_, err = io.Copy(gzipWriter, file)
		if err != nil {
			log.Fatal(err)
		}
		gzipWriter.Close()
		fileMetadata.CompressionType = "gzip"
		file = gzipFile
	}

	// Разбиваем сжатый файл на сегменты и отправляем в потоке
	batchNumber := 1
	segmentSize := 102400
	bufSender := make([]byte, segmentSize)
	resendCounter := 0

	for {
		num, err := file.Read(bufSender)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		chunk := bufSender[:num]
		metadata, err := json.Marshal(fileMetadata)
		if err != nil {
			log.Fatal(err)
		}
		if err := stream.Send(&v1.StorePrivateDataRequest{
			Data:     chunk,
			DataType: v1.DataType_DATA_TYPE_BINARY,
			Metadata: string(metadata),
		}); err != nil {
			if resendCounter < 3 {
				resendCounter += 1
				continue
			}
			log.Fatal(err)
		}
		log.Printf("Sent - batch #%v - size - %v\n", batchNumber, len(chunk))
		resendCounter = 0
		batchNumber += 1
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent - %v bytes - %s\n", res.Success, res.Message)
}
