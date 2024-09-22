package entities

import (
	"encoding/json"
	"os"
)

type FilePassword struct {
	BaseFile
	Uuid        string
	Title       string `json:"title"`
	Description string `json:"description"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

func (f *FilePassword) ToFile() (*os.File, error) {
	return f.toFile(f)
}

func (f *FilePassword) GetName() string {
	return f.Title
}
func (f *FilePassword) FromFile(file *os.File) error {
	decoder := json.NewDecoder(file)
	return decoder.Decode(f)
}
