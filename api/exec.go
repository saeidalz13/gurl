package api

import (
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/methodparser"
)

// Entry point of the application execution.
//
// Depending on the url protocol, it can be
// a secure or insecure request.
func ExecGurl() {
	ipCacheDir := appconstants.MustMakeIpCacheDir()

	// Input params from CLI
	cp := initCli()

	method, err := methodparser.ParseMethod(cp.method)
	errutils.CheckErr(err)

	// Parsing domain
	dp := domainparser.NewDomainParser(cp.domain)
	err = dp.Parse()
	errutils.CheckErr(err)

	// Fetching IP and port of the remote address.
	ram := newRemoteAddrManager(ipCacheDir, dp.Domain, dp.DomainSegment)
	ip, port, isConnTls := ram.resolveConnectionInfo()

	// Initializing the TCP connection manager for
	// TCP conn and its management.
	tcm := newTCPConnManager(ip, port, isConnTls, dp.Domain)
	err = tcm.InitTCPConn()
	errutils.CheckErr(err)

	if hwp.isWs {
		wsRequest, err := createWsRequest(hwp.path, hwp.domain)
		errutils.CheckErr(err)
		go tcm.readWebSocketData()
		tcm.writeWebSocketData([]byte(wsRequest))

	} else {
		httpRequest := newHTTPRequestCreator(hwp.domain, hwp.path, method, hwp.headers).create()
		respBytes := tcm.dispatchHTTPRequest(httpRequest)
		newHTTPResponseParser(respBytes).parse().printPretty(hwp.verbose)
	}
}
