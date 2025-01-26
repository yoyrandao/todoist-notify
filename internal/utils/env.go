package utils

import "os"

func GetenvOrDefault(key, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}

	return value
}
