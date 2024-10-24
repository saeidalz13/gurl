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
	ip, port, isConnTls := newRemoteAddrManager(ipCacheDir, gp.domain).determineRemoteIpPort()

	// Initializing the TCP connection manager for
	// TCP conn and its management.
	tcm := newTCPConnManager(ip, port, isConnTls)
	err := tcm.initTCPConn(gp)
	errutils.CheckErr(err)

	if gp.isWs {
		wsRequest, err := createWsRequest(gp.path, gp.domain)
		errutils.CheckErr(err)
		go tcm.manageReadTCPConnWS()
		tcm.manageWriteTCPConnWS([]byte(wsRequest))
	} else {
		tcm.setDeadlineToConn()
		httpRequest := createHTTPRequest(gp)
		respBytes := tcm.makeHTTPRequest(httpRequest)
		newHTTPResponse(respBytes).parse().printPretty(gp.verbose)
	}
}
