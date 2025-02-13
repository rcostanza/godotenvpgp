package core

import (
	"os"
	"slices"
	"strings"
)

const DEFAULT_ENVIRONMENT = "__DEFAULT"

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

	// Sort environments so the default one's variables are set first, so they can be overriden
	// by environment-specific values if needed.
	envNames := make([]string, 0, len(envVars))
	for env := range envVars {
		envNames = append(envNames, env)
	}
	slices.Sort(envNames)

	for _, envName := range envNames {
		if envName == currentEnv || envName == DEFAULT_ENVIRONMENT {
			for key, value := range envVars[envName] {
				os.Setenv(key, value)
			}
		}
	}

}
