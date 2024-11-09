package api

import (
	"github.com/saeidalz13/gurl/api/cli"
	"github.com/saeidalz13/gurl/api/conninfo"
	"github.com/saeidalz13/gurl/api/http"
	"github.com/saeidalz13/gurl/api/tcp"
	"github.com/saeidalz13/gurl/api/ws"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/methodparser"
	"github.com/saeidalz13/gurl/internal/pathutils"
	"github.com/saeidalz13/gurl/internal/terminalutils"
)

func ExecGurl() {
	ipCacheDir := pathutils.MustMakeIpCacheDir()

	cp := cli.InitCli()

	method, err := methodparser.ParseMethod(cp.Method)
	errutils.CheckErr(err)

	dp := domainparser.NewDomainParser(cp.Domain)
	err = dp.Parse()
	errutils.CheckErr(err)

	connInfo := conninfo.NewConnInfoResolver(
		ipCacheDir,
		dp.Domain,
		dp.DomainSegment,
	).Resolve()

	tcm := tcp.NewTCPConnManager(connInfo, dp.Domain)
	err = tcm.InitTCPConn()
	errutils.CheckErr(err)

	switch dp.Protocol {
	case domainparser.ProtocolWS:
		secWsKey, wsRequest := ws.NewWebSocketRequestGenerator(
			dp.Domain,
			dp.Path,
			connInfo.IP.String(),
			cp.Verbose,
		).Generate()

		go tcm.ReadWebSocketData(secWsKey, cp.Verbose)
		tcm.WriteWebSocketData([]byte(wsRequest))

	case domainparser.ProtocolHTTP:
		httpRequest := http.NewHTTPRequestGenerator(
			dp.Domain,
			dp.Path,
			cp.Cookies,
			method,
			cp.Data,
			cp.DataType,
		).Generate()

		if cp.Verbose {
			terminalutils.PrintHTTPClientInfo(connInfo.IP.String(), httpRequest)
		}

		respBytes := tcm.DispatchHTTPRequest(httpRequest)
		http.NewHTTPResponseParser(respBytes).Parse().Print(cp.Verbose)
	}
}
