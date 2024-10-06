package file_helper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileMetadata(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write some content to the file
	content := []byte("Hello, World!")
	_, err = tmpFile.Write(content)
	assert.NoError(t, err)

	// Close the file
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Test GetFileMetadata
	fileMetadata, err := GetFileMetadata(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, fileMetadata)
	assert.Equal(t, "application/octet-stream", fileMetadata.MemType)
	assert.Equal(t, int64(len(content)), fileMetadata.FileSize)
	assert.False(t, fileMetadata.IsCompressed)
}

func TestCanCompressBytes(t *testing.T) {
	// Test with ASCII printable characters (should return true)
	printableData := []byte("Hello, World!")
	assert.True(t, canCompressBytes(printableData))

	// Test with non-printable characters (should return false)
	nonPrintableData := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	assert.False(t, canCompressBytes(nonPrintableData))
}
