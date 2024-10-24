package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/terminalutils"
	"github.com/saeidalz13/gurl/internal/wsutils"
)

const (
	headerCloseConn uint8 = iota
	headerChunk
	headerContentLength
)

type TCPConnManager struct {
	ip   net.IP
	port int

	isConnTls bool
	conn      net.Conn
}

func newTCPConnManager(ip net.IP, port int, isConnTls bool) TCPConnManager {
	return TCPConnManager{
		ip:        ip,
		port:      port,
		isConnTls: isConnTls,
	}
}

func (tcm TCPConnManager) setDeadlineToConn() {
	tcm.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	tcm.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
}

// For the requests that are not made with
// TLS handshake.
func (tcm *TCPConnManager) initTCPConn(hwp httpWsParams) error {
	if tcm.isConnTls {
		certPool := mustPrepareCertPool()
		conn, err := tls.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", tcm.ip.String(), tcm.port),
			&tls.Config{RootCAs: certPool, ServerName: hwp.domain},
		)
		if err != nil {
			return err
		}
		tcm.conn = conn
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: tcm.ip, Port: tcm.port})
	if err != nil {
		return err
	}
	tcm.conn = conn
	return nil
}

// Write the prepare http request to TCP connection
// and returns the response bytes.
func (tcm TCPConnManager) dispatchHTTPRequest(httpRequest string) []byte {
	tcm.setDeadlineToConn()

	_, err := tcm.conn.Write([]byte(httpRequest))
	if err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		os.Exit(1)
	}

	return tcm.readHTTPRespFromConn()
}

// func handleWSSReq(conn *tls.Conn, wsRequest string) {
// 	_, err := conn.Write([]byte(wsRequest))
// 	if err != nil {
// 		fmt.Printf("write tcp read: %v\n", err)
// 		os.Exit(1)
// 	}
// 	readFromTLSConnWSS(conn)
// }

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
func (tcm TCPConnManager) readHTTPRespFromConn() []byte {
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

		n, err := tcm.conn.Read(buf)
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
		case headerCloseConn:
			return response.Bytes()

		case headerContentLength:
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

		case headerChunk:
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

// Reads the content of websocket frame stream
// on a separate goroutine to be able to both
// read from and write to TCP conn concurrently.
func (tcm TCPConnManager) readWebSocketData() {
	headerIteration := true

	for {
		buf := make([]byte, 2<<15)
		n, err := tcm.conn.Read(buf)
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

		if headerIteration {
			respHeader := string(buf[:n])
			// 101 is code for switching protocol showing
			// that server is ready to be serving WS
			if !strings.Contains(respHeader, "101") {
				fmt.Println(respHeader)
				os.Exit(1)
			}
			headerIteration = false

		} else {
			payload, err := wsutils.ParseWsFrame(buf[:n])
			if err != nil {
				fmt.Println(err)
				continue
			}
			terminalutils.PrintWsServerMsg(string(payload))
		}
	}

	fmt.Println("Server closed connection")
	os.Exit(1)
}

func (tcm TCPConnManager) writeWebSocketData(msgByte []byte) {
	for {
		_, err := tcm.conn.Write(msgByte)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		input := terminalutils.GetWsInputFromStdin()
		terminalutils.PrintWsClientMsg(string(input))

		frame := wsutils.CreateWsFrame(input)
		msgByte = frame
	}
}

// Identifies which parameter exists in the
// http response header so we know if we should
// close the connection right away or keep it
// alive and handle the rest of streams of data.
func identifyHeaderParam(bufContainingHeader []byte) uint8 {
	if bytes.Contains(bufContainingHeader, []byte("Content-Length")) {
		return headerContentLength
	}

	if bytes.Contains(bufContainingHeader, []byte("chunked")) {
		return headerChunk
	}

	return headerCloseConn
}

// Preparing the certificaion info for
// the TLS handshake on TCP. Some systems
// don't automatically load certificates.
//
// It is included in the binary package.
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
