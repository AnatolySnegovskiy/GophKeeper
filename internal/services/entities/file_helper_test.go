package entities

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFile struct {
	BaseFile
	Name string `json:"name"`
	Data string `json:"data"`
}

func (t *TestFile) ToFile() (*os.File, error) {
	return t.toFile(t)
}

func (t *TestFile) GetName() string {
	return t.Name
}

func (t *TestFile) FromFile(file *os.File) error {
	decoder := json.NewDecoder(file)
	return decoder.Decode(t)
}

func TestBaseFile_toFile(t *testing.T) {
	testFile := &TestFile{
		Name: "testfile",
		Data: "This is some test data",
	}

	file, err := testFile.ToFile()
	assert.NoError(t, err)
	defer os.Remove(file.Name()) // Ensure the temporary file is removed
	defer file.Close()
	file, _ = os.Open(file.Name())
	// Read the content and unmarshal it to a new TestFile
	content, err := io.ReadAll(file)
	assert.NoError(t, err)

	var newTestFile TestFile
	err = json.Unmarshal(content, &newTestFile)
	assert.NoError(t, err)
	assert.Equal(t, testFile, &newTestFile)
}

func TestTestFile_GetName(t *testing.T) {
	testFile := &TestFile{
		Name: "testfile",
	}
	assert.Equal(t, "testfile", testFile.GetName())
}

func TestTestFile_FromFile(t *testing.T) {
	originalTestFile := &TestFile{
		Name: "testfile",
		Data: "This is some test data",
	}

	// Create a temporary file and write the originalTestFile data to it
	tmpFile, err := os.CreateTemp("", "testfile_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = json.NewEncoder(tmpFile).Encode(originalTestFile)
	assert.NoError(t, err)
	tmpFile.Close()

	// Reopen the file for reading
	tmpFile, err = os.Open(tmpFile.Name())
	assert.NoError(t, err)
	defer tmpFile.Close()

	// Test FromFile method
	newTestFile := &TestFile{}
	err = newTestFile.FromFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, originalTestFile, newTestFile)
}
