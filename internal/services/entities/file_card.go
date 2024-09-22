package entities

import (
	"encoding/json"
	"os"
)

type FileCard struct {
	BaseFile
	Uuid        string
	CardName    string `json:"card_name"`
	Description string `json:"description"`
	CardNumber  string `json:"card_number"`
	ExpiryDate  string `json:"expiry_date"`
	CVV         string `json:"cvv"`
	CardHolder  string `json:"card_holder"`
}

func (f *FileCard) ToFile() (*os.File, error) {
	return f.toFile(f)
}

func (f *FileCard) GetName() string {
	return f.CardName
}
func (f *FileCard) FromFile(file *os.File) error {
	decoder := json.NewDecoder(file)
	return decoder.Decode(f)
}
