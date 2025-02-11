package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
)

type SshKeyGen struct {
}

func NewSshKeyGen() *SshKeyGen {
	return &SshKeyGen{}
}

func (s *SshKeyGen) Generate(login string) (string, error) {
	dir := "./.ssh/" + login
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", err
	}

	privateKeyFile, err := os.Create(dir + "/private_key.pem")
	if err != nil {
		return "", err
	}

	defer privateKeyFile.Close()
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return "", err
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return "", err
	}

	publicKeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}
	publicKeyFile, err := os.Create(dir + "/public_key.pem")
	if err != nil {
		return "", err
	}

	defer publicKeyFile.Close()
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return "", err
	}

	if _, err := publicKeyFile.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	content, _ := io.ReadAll(publicKeyFile)
	return string(content), nil
}

func (s *SshKeyGen) EncryptMessage(message []byte, publicKeyData string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKeyData))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, _ := publicKey.(*rsa.PublicKey)
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, message)
	if err != nil {
		return nil, err
	}

	encodedMessage := base64.StdEncoding.EncodeToString(encryptedMessage)
	return []byte(encodedMessage), nil
}

func (s *SshKeyGen) DecryptionFunction(data []byte, privateKeyData string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyData))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedData)
	if err != nil {
		return nil, err
	}

	return decryptedMessage, nil
}
