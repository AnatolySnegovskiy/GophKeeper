package file_helper

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func CompressGZIP(file *os.File) (*os.File, error) {
	gzipFile, err := os.CreateTemp("", "compressed_*.gz")
	defer func(gzipFile *os.File) {
		err := gzipFile.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(gzipFile)
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(gzipFile)
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return gzipFile, nil
}
