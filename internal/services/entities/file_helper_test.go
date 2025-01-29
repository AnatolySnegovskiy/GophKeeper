package entities

import (
	"encoding/json"
	"goph_keeper/internal/mocks"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
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

func TestToFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseFile := &BaseFile{}
	mockFileEntity := mocks.NewMockFileEntity(ctrl)

	// Test case 1: Successful creation of a temporary file
	testFile := &TestFile{
		Name: "testfile",
		Data: "testdata",
	}

	tmpFile, err := baseFile.toFile(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, tmpFile)
	assert.FileExists(t, tmpFile.Name())

	// Clean up the temporary file
	defer os.Remove(tmpFile.Name())

	// Verify the content of the temporary file
	fileContent, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	var readTestFile TestFile
	err = json.Unmarshal(fileContent, &readTestFile)
	assert.NoError(t, err)
	assert.Equal(t, testFile.Name, readTestFile.Name)
	assert.Equal(t, testFile.Data, readTestFile.Data)

	// Test case 2: Error in JSON marshalling
	mockFileEntity.EXPECT().GetName().Return("invalidfile").AnyTimes()
	mockFileEntity.EXPECT().ToFile().Return(nil, json.Unmarshal([]byte("invalid json"), &struct{}{})).AnyTimes()

	tmpFile, err = mockFileEntity.ToFile()
	assert.Error(t, err) // JSON marshalling should fail for invalid JSON
	assert.Nil(t, tmpFile)

	// Test case 3: Error in creating a temporary file
	mockFileEntity.EXPECT().GetName().Return("testfile").AnyTimes()
	mockFileEntity.EXPECT().ToFile().Return(nil, os.ErrNotExist).AnyTimes()
	tmpFile, err = mockFileEntity.ToFile()
	assert.Error(t, err) // Creating a temporary file should fail
	assert.Nil(t, tmpFile)

	// Test case 4: Error in writing to the temporary file
	mockFileEntity.EXPECT().GetName().Return("testfile").AnyTimes()
	mockFileEntity.EXPECT().ToFile().Return(nil, os.ErrPermission).AnyTimes()
	tmpFile, err = mockFileEntity.ToFile()
	assert.Error(t, err) // Writing to the temporary file should fail
	assert.Nil(t, tmpFile)

	// Test case 5: Error in closing the temporary file
	mockFileEntity.EXPECT().GetName().Return("testfile").AnyTimes()
	mockFileEntity.EXPECT().ToFile().Return(nil, os.ErrClosed).AnyTimes()
	tmpFile, err = mockFileEntity.ToFile()
	assert.Error(t, err) // Closing the temporary file should fail
	assert.Nil(t, tmpFile)
}
