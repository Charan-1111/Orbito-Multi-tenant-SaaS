package utils

import (
	"encoding/base64"
	"fmt"
)

func Base64Decode(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	return string(decodedBytes), nil
}
