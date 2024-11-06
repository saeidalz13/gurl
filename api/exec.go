package api

import (
	"github.com/saeidalz13/gurl/api/cli"
	"github.com/saeidalz13/gurl/api/http"
	"github.com/saeidalz13/gurl/api/tcp"
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

	tcm := tcp.NewTCPConnManager(connInfo, dp.Domain)
	err = tcm.InitTCPConn()
	errutils.CheckErr(err)

	if dp.IsWebSocket {
		manageWebSocket(dp, tcm, cp.Verbose, connInfo.IP.String())
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
		terminalutils.PrintHTTPClientInfo(connInfo.IP.String(), httpRequest)
	}

	respBytes := tcm.DispatchHTTPRequest(httpRequest)
	http.NewHTTPResponseParser(respBytes).Parse().Print(cp.Verbose)
}

func manageWebSocket(dp domainparser.DomainParser, tcm tcp.TCPConnManager, verbose bool, ip string) {
	secWsKey, err := generateSecWsKey()
	errutils.CheckErr(err)

	wsRequest := createWsRequest(dp.Path, dp.Domain, secWsKey)

	if verbose {
		terminalutils.PrintWebSocketClientInfo(ip, wsRequest)
	}

	go tcm.ReadWebSocketData(secWsKey, verbose)
	tcm.WriteWebSocketData([]byte(wsRequest))
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
