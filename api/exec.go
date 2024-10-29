package api

import (
	"github.com/saeidalz13/gurl/internal/appconstants"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
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
		secWsKey, err := generateSecWsKey()
		errutils.CheckErr(err)

		wsRequest := createWsRequest(dp.Path, dp.Domain, secWsKey)

		go tcm.readWebSocketData(secWsKey)
		tcm.writeWebSocketData([]byte(wsRequest))
		return
	}

	contentType := determineContentType(cp.dataType)

	hrg := NewHTTPRequestGenerator(dp.Domain, dp.Path, cp.cookies)
	var httpRequest string
methodBlock:
	switch method {
	case httpconstants.MethodGET:
		httpRequest = hrg.GenerateGETRequest()

	case httpconstants.MethodPOST:
		httpRequest = hrg.GeneratePOSTRequest(cp.data, contentType)
		break methodBlock

	case httpconstants.MethodPUT, httpconstants.MethodPATCH:
		httpRequest = hrg.GeneratePUTPATCHRequest(cp.data, contentType)
		break methodBlock

	case httpconstants.MethodDELETE:
		httpRequest = hrg.GenerateDELETERequest()
	}

	respBytes := tcm.dispatchHTTPRequest(httpRequest)
	newHTTPResponseParser(respBytes).parse().printPretty(cp.verbose)
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
