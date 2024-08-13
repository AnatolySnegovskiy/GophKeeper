package file_helper

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressGZIP(file *os.File) (*os.File, error) {
	gzipFile, err := os.CreateTemp("", "compressed_*.gz")
	defer gzipFile.Close()
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(gzipFile)
	gzipWriter.Close()
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return nil, err
	}

	return gzipFile, nil
}
