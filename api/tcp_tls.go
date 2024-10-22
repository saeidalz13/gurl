package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
)

const (
	HeaderCloseConn uint8 = iota
	HeaderChunk
	HeaderContentLength
)

func writeToTLSConn(tlsConn *tls.Conn, b []byte) error {
	_, err := tlsConn.Write(b)
	return err
}

// Identifies which parameter exists in the
// http response header so we know if we should
// close the connection right away or keep it
// alive and handle the rest of streams of data.
func identifyHeaderParam(bufContainingHeader []byte) uint8 {
	if bytes.Contains(bufContainingHeader, []byte("Content-Length")) {
		return HeaderContentLength
	}

	if bytes.Contains(bufContainingHeader, []byte("chunked")) {
		return HeaderChunk
	}

	return HeaderCloseConn
}

// in HTTP 1.1 the default header setting of
// connection is "keep-alive".
//
// It is the case that if the server wants to send
// the response in different streams of tcp, it
// indiciates either "Content-Length" or sets the
// "Transfer-Encoding" to "chunked"
//
// If none of those options provided, that means the
// entire data is sent in one single stream and we can
// close the TCP conn after the first read.
func readFromTLSConn(tlsConn *tls.Conn) []byte {
	var response bytes.Buffer
	var readContentLength int = -1
	var contentLength int
	var headerParam uint8

	bufSize := 2 << 12 // 4kb
	var readIteration int
	headerIteration := 1

tcpLoop:
	for readContentLength < contentLength {
		readIteration++
		buf := make([]byte, bufSize)

		n, err := tlsConn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			if strings.Contains(err.Error(), "i/o timeout") {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("error read: %v\n", err)
			os.Exit(1)
		}

		response.Write(buf[:n])

		if readIteration == headerIteration {
			headerParam = identifyHeaderParam(buf[:n])
		}

		switch headerParam {
		case HeaderCloseConn:
			return response.Bytes()

		case HeaderContentLength:
			if readIteration == headerIteration {
				var bodyPos, shouldBreakNum int
				bytesLines := bytes.Split(buf[:n], []byte("\r\n"))
			lineLoop:
				for _, line := range bytesLines {
					// + 2 is for \n and \r
					bodyPos = bodyPos + len(line) + 2
					if string(line) == "" {
						shouldBreakNum++
					}

					if bytes.Contains(line, []byte("Content-Length")) {
						contentLengthBytes := bytes.TrimSpace(bytes.Split(line, []byte(":"))[1])

						num, err := strconv.ParseInt(string(contentLengthBytes), 10, 64)
						if err != nil {
							fmt.Println("invalid response from server")
							os.Exit(1)
						}
						contentLength = int(num)
						shouldBreakNum++
					}

					if shouldBreakNum == 2 {
						break lineLoop
					}
				}
				readContentLength = n - bodyPos

				continue tcpLoop
			} else {
				readContentLength += n
			}

		case HeaderChunk:
			bytesLines := bytes.Split(buf[:n], []byte("\r\n"))
			for _, line := range bytesLines {

				// If "0" was found at the end of the body
				// it shows that there's no more bytes to
				// be sent from the server.
				if bytes.Equal(line, []byte{48}) {
					return response.Bytes()
				}
			}

			// This is unncessary but to show explictely
			// what needs to happen. If "0" was not found
			// at the end of body, it means more data
			// will be streamed from server. So tcpLoop
			// shall live on!
			continue tcpLoop
		}
	}

	return response.Bytes()
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
