package cli

import (
	"fmt"
	"godotenvpgp/internal/core"
	"os"
	"strings"
)

// Make it test-mockable
var (
	bail         = os.Exit
	findEnvFiles = core.FindEnvFiles
	osStat       = os.Stat
)

func Cli() {
	if len(os.Args) == 1 {
		showHelp()
		return
	}
	core.LoadUnencrypted()

	command := strings.ToLower(os.Args[1])

	var err error

	switch command {
	case "encrypt", "decrypt":
		err = encryptDecrypt(command)
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage: godotenvpgp show <file>")
			bail(0)
			return
		}
		fileName := os.Args[2]
		err = listVars(fileName)
	default:
		showHelp()
	}

	if err != nil {
		fmt.Println(err)
		bail(1)
		return
	}
}

func showHelp() {
	fmt.Print(`
Tool to interact with PGP-encrypted .env files

Usage: godotenvpgp <command>

  encrypt       Encrypt all unencrypted .env files (.env*.unencrypted) to .env.*.encrypted files
  decrypt       Decrypt all encrypted .env files (.env*.encrypted) to .env.*.unencrypted files
  show <file>	Show the content of an encrypted .env file

`)
	bail(0)
}

func encryptDecrypt(command string) error {
	var envFiles []string
	var fn func(string) error

	if command == "encrypt" {
		envFiles = findEnvFiles("unencrypted")
		fn = saveEncryptedFile
	} else if command == "decrypt" {
		envFiles = findEnvFiles("encrypted")
		fn = saveDecryptedFile
	}

	if len(envFiles) == 0 {
		return fmt.Errorf("no files found")
	}

	for _, fileName := range envFiles {
		if err := fn(fileName); err != nil {
			fmt.Printf("> %s failed: %s\n", fileName, err)
		} else {
			fmt.Printf("> %sed: %s\n", command, fileName)
		}
	}
	return nil
}

func listVars(fileName string) error {
	if !strings.HasSuffix(fileName, ".encrypted") {
		return fmt.Errorf("file %s is not encrypted", fileName)
	}

	if _, err := osStat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file %s not found", fileName)
	}

	fileBytes, err := decryptFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot list vars: %s", err)
	}

	fmt.Println(string(fileBytes))

	return nil
}
