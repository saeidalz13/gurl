package ws

import (
	"strings"
)

const (
	wsKeyLength = 16
)

type WebSocketRequestGenerator struct {
	path   string
	domain string
}

func NewWebSocketRequestGenerator(domain, path string) WebSocketRequestGenerator {
	return WebSocketRequestGenerator{
		domain: domain,
		path:   path,
	}
}

func (w WebSocketRequestGenerator) createWsRequest(secWsKey string) string {
	sb := strings.Builder{}
	sb.Grow(50)

	// Method (always GET for ws)
	sb.WriteString("GET")
	sb.WriteString(" ")

	sb.WriteString(w.path)
	sb.WriteString(" ")

	sb.WriteString("HTTP/1.1\r\n")

	sb.WriteString("Host: ")
	sb.WriteString(w.domain)
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

func (w WebSocketRequestGenerator) Generate() (string, string, error) {
	key := make([]byte, wsKeyLength)
	n, err := generateRandomKey(key)
	if err != nil {
		return "", "", err
	}

	secWsKey := generateSecWsKey(key[:n])
	wsRequest := w.createWsRequest(secWsKey)

	return secWsKey, wsRequest, nil
}
