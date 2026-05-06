package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/envlayer/internal/encryptor"
	"github.com/your-org/envlayer/internal/loader"
)

// EncryptOptions configures the encrypt/decrypt CLI commands.
type EncryptOptions struct {
	InputFile  string
	OutputFile string
	Passphrase string
	Decrypt    bool
}

// RunEncrypt loads a .env file, encrypts (or decrypts) its values and writes
// the result as a JSON object to OutputFile (or stdout when empty).
func RunEncrypt(opts EncryptOptions) error {
	if opts.Passphrase == "" {
		return fmt.Errorf("encrypt: passphrase must not be empty")
	}

	vars, err := loader.LoadFile(opts.InputFile)
	if err != nil {
		return fmt.Errorf("encrypt: load %q: %w", opts.InputFile, err)
	}

	var result map[string]string
	if opts.Decrypt {
		result, err = encryptor.Decrypt(vars, opts.Passphrase)
		if err != nil {
			return fmt.Errorf("encrypt: decrypt: %w", err)
		}
	} else {
		result, err = encryptor.Encrypt(vars, opts.Passphrase)
		if err != nil {
			return fmt.Errorf("encrypt: encrypt: %w", err)
		}
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("encrypt: marshal: %w", err)
	}

	if opts.OutputFile == "" {
		fmt.Println(string(data))
		return nil
	}

	if err := os.WriteFile(opts.OutputFile, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("encrypt: write %q: %w", opts.OutputFile, err)
	}
	return nil
}
