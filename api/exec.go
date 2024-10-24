package api

import (
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/errutils"
)

// Entry point of the application execution.
//
// Depending on the url protocol, it can be
// a secure or insecure request.
func ExecGurl() {
	ipCacheDir := appconstants.MustMakeIpCacheDir()

	// Preparing the parameters for gurl app.
	gp := newGurlParams()

	// Fetching IP and port of the remote address.
	ip, port, isConnTls := newRemoteAddrManager(ipCacheDir, gp.domain).resolveConnectionInfo()

	// Initializing the TCP connection manager for
	// TCP conn and its management.
	tcm := newTCPConnManager(ip, port, isConnTls)
	err := tcm.initTCPConn(gp)
	errutils.CheckErr(err)

	if gp.isWs {
		wsRequest, err := createWsRequest(gp.path, gp.domain)
		errutils.CheckErr(err)
		go tcm.readWebSocketData()
		tcm.writeWebSocketData([]byte(wsRequest))

	} else {
		httpRequest := newHTTPRequestCreator(gp.domain, gp.path, gp.method, gp.headers).create()
		respBytes := tcm.dispatchHTTPRequest(httpRequest)
		newHTTPResponseParser(respBytes).parse().printPretty(gp.verbose)
	}
}
