package services_test

import (
	"goph_keeper/internal/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	sshKeyGen := services.NewSshKeyGen()
	login := "testuser"

	// Clean up after test
	defer os.RemoveAll("./.ssh/" + login)

	publicKey, err := sshKeyGen.Generate(login)
	assert.NoError(t, err)
	assert.NotEmpty(t, publicKey)

	// Confirm private key file exists
	privateKeyPath := "./.ssh/" + login + "/private_key.pem"
	_, err = os.Stat(privateKeyPath)
	assert.NoError(t, err)

	// Confirm public key file exists
	publicKeyPath := "./.ssh/" + login + "/public_key.pem"
	_, err = os.Stat(publicKeyPath)
	assert.NoError(t, err)

	// Confirm public key content matches the returned public key
	content, err := os.ReadFile(publicKeyPath)
	assert.NoError(t, err)
	assert.Equal(t, publicKey, string(content))
}

func TestEncryptMessage(t *testing.T) {
	sshKeyGen := services.NewSshKeyGen()
	login := "testuser"

	// Generate keys
	publicKey, err := sshKeyGen.Generate(login)
	assert.NoError(t, err)

	message := []byte("Hello, World!")
	encryptedMessage, err := sshKeyGen.EncryptMessage(message, publicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedMessage)
}

func TestDecryptionFunction(t *testing.T) {
	sshKeyGen := services.NewSshKeyGen()
	login := "testuser"

	// Generate keys
	_, err := sshKeyGen.Generate(login)
	assert.NoError(t, err)

	// Read private key
	privateKeyPath := "./.ssh/" + login + "/private_key.pem"
	privateKey, err := os.ReadFile(privateKeyPath)
	assert.NoError(t, err)

	// Read public key
	publicKeyPath := "./.ssh/" + login + "/public_key.pem"
	publicKey, err := os.ReadFile(publicKeyPath)
	assert.NoError(t, err)

	message := []byte("Hello, World!")
	encryptedMessage, err := sshKeyGen.EncryptMessage(message, string(publicKey))
	assert.NoError(t, err)

	decryptedMessage, err := sshKeyGen.DecryptionFunction(encryptedMessage, string(privateKey))
	assert.NoError(t, err)
	assert.Equal(t, message, decryptedMessage)
}
