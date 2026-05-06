// Package encryptor provides AES-GCM encryption and decryption for
// environment variable values, enabling secure storage of sensitive .env files.
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("encryptor: invalid ciphertext")

// deriveKey produces a 32-byte AES-256 key from the given passphrase.
func deriveKey(passphrase string) []byte {
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

// Encrypt encrypts each value in vars using AES-256-GCM with the provided
// passphrase. Returns a new map with base64-encoded ciphertext values.
func Encrypt(vars map[string]string, passphrase string) (map[string]string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
		sealed := gcm.Seal(nonce, nonce, []byte(v), nil)
		out[k] = base64.StdEncoding.EncodeToString(sealed)
	}
	return out, nil
}

// Decrypt decrypts each value in vars using AES-256-GCM with the provided
// passphrase. Returns a new map with plaintext values.
func Decrypt(vars map[string]string, passphrase string) (map[string]string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(vars))
	for k, v := range vars {
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil || len(data) < gcm.NonceSize() {
			return nil, ErrInvalidCiphertext
		}
		nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
		plain, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, ErrInvalidCiphertext
		}
		out[k] = string(plain)
	}
	return out, nil
}
