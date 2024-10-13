package client

import (
	"fmt"
	"io"
	"os"
)

func GetPublicKey(login string) (string, error) {
	filePath := "./.ssh/" + login + "/public_key.pem"
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open private key file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}

	return string(content), nil
}

func GetPrivateKey(login string) (string, error) {
	filePath := "./.ssh/" + login + "/private_key.pem"
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open private key file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %w", err)
	}

	return string(content), nil
}
