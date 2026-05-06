package encryptor_test

import (
	"testing"

	"github.com/your-org/envlayer/internal/encryptor"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123xyz",
	}
	passphrase := "my-strong-passphrase"

	encrypted, err := encryptor.Encrypt(vars, passphrase)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	for k, v := range vars {
		if encrypted[k] == v {
			t.Errorf("key %q: expected ciphertext, got plaintext", k)
		}
	}

	decrypted, err := encryptor.Decrypt(encrypted, passphrase)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	for k, want := range vars {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncrypt_UniqueNonces(t *testing.T) {
	vars := map[string]string{"SECRET": "value"}
	passphrase := "passphrase"

	a, _ := encryptor.Encrypt(vars, passphrase)
	b, _ := encryptor.Encrypt(vars, passphrase)

	if a["SECRET"] == b["SECRET"] {
		t.Error("expected different ciphertexts due to random nonces")
	}
}

func TestDecrypt_WrongPassphrase_ReturnsError(t *testing.T) {
	vars := map[string]string{"KEY": "value"}
	encrypted, _ := encryptor.Encrypt(vars, "correct")

	_, err := encryptor.Decrypt(encrypted, "wrong")
	if err == nil {
		t.Error("expected error with wrong passphrase, got nil")
	}
}

func TestDecrypt_InvalidBase64_ReturnsError(t *testing.T) {
	vars := map[string]string{"KEY": "!!!not-base64!!!"}
	_, err := encryptor.Decrypt(vars, "passphrase")
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}
}

func TestEncryptDecrypt_EmptyMap(t *testing.T) {
	vars := map[string]string{}
	enc, err := encryptor.Encrypt(vars, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dec, err := encryptor.Decrypt(enc, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dec) != 0 {
		t.Errorf("expected empty map, got %v", dec)
	}
}
