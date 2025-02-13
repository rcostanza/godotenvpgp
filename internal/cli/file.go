package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/rcostanza/godotenvpgp/internal/core"
)

// Make it test-mockable
var (
	writeFile   = os.WriteFile
	encryptFile = core.EncryptFile
	decryptFile = core.DecryptFile
)

func saveEncryptedFile(fileName string) error {
	fileBytes, err := encryptFile(fileName)
	if err != nil {
		return fmt.Errorf("[godotenvpgp][Error] Cannot encrypt env file: %s", err)
	}

	targetFile := strings.Replace(fileName, ".unencrypted", ".encrypted", 1)
	err = writeFile(targetFile, []byte(fileBytes), 0644)
	if err != nil {
		return fmt.Errorf("[godotenvpgp][Error] Cannot write encrypted env file: %s", err)
	}
	return nil
}

func saveDecryptedFile(fileName string) error {
	fileContents, err := decryptFile(fileName)
	if err != nil {
		return fmt.Errorf("[godotenvpgp][Error] Cannot decrypt env file: %s", err)
	}
	targetFile := strings.Replace(fileName, ".encrypted", ".unencrypted", 1)
	err = writeFile(targetFile, []byte(fileContents), 0644)
	if err != nil {
		return fmt.Errorf("[godotenvpgp][Error] Cannot write decrypted env file: %s", err)
	}
	return nil
}
