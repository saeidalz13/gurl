package api

import (
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/methodparser"
)

func ExecGurl() {
	ipCacheDir := appconstants.MustMakeIpCacheDir()

	cp := initCli()

	method, err := methodparser.ParseMethod(cp.method)
	errutils.CheckErr(err)

	dp := domainparser.NewDomainParser(cp.domain)
	err = dp.Parse()
	errutils.CheckErr(err)

	ram := newRemoteAddrManager(ipCacheDir, dp.Domain, dp.DomainSegment)
	ip, port, isConnTls := ram.resolveConnectionInfo()

	tcm := newTCPConnManager(ip, port, isConnTls, dp.Domain)
	err = tcm.InitTCPConn()
	errutils.CheckErr(err)

	if dp.IsWebSocket {
		wsRequest, err := createWsRequest(dp.Path, dp.Domain)
		errutils.CheckErr(err)
		go tcm.readWebSocketData()
		tcm.writeWebSocketData([]byte(wsRequest))
		return
	}

	hrc := newHTTPRequestCreator(dp.Domain, dp.Path, method)
	if cp.ctJson {
		hrc.AddContentTypeJson()
	}
	if cp.cookies != "" {
		hrc.AddCookie(cp.cookies)
	}

	httpRequest := hrc.Create()
	respBytes := tcm.dispatchHTTPRequest(httpRequest)
	newHTTPResponseParser(respBytes).parse().printPretty(cp.verbose)
}
