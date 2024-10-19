package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
)

func writeToTLSConn(tlsConn *tls.Conn, b []byte) error {
	_, err := tlsConn.Write(b)
	return err
}

func readFromTLSConn(tlsConn *tls.Conn) []byte {
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
		os.Exit(1)
	}

	return buf[:n]
}

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

func initTlsConn(ip, domain string) *tls.Conn {
	certPool := mustPrepareCertPool()

	tlsConn, err := tls.Dial(
		"tcp",
		ip+":"+httpconstants.PortHTTPS,
		&tls.Config{RootCAs: certPool, ServerName: domain},
	)
	errutils.CheckErr(err)

	tlsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	tlsConn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	return tlsConn
}
