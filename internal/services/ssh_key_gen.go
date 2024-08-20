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

func (s *SshKeyGen) Generate() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	privateKeyFile, err := os.Create("./.ssh/private_key.pem")
	if err != nil {
		panic(err)
	}

	defer privateKeyFile.Close()
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes}
	publicKeyFile, err := os.Create("./.ssh/public_key.pem")
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		panic(err)
	}

	if _, err := publicKeyFile.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	content, _ := io.ReadAll(publicKeyFile)
	return string(content), nil
}

func (s *SshKeyGen) EncryptMessage(message string, publicKeyData string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyData))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPubKey, _ := publicKey.(*rsa.PublicKey)
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(message))
	encodedMessage := base64.StdEncoding.EncodeToString(encryptedMessage)
	return encodedMessage, nil
}

func (s *SshKeyGen) DecryptionFunction(data string, privateKeyData string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyData))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedData)
	if err != nil {
		return "", err
	}

	return string(decryptedMessage), nil
}
