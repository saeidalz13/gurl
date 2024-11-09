package ws

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomKey(key []byte) (int, error) {
	n, err := rand.Read(key)
	if err != nil {
		return -1, err
	}
	return n, nil
}

func generateSecWsKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
