package api

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
)

func execGetHttpReq(tlsConn *tls.Conn, domain string) {
	for {
		_, err := tlsConn.Write([]byte("GET / HTTP/1.1\r\n" +
			"Host: " + domain + "\r\nUser-Agent: Client\r\nAccept: */*\r\nConnection: close\r\n\r\n"))

		if err != nil {
			fmt.Printf("error write: %v\n", err)
			break
		}

		buf := make([]byte, 2<<12)
		n, err := tlsConn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				os.Exit(0)
			}

			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("error read: %v\n", err)
			continue
		}

		fmt.Println(string(buf[:n]))
	}
}
