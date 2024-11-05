package api

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
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

//go:embed cacert.pem
var cacertsPEM []byte

type TCPConnManager struct {
	domain   string
	connInfo ConnInfo
	conn     net.Conn
}

func newTCPConnManager(connInfo ConnInfo, domain string) TCPConnManager {
	return TCPConnManager{
		connInfo: connInfo,
		domain:   domain,
	}
}

func (tcm TCPConnManager) setDeadlineToConn() {
	tcm.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	tcm.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
}

// For the requests that are not made with
// TLS handshake.
func (tcm *TCPConnManager) InitTCPConn() error {
	if tcm.connInfo.isTls {
		addr := fmt.Sprintf("%s:%d", tcm.connInfo.ip.String(), tcm.connInfo.port)
		if tcm.connInfo.ipType == 1 {
			addr = fmt.Sprintf("[%s]:%d", tcm.connInfo.ip.String(), tcm.connInfo.port)
		}

		certPool := mustPrepareCertPool()
		conn, err := tls.Dial(
			"tcp",
			addr,
			&tls.Config{RootCAs: certPool, ServerName: tcm.domain},
		)
		if err != nil {
			return err
		}
		tcm.conn = conn

		return nil
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: tcm.connInfo.ip, Port: tcm.connInfo.port})
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
				cl, bodyIdx := extractContentLengthBodyStartingIdx(buf[:n])
				contentLength = cl
				readContentLength = n - bodyIdx

				continue tcpLoop
			} else {
				readContentLength += n
			}

		case headerChunk:
			// bytes of final \r\n + bytes of end of body \r\n + offset of len
			// 2 + 2 + 1 = 5
			// position of potential '0' = 5
			if buf[n-5] == '0' {
				return response.Bytes()
			}

			// unncessary `continue` but to show explicitely
			// what needs to happen. If '0' was not found
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
func (tcm TCPConnManager) readWebSocketData(secWsKey string, verbose bool) {
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
				fmt.Println("Server did not accept WS req:")
				fmt.Println(respHeader)
				os.Exit(1)
			}
			headerIteration = false
			if verbose {
				fmt.Println(respHeader)
			}

			if !isServerVerified(buf[:n], secWsKey) {
				fmt.Println("server key not verified")
				break
			}

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

// The initial msgByte is the request sent to the
// server to initiate the WS connection.
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
	certPool, err := x509.SystemCertPool()
	errutils.CheckErr(err)

	if !certPool.AppendCertsFromPEM(cacertsPEM) {
		fmt.Printf("failed to load the certificates")
		os.Exit(1)
	}

	return certPool
}

func extractContentLengthBodyStartingIdx(httpResp []byte) (int, int) {
	var bodyPos, shouldBreakNum, contentLength int
	bytesLines := bytes.Split(httpResp, []byte("\r\n"))
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

	return contentLength, bodyPos
}

// When the WebSocket server sends the 101 Code, it
// includes `Sec-Weboscket-Accept: VALUE`. `VALUE` is
// the base64 encoded value of SHA-1 hash of the
// client key + special GUID. This GUID is a unversal
// constant.
//
// This function checks the `VALUE` seny by the server
// with the base64 SHA-1 hash of client key. If they match
// the response was sent from the requested server and not
// and itermediary malicious middle man.
func isServerVerified(respHeaser []byte, key string) bool {
	specialGUID := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	var secWsAccept []byte

	for _, line := range bytes.Split(respHeaser, []byte("\r\n")) {
		lineSegments := bytes.Split(line, []byte(":"))

		if len(lineSegments) != 2 {
			continue
		}

		if bytes.Equal(bytes.ToLower(lineSegments[0]), []byte("sec-websocket-accept")) {
			secWsAccept = bytes.TrimSpace(lineSegments[1])
			break
		}
	}

	if secWsAccept == nil {
		return false
	}

	h := sha1.New()
	h.Write([]byte(key + specialGUID))
	hashed := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashed) == string(secWsAccept)
}
