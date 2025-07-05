// encryption/encryption.go
// This package encapsulates all encryption and decryption logic using AES-GCM.
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts data using AES-GCM.
// It takes plaintext data and a key, and returns a hex-encoded ciphertext string.
func Encrypt(data []byte, key []byte) (string, error) {
	// Create a new AES cipher block from the key.
	// The key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM block cipher, which provides authenticated encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce (number used once). Its size is determined by the GCM implementation.
	// A random nonce is used here for security.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data. The nonce is prepended to the ciphertext.
	// This is a common practice as the nonce is required for decryption but doesn't need to be secret.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM.
// It takes a hex-encoded ciphertext string and the key, and returns the original plaintext data.
func Decrypt(encryptedData string, key []byte) ([]byte, error) {
	// Decode the hex-encoded string back into bytes.
	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Create a new AES cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM block cipher.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Get the nonce size.
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	// Extract the nonce from the beginning of the ciphertext.
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data. If the authentication tag (part of the ciphertext)
	// doesn't match, this will return an error, protecting against tampering.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
