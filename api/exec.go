package api

import (
	"crypto/tls"
	"fmt"
	"os"
)

func handleHTTPSReq(conn *tls.Conn, httpRequest string) {
	if err := writeToTLSConn(conn, []byte(httpRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		os.Exit(1)
	}

	tcpRespBytes := readFromTLSConnHTTPS(conn)
	httpResp := newHTTPResponse(string(tcpRespBytes))
	printPretty(httpResp)
}

func handleWSSReq(conn *tls.Conn, wsRequest string) {
	if err := writeToTLSConn(conn, []byte(wsRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		os.Exit(1)
	}

	readFromTLSConnWSS(conn)
}

func execSecure(gp gurlParams) {
	tlsConn := makeTlsTcpConn(gp.serverIP.String(), gp.domain)
	defer tlsConn.Close()

	if gp.isWs {
		wsRequest := mustCreateWsRequest(gp.path, gp.domain)
		handleWSSReq(tlsConn, wsRequest)
	} else {
		httpRequest := createHTTPRequest(gp)
		handleHTTPSReq(tlsConn, httpRequest)
	}
}

// For the requests that are not made with
// TLS handshake.
func execInSecure(gp gurlParams) {
	tcpConn := makeTcpConn(gp.serverIP, gp.port)
	defer tcpConn.Close()

	if gp.isWs {
		wsRequest := mustCreateWsRequest(gp.path, gp.domain)
		manageTCPConnWS(tcpConn, []byte(wsRequest))
	}
}

// Entry point of the application execution.
//
// Depending on the url protocol, it can be
// a secure or insecure request.
func ExecGurl() {
	gp := newGurlParams()

	if gp.isConnTls {
		execSecure(gp)
	} else {
		execInSecure(gp)
	}
}
