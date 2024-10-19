package api

import (
	"crypto/tls"
	"fmt"
)

func execHTTPReq(tlsConn *tls.Conn, httpRequest string) {
	if err := writeToTLSConn(tlsConn, []byte(httpRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		return
	}

	readBytes := readFromTLSConn(tlsConn)
	fmt.Println(string(readBytes))
}
