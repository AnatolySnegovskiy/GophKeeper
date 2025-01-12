package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestFiles(login string) (cleanup func(), err error) {
	dir := "./.ssh/" + login
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	publicKeyPath := dir + "/public_key.pem"
	privateKeyPath := dir + "/private_key.pem"

	if err := os.WriteFile(publicKeyPath, []byte("public_key_content"), 0600); err != nil {
		return nil, err
	}
	if err := os.WriteFile(privateKeyPath, []byte("private_key_content"), 0600); err != nil {
		return nil, err
	}

	cleanup = func() {
		os.RemoveAll("./.ssh/" + login)
	}

	return cleanup, nil
}

func setupUnreadableFile(filePath string) (cleanup func(), err error) {
	if err := os.WriteFile(filePath, []byte("content"), 0000); err != nil {
		return nil, err
	}

	cleanup = func() {
		os.Remove(filePath)
	}

	return cleanup, nil
}

func TestGetPublicKey(t *testing.T) {
	login := "testuser"
	cleanup, err := setupTestFiles(login)
	if err != nil {
		t.Fatalf("setupTestFiles failed: %v", err)
	}
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		publicKey, err := GetPublicKey(login)
		assert.NoError(t, err)
		assert.Equal(t, "public_key_content", publicKey)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := GetPublicKey("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open private key file")
	})

	t.Run("read error", func(t *testing.T) {
		unreadableFilePath := "./.ssh/" + login + "/unreadable_public_key.pem"
		cleanupUnreadable, err := setupUnreadableFile(unreadableFilePath)
		if err != nil {
			t.Fatalf("setupUnreadableFile failed: %v", err)
		}
		defer cleanupUnreadable()

		_, err = GetPublicKey(login + "/unreadable_public_key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open private key")
	})
}

func TestGetPrivateKey(t *testing.T) {
	login := "testuser"
	cleanup, err := setupTestFiles(login)
	if err != nil {
		t.Fatalf("setupTestFiles failed: %v", err)
	}
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		privateKey, err := GetPrivateKey(login)
		assert.NoError(t, err)
		assert.Equal(t, "private_key_content", privateKey)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := GetPrivateKey("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open private key file")
	})

	t.Run("read error", func(t *testing.T) {
		unreadableFilePath := "./.ssh/" + login + "/unreadable_private_key.pem"
		cleanupUnreadable, err := setupUnreadableFile(unreadableFilePath)
		if err != nil {
			t.Fatalf("setupUnreadableFile failed: %v", err)
		}
		defer cleanupUnreadable()

		_, err = GetPrivateKey(login + "/unreadable_private_key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open private key")
	})
}
