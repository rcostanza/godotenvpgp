package core

import (
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// Make it test-mockable
var (
	readFile = os.ReadFile
	readDir  = os.ReadDir
)

func EncryptFile(fileName string) (fileBytes []byte, err error) {
	password, err := getFilePassword(fileName)
	if err != nil {
		err = fmt.Errorf("[godotenvpgp][Error] Cannot get password: %w", err)
		return
	}

	bFileContents, err := readFile(fileName)
	if err != nil {
		err = fmt.Errorf("[godotenvpgp][Error] Cannot load unencrypted env file: %w", err)
		return
	}

	pgp := crypto.PGP()
	encHandle, _ := pgp.Encryption().Password([]byte(password)).New()
	pgpMessage, _ := encHandle.Encrypt([]byte(bFileContents))
	fileBytes, _ = pgpMessage.ArmorBytes()

	return
}

func DecryptFile(fileName string) (fileBytes []byte, err error) {
	password, err := getFilePassword(fileName)
	if err != nil {
		err = fmt.Errorf("[godotenvpgp][Error] Cannot get password: %w", err)
		return
	}

	bFileContents, err := readFile(fileName)
	if err != nil {
		err = fmt.Errorf("[godotenvpgp][Error] Cannot load encrypted env file: %w", err)
		return
	}

	pgp := crypto.PGP()
	decHandle, _ := pgp.Decryption().Password([]byte(password)).New()
	decrypted, err := decHandle.Decrypt(bFileContents, crypto.Armor)
	if err != nil {
		err = fmt.Errorf("[godotenvpgp][Error] Cannot decrypt env file: %w", err)
		return
	}

	fileBytes = decrypted.Bytes()
	return
}
