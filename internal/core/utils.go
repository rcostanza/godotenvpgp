package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const envFilePattern = `(?m)^\.env(\..+)*\.(?:un){0,1}encrypted$`

func getFilePassword(fileName string) (password string, err error) {
	environment, err := getFileEnvironment(fileName)
	if err != nil {
		return
	}

	key := "ENVFILE_PASSWORD"
	if environment != "" && environment != DEFAULT_ENVIRONMENT {
		key += "_" + strings.ToUpper(environment)
	}

	password = os.Getenv(key)
	if password == "" {
		err = fmt.Errorf("password not found; environment variable %s missing", key)
	}
	return
}

func getFileEnvironment(fileName string) (string, error) {
	var re = regexp.MustCompile(envFilePattern)

	if !re.MatchString(fileName) {
		return "", fmt.Errorf("not a valid .env file")
	}

	groups := re.FindStringSubmatch(fileName)
	env, found := strings.CutPrefix(groups[1], ".")
	if !found {
		env = DEFAULT_ENVIRONMENT
	}

	return env, nil
}
