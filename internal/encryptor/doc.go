// Package encryptor provides symmetric encryption and decryption of
// environment variable maps using AES-256-GCM.
//
// Keys are derived from a user-supplied passphrase via SHA-256 so no
// separate key-management infrastructure is required for local use.
//
// Typical usage:
//
//	// Encrypt before persisting sensitive vars
//	sealed, err := encryptor.Encrypt(vars, passphrase)
//
//	// Decrypt after loading from disk
//	plain, err := encryptor.Decrypt(sealed, passphrase)
//
// Each call to Encrypt generates a fresh random nonce per value, so
// identical plaintexts produce different ciphertexts across calls.
package encryptor
