package entities_test

import (
	"encoding/json"
	"goph_keeper/internal/services/entities"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileCard_ToFile(t *testing.T) {
	fileCard := &entities.FileCard{
		Uuid:        "1234",
		CardName:    "Test Card",
		Description: "This is a test card",
		CardNumber:  "1111-2222-3333-4444",
		ExpiryDate:  "12/23",
		CVV:         "123",
		CardHolder:  "John Doe",
	}

	file, err := fileCard.ToFile()
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // Ensure the temporary file is removed
	defer file.Close()
	file, _ = os.Open(file.Name())
	// Read the content and unmarshal it to a new FileCard
	content, err := io.ReadAll(file)
	assert.NoError(t, err)

	var newFileCard entities.FileCard
	err = json.Unmarshal(content, &newFileCard)
	assert.NoError(t, err)
	assert.Equal(t, fileCard, &newFileCard)
}

func TestFileCard_GetName(t *testing.T) {
	fileCard := &entities.FileCard{
		CardName: "Test Card",
	}
	assert.Equal(t, "Test Card", fileCard.GetName())
}

func TestFileCard_FromFile(t *testing.T) {
	originalFileCard := &entities.FileCard{
		Uuid:        "1234",
		CardName:    "Test Card",
		Description: "This is a test card",
		CardNumber:  "1111-2222-3333-4444",
		ExpiryDate:  "12/23",
		CVV:         "123",
		CardHolder:  "John Doe",
	}

	// Create a temporary file and write the originalFileCard data to it
	tmpFile, err := os.CreateTemp("", "testfile_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = json.NewEncoder(tmpFile).Encode(originalFileCard)
	assert.NoError(t, err)
	tmpFile.Close()

	// Reopen the file for reading
	tmpFile, err = os.Open(tmpFile.Name())
	assert.NoError(t, err)
	defer tmpFile.Close()

	// Test FromFile method
	newFileCard := &entities.FileCard{}
	err = newFileCard.FromFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, originalFileCard, newFileCard)
}
