package ws

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
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

func GenerateWebSocketRequest(domain, path, secWsKey string) string {
	sb := strings.Builder{}
	sb.Grow(50)

	// Method (always GET for ws)
	sb.WriteString("GET")
	sb.WriteString(" ")

	sb.WriteString(path)
	sb.WriteString(" ")

	sb.WriteString("HTTP/1.1\r\n")

	sb.WriteString("Host: ")
	sb.WriteString(domain)
	sb.WriteString("\r\n")

	sb.WriteString("User-Agent: gurl/1.2.0\r\n")

	// To signal the server this should be ws
	sb.WriteString("Connection: Upgrade")
	sb.WriteString("\r\n")

	sb.WriteString("Upgrade: websocket")
	sb.WriteString("\r\n")

	sb.WriteString("Sec-Websocket-Key: ")
	sb.WriteString(secWsKey)
	sb.WriteString("\r\n")

	sb.WriteString("Sec-WebSocket-Version: 13\r\n")

	// Ending of request based on HTTP
	sb.WriteString("\r\n")

	return sb.String()
}
