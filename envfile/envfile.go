package envfile

import (
	"fmt"
	"log"

	"github.com/rcostanza/godotenvpgp/internal/core"
)

var (
	loadUnencrypted = core.LoadUnencrypted
	findEnvFiles    = core.FindEnvFiles
	loadEncrypted   = core.LoadEncrypted
)

// Loads env files, encrypted and unencrypted.
func Load() error {
	// Tries to load a regular, unencrypted env file, that could contain the PGP_PASSWORD
	if err := loadUnencrypted(); err != nil {
		log.Println("[godotenvpgp][Error] Cannot parse unencrypted env file: ", err)
	}

	// Look for encrypted env files
	encryptedFiles := findEnvFiles("encrypted")

	if len(encryptedFiles) == 0 {
		log.Println("[godotenvpgp][Warn] No encrypted env files found")
	} else {
		for _, fileName := range encryptedFiles {
			// Skip file if it's environment-specific and doesn't match the current environment
			if err := loadEncrypted(fileName); err != nil {
				return fmt.Errorf("[godotenvpgp][Error] Cannot load encrypted env file %s: %s", fileName, err)
			}
		}
	}

	return nil
}
