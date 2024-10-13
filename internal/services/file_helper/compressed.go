package file_helper

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressGZIP(file *os.File) (*os.File, error) {
	gzipFile, err := os.CreateTemp("", "compressed_*.gz")
	if err != nil {
		return nil, err
	}

	defer gzipFile.Close()

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
