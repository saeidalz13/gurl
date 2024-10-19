package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
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
