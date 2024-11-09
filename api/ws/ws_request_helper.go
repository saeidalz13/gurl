package ws

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	wsKeyLength = 16
)

func generateRandomKey(key []byte) (int, error) {
	n, err := rand.Read(key)
	if err != nil {
		return -1, err
	}
	return n, nil
}

func GenerateSecWebSocketKey() (string, error) {
	key := make([]byte, wsKeyLength)

	n, err := generateRandomKey(key)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(key[:n]), nil
}
