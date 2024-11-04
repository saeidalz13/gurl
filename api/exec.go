package api

import (
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/methodparser"
	"github.com/saeidalz13/gurl/internal/terminalutils"
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
		manageWebSocket(dp, tcm)
		return
	}

	contentType := determineContentType(cp.dataType)

	httpRequest := NewHTTPRequestGenerator(
		dp.Domain,
		dp.Path,
		cp.cookies,
		method,
		contentType,
		cp.data,
	).Generate()

	if cp.verbose {
		terminalutils.PrintClientInfo(ip.String(), httpRequest)
	}

	respBytes := tcm.dispatchHTTPRequest(httpRequest)
	newHTTPResponseParser(respBytes).parse().printPretty(cp.verbose)
}

func manageWebSocket(dp domainparser.DomainParser, tcm TCPConnManager) {
	secWsKey, err := generateSecWsKey()
	errutils.CheckErr(err)

	wsRequest := createWsRequest(dp.Path, dp.Domain, secWsKey)

	go tcm.readWebSocketData(secWsKey)
	tcm.writeWebSocketData([]byte(wsRequest))
}

func determineContentType(dataType uint8) string {
	switch dataType {
	case dataTypeJson:
		return "application/json"
	case dataTypeText:
		return "text/plain"
	case dataTypeImage:
		return "image/jpeg"
	}

	return ""
}
