package ws

import (
	"testing"
)

const desiredRequest = `GET / HTTP/1.1\r\nHost: echo.websocket.org\r\nUser-Agent: gurl/1.2.0\r\nConnection: Upgrade\r\nUpgrade: websocket\r\nSec-Websocket-Key: liba47zQ7ldW1aOBJt7Mjw==\r\nSec-WebSocket-Version: 13\r\n\r\n`

// var testWsReqGenerator = WebSocketRequestGenerator{
// 	path:   "/",
// 	domain: "echo.websocket.org",
// }

func TestCreateWsRequest(t *testing.T) {
	// wsSecKey, wsReq, err := testWsReqGenerator.Generate()

	// if err != nil {
	//     t.Fatal(err)
	// }
}
