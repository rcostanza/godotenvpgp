package core

import (
	"os"
	"strings"
)

const DEFAULT_ENVIRONMENT = "DEFAULT"

func ParseEnv(fileContents string) map[string]map[string]string {
	// Values are separated by specific environments. By default variables belong to "all",
	// unless declared under an [environment] section in the env file.
	varEnv := DEFAULT_ENVIRONMENT
	envVars := make(map[string]map[string]string)
	lines := strings.Split(fileContents, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		} else if strings.Contains(line, "[") && strings.Contains(line, "]") {
			varEnv = strings.Trim(line, "[]")
		} else if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "" || value == "" {
				continue
			}
			if _, ok := envVars[varEnv]; !ok {
				envVars[varEnv] = make(map[string]string)
			}
			envVars[varEnv][key] = value
		}
	}
	return envVars
}

// Try to get an environment variable that defines the environment
func GetCurrentEnvironment() string {
	for _, name := range []string{"ENVIRONMENT", "environment", "ENV", "env"} {
		environment := os.Getenv(name)
		if environment != "" {
			return environment
		}
	}
	return DEFAULT_ENVIRONMENT
}

func SetEnv(envVars map[string]map[string]string) {
	currentEnv := GetCurrentEnvironment()
	for env, vars := range envVars {
		if env == currentEnv || env == DEFAULT_ENVIRONMENT {
			for key, value := range vars {
				os.Setenv(key, value)
			}
		}
	}
}
