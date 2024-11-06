package api

import (
	"github.com/saeidalz13/gurl/api/cli"
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
	"github.com/saeidalz13/gurl/internal/methodparser"
	"github.com/saeidalz13/gurl/internal/terminalutils"
)

func ExecGurl() {
	ipCacheDir := appconstants.MustMakeIpCacheDir()

	cp := cli.InitCli()

	method, err := methodparser.ParseMethod(cp.Method)
	errutils.CheckErr(err)

	dp := domainparser.NewDomainParser(cp.Domain)
	err = dp.Parse()
	errutils.CheckErr(err)

	ram := newRemoteAddrManager(ipCacheDir, dp.Domain, dp.DomainSegment)
	connInfo := ram.resolveConnectionInfo()

	tcm := newTCPConnManager(connInfo, dp.Domain)
	err = tcm.InitTCPConn()
	errutils.CheckErr(err)

	if dp.IsWebSocket {
		manageWebSocket(dp, tcm, cp.Verbose)
		return
	}

	contentType := determineContentType(cp.DataType)

	httpRequest := NewHTTPRequestGenerator(
		dp.Domain,
		dp.Path,
		cp.Cookies,
		method,
		contentType,
		cp.Data,
	).Generate()

	if cp.Verbose {
		terminalutils.PrintHTTPClientInfo(connInfo.ip.String(), httpRequest)
	}

	respBytes := tcm.dispatchHTTPRequest(httpRequest)
	newHTTPResponseParser(respBytes).parse().printPretty(cp.Verbose)
}

func manageWebSocket(dp domainparser.DomainParser, tcm TCPConnManager, verbose bool) {
	secWsKey, err := generateSecWsKey()
	errutils.CheckErr(err)

	wsRequest := createWsRequest(dp.Path, dp.Domain, secWsKey)

	if verbose {
		terminalutils.PrintWebSocketClientInfo(tcm.connInfo.ip.String(), wsRequest)
	}

	go tcm.readWebSocketData(secWsKey, verbose)
	tcm.writeWebSocketData([]byte(wsRequest))
}

func determineContentType(dataType uint8) string {
	switch dataType {
	case httpconstants.DataTypeJson:
		return "application/json"
	case httpconstants.DataTypeText:
		return "text/plain"
	case httpconstants.DataTypeImage:
		return "image/jpeg"
	}

	return ""
}
