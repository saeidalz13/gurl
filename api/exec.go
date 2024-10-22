package api

import (
	"crypto/tls"
	"fmt"
	"strings"
)

func createHTTPRequest(gp gurlParams) string {
	sb := strings.Builder{}
	sb.Grow(50)

	// Method
	sb.WriteString(gp.method)
	sb.WriteString(" ")

	// Path
	sb.WriteString(gp.path)
	sb.WriteString(" ")

	// Protocol and version
	sb.WriteString("HTTP/1.1\r\n")

	// Host
	sb.WriteString("Host: ")
	sb.WriteString(gp.domain)
	sb.WriteString("\r\n")

	// User-Agent
	sb.WriteString("User-Agent: gurl/1.0.0\r\n")

	// Accept headers
	sb.WriteString("Accept: */*\r\n")

	// User Headers
	for _, header := range gp.headers {
		sb.WriteString(header)
		sb.WriteString("\r\n")
	}

	// Ending of request based on HTTP
	sb.WriteString("\r\n")

	return sb.String()
}

func execHTTPReq(tlsConn *tls.Conn, httpRequest string) {
	if err := writeToTLSConn(tlsConn, []byte(httpRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		return
	}

	tcpRespBytes := readFromTLSConn(tlsConn)
	tcpResponse := string(tcpRespBytes)

	httpResp := newHTTPResponse(tcpResponse)
	printPretty(httpResp)
}

func ExecGurl() {
	gp := initGurlCli()
	ip := mustFetchDomainIp(gp.domain)

	tlsConn := initTlsConn(ip.String(), gp.domain)
	defer tlsConn.Close()

	httpRequest := createHTTPRequest(gp)

	execHTTPReq(tlsConn, httpRequest)
}
