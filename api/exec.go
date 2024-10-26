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
	hwp := initCli()

	// Fetching IP and port of the remote address.
	ip, port, isConnTls := newRemoteAddrManager(ipCacheDir, hwp.domain).resolveConnectionInfo()

	// Initializing the TCP connection manager for
	// TCP conn and its management.
	tcm := newTCPConnManager(ip, port, isConnTls)
	err := tcm.initTCPConn(hwp)
	errutils.CheckErr(err)

	if hwp.isWs {
		wsRequest, err := createWsRequest(hwp.path, hwp.domain)
		errutils.CheckErr(err)
		go tcm.readWebSocketData()
		tcm.writeWebSocketData([]byte(wsRequest))

	} else {
		httpRequest := newHTTPRequestCreator(hwp.domain, hwp.path, hwp.method, hwp.headers).create()
		respBytes := tcm.dispatchHTTPRequest(httpRequest)
		newHTTPResponseParser(respBytes).parse().printPretty(hwp.verbose)
	}
}
