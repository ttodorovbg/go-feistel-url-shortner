package cli

import (
	"errors"
	"os"
)

func getKey(flagKey string) (string, error) {

	if flagKey != "" {
		return flagKey, nil
	}

	envKey := os.Getenv("FEISTEL_URL_SHORTENER_KEY")
	if envKey != "" {
		return envKey, nil
	}

	return "", errors.New("missing secret key")
}
