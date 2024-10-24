package api

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
	"github.com/saeidalz13/gurl/internal/stringutils"
)

type httpWsParams struct {
	// Http
	domain  string
	method  string
	path    string
	headers []string
	verbose bool

	// Is this a WebSocket request
	isWs       bool
	wsProtocol uint8
}

func (hwp *httpWsParams) mustParseDomainHTTP(domain string) {
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)
	hwp.domain, hwp.path = stringutils.ExtractPathFromDomain(domain)
}

func (hwp *httpWsParams) mustParseDomainWS(domain string) {
	wsProtocol, trimmedDomain, err := stringutils.ExtractWsProtocol(domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	hwp.wsProtocol = wsProtocol
	hwp.domain, hwp.path = stringutils.ExtractPathFromDomain(trimmedDomain)
}

func (hwp *httpWsParams) mustParseMethod(rawMethod string) {
	method := strings.TrimSpace(rawMethod)
	method = strings.ToUpper(method)

	_, ok := httpconstants.ValidHttpMethods[method]
	if !ok {
		fmt.Printf("invalid method: %s\n", method)
		os.Exit(1)
	}

	hwp.method = method
}

func (hwp *httpWsParams) parseDomain() {
	domain := os.Args[1]
	domain = stringutils.TrimDomainSpace(domain)

	// Decide if the input wants a websocket req.
	// User MUST include the procotol in the input
	// domain for websocket requests.
	// (ws:// or wss://)
	hwp.isWs = stringutils.IsDomainForWebsocket(domain)

	if hwp.isWs {
		hwp.mustParseDomainWS(domain)
	} else {
		hwp.mustParseDomainHTTP(domain)
	}
}

func (hwp *httpWsParams) adjustHeaders(contentLengthJson bool) {
	if contentLengthJson {
		hwp.headers = append(hwp.headers, "Content-Type: application/json")
	}
}

func newHTTPWSParams() httpWsParams {
	// CLI init section
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	ctJsonPtr := domainCmd.Bool("json", false, "Add HTTP Request Header -> Content-Type: application/json")
	verbose := domainCmd.Bool("v", false, "Prints metadata and steps of request")
	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}
	domainCmd.Parse(os.Args[2:])

	// Initialize gurlParams struct
	hwp := httpWsParams{
		headers: make([]string, 0),
		verbose: *verbose,
	}

	hwp.parseDomain()
	hwp.mustParseMethod(*methodPtr)
	hwp.adjustHeaders(*ctJsonPtr)

	return hwp
}
