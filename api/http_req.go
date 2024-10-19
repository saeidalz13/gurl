package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
)

func mustPrepareCertPool() *x509.CertPool {
	readBytes, err := os.ReadFile(os.Getenv("CERTS_DIR"))
	errutils.CheckErr(err)

	certPool, err := x509.SystemCertPool()
	errutils.CheckErr(err)

	if !certPool.AppendCertsFromPEM(readBytes) {
		fmt.Printf("failed to load the certificates")
		os.Exit(1)
	}

	return certPool
}

func ExecGetHttpReq(ip net.IP, domain string) {
	certPool := mustPrepareCertPool()

	tcpConn, err := tls.Dial(
		"tcp",
		ip.String()+":"+httpconstants.PortHTTPS,
		&tls.Config{RootCAs: certPool, ServerName: domain},
	)
	errutils.CheckErr(err)

	tcpConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	tcpConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	defer tcpConn.Close()

	for {
		_, err := tcpConn.Write([]byte("GET / HTTP/1.1\r\n" +
			"Host: " + domain + "\r\nUser-Agent: Client\r\nAccept: */*\r\nConnection: close\r\n\r\n"))

		if err != nil {
			fmt.Printf("error write: %v\n", err)
			break
		}

		buf := make([]byte, 2<<12)
		n, err := tcpConn.Read(buf)
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
