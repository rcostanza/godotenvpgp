package core

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func LoadUnencrypted() (err error) {
	bFileContents, err := readFile(".env")
	if err != nil {
		return
	}
	fileContents := string(bFileContents)
	SetEnv(ParseEnv(fileContents))
	return
}

func LoadEncrypted(fileName string) (err error) {
	fileBytes, err := DecryptFile(fileName)
	if err != nil {
		return
	}
	fileContents := string(fileBytes)
	SetEnv(ParseEnv(fileContents))
	return
}

// Return a list of encrypted env files found in the root directory,
// and their estimated specific environment if applicable
func FindEnvFiles(fileType string) []string {
	var envFiles []string

	files, err := readDir(".")
	if err != nil {
		fmt.Println(fmt.Errorf("[godotenvpgp][Error] Cannot read root directory: %w", err))
		return nil
	}

	re := regexp.MustCompile(envFilePattern)

	for _, file := range files {
		fileName := file.Name()
		if !re.MatchString(fileName) || !strings.HasSuffix(fileName, "."+fileType) {
			continue
		}
		envFiles = append(envFiles, file.Name())
	}

	// Reverse sort so that the regular .env file is parsed first
	sort.Sort(sort.Reverse(sort.StringSlice(envFiles)))

	return envFiles
}
