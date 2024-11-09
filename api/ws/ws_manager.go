package ws

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/terminalutils"
)

type WebSocketRequestGenerator struct {
	verbose  bool
	path     string
	domain   string
	ip       string
	secWsKey string
}

func NewWebSocketRequestGenerator(domain, path, ip string, verbose bool) WebSocketRequestGenerator {
	return WebSocketRequestGenerator{
		domain:  domain,
		path:    path,
		ip:      ip,
		verbose: verbose,
	}
}

func (w *WebSocketRequestGenerator) generateSecWsKey() {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	errutils.CheckErr(err)

	w.secWsKey = base64.StdEncoding.EncodeToString(key)
}

func (w WebSocketRequestGenerator) createWsRequest() string {
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

	sb.WriteString("User-Agent: gurl/1.0.0\r\n")

	// To signal the server this should be ws
	sb.WriteString("Connection: Upgrade")
	sb.WriteString("\r\n")

	sb.WriteString("Upgrade: websocket")
	sb.WriteString("\r\n")

	sb.WriteString("Sec-Websocket-Key: ")
	sb.WriteString(w.secWsKey)
	sb.WriteString("\r\n")

	sb.WriteString("Sec-WebSocket-Version: 13\r\n")

	// Ending of request based on HTTP
	sb.WriteString("\r\n")

	return sb.String()
}

func (w WebSocketRequestGenerator) Generate() (string, string) {
	w.generateSecWsKey()
	wsRequest := w.createWsRequest()

	if w.verbose {
		terminalutils.PrintWebSocketClientInfo(w.ip, wsRequest)
	}

	return w.secWsKey, wsRequest
}
