package entities_test

import (
	"encoding/json"
	"goph_keeper/internal/services/entities"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilePassword_ToFile(t *testing.T) {
	filePassword := &entities.FilePassword{
		Uuid:        "1234",
		Title:       "Test Password",
		Description: "This is a test password",
		Login:       "test_login",
		Password:    "test_password",
	}

	file, err := filePassword.ToFile()
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // Ensure the temporary file is removed
	defer file.Close()
	file, _ = os.Open(file.Name())
	// Read the content and unmarshal it to a new FilePassword
	content, err := io.ReadAll(file)
	assert.NoError(t, err)

	var newFilePassword entities.FilePassword
	err = json.Unmarshal(content, &newFilePassword)
	assert.NoError(t, err)
	assert.Equal(t, filePassword, &newFilePassword)
}

func TestFilePassword_GetName(t *testing.T) {
	filePassword := &entities.FilePassword{
		Title: "Test Password",
	}
	assert.Equal(t, "Test Password", filePassword.GetName())
}

func TestFilePassword_FromFile(t *testing.T) {
	originalFilePassword := &entities.FilePassword{
		Uuid:        "1234",
		Title:       "Test Password",
		Description: "This is a test password",
		Login:       "test_login",
		Password:    "test_password",
	}

	// Create a temporary file and write the originalFilePassword data to it
	tmpFile, err := os.CreateTemp("", "testfile_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = json.NewEncoder(tmpFile).Encode(originalFilePassword)
	assert.NoError(t, err)
	tmpFile.Close()

	// Reopen the file for reading
	tmpFile, err = os.Open(tmpFile.Name())
	assert.NoError(t, err)
	defer tmpFile.Close()

	// Test FromFile method
	newFilePassword := &entities.FilePassword{}
	err = newFilePassword.FromFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, originalFilePassword, newFilePassword)
}
