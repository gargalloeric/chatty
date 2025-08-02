package identity

import (
	"crypto/rand"
	"encoding/base32"
)

func GenerateRandomID(numberBytes int) (string, error) {
	buffer := make([]byte, numberBytes)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	id := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buffer)

	return id, nil
}
