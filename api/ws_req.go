package api

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func generateSecWsKey() (string, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	return base64.StdEncoding.EncodeToString(key), err
}

func createWsRequest(path, domain, secWsKey string) string {
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

	sb.WriteString("User-Agent: gurl/1.0.0\r\n")

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
