package file_helper_test

import (
	"compress/gzip"
	"goph_keeper/internal/services/file_helper"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressGZIP(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write some content to the file
	content := []byte("Hello, World!")
	_, err = tmpFile.Write(content)
	assert.NoError(t, err)

	// Close the file
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Reopen the file for reading
	tmpFile, err = os.Open(tmpFile.Name())
	assert.NoError(t, err)
	defer tmpFile.Close()

	// Test CompressGZIP function
	compressedFile, err := file_helper.CompressGZIP(tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, compressedFile)
	defer compressedFile.Close()
	defer os.Remove(compressedFile.Name())
	compressedFilePath := compressedFile.Name()
	compressedFile, err = os.Open(compressedFilePath) // Reopen the file for reading

	// Verify the compressed file
	compressedFileStat, err := compressedFile.Stat()
	assert.NoError(t, err)
	assert.True(t, compressedFileStat.Size() > 0)

	// Reopen compressed file for reading
	compressedFile, err = os.Open(compressedFile.Name())
	assert.NoError(t, err)

	// Decompress and verify content
	gzipReader, err := gzip.NewReader(compressedFile)
	assert.NoError(t, err)
	defer gzipReader.Close()

	decompressedContent, err := io.ReadAll(gzipReader)
	assert.NoError(t, err)
	assert.Equal(t, content, decompressedContent)
}
