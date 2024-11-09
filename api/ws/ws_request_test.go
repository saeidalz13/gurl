package ws

import (
	"fmt"
	"strings"
	"testing"
)

func createExpectedRequest(key string) string {
	sb := strings.Builder{}
	sb.Grow(50)

	sb.WriteString("GET / HTTP/1.1\r\n")
	sb.WriteString("Host: echo.websocket.org\r\nUser-Agent: gurl/1.2.0\r\nConnection: Upgrade\r\n")
	sb.WriteString("Upgrade: websocket\r\n")
	sb.WriteString(fmt.Sprintf("Sec-Websocket-Key: %s\r\n", key))
	sb.WriteString("Sec-WebSocket-Version: 13\r\n\r\n")

	return sb.String()
}

func TestCreateWsRequest(t *testing.T) {
	key, err := GenerateSecWebSocketKey()
	if err != nil {
		t.Fatal(err)
	}

	expectedRequest := createExpectedRequest(key)

	gotRequest := GenerateWebSocketRequest("echo.websocket.org", "/", key)

	if expectedRequest != gotRequest {
		t.Fatal("expected request did not match got request")
		fmt.Printf("Expected:\n----\n%s\n\n", expectedRequest)
		fmt.Printf("Got:\n----\n%s\n", gotRequest)
	}
}
