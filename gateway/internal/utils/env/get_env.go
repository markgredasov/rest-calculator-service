package utils_env

import (
	"fmt"
	"os"
)

func GetEnvironments(envs ...string) (map[string]string, error) {
	result := make(map[string]string, len(envs))
	for _, env := range envs {
		value := os.Getenv(env)
		if value == "" {
			return map[string]string{}, fmt.Errorf("failed to get env: %s", env)
		}

		result[env] = value
	}

	return result, nil
}
