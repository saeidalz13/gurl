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

type gurlParams struct {
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

func (gp *gurlParams) mustParseDomainHTTP(domain string) {
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)
	gp.domain, gp.path = stringutils.ExtractPathFromDomain(domain)
}

func (gp *gurlParams) mustParseDomainWS(domain string) {
	wsProtocol, trimmedDomain, err := stringutils.ExtractWsProtocol(domain)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gp.wsProtocol = wsProtocol
	gp.domain, gp.path = stringutils.ExtractPathFromDomain(trimmedDomain)
}

func (gp *gurlParams) mustParseMethod(rawMethod string) {
	method := strings.TrimSpace(rawMethod)
	method = strings.ToUpper(method)

	_, ok := httpconstants.ValidHttpMethods[method]
	if !ok {
		fmt.Printf("invalid method: %s\n", method)
		os.Exit(1)
	}

	gp.method = method
}

func (gp *gurlParams) parseDomain() {
	domain := os.Args[1]
	domain = stringutils.TrimDomainSpace(domain)

	// Decide if the input wants a websocket req.
	// User MUST include the procotol in the input
	// domain for websocket requests.
	// (ws:// or wss://)
	gp.isWs = stringutils.IsDomainForWebsocket(domain)

	if gp.isWs {
		gp.mustParseDomainWS(domain)
	} else {
		gp.mustParseDomainHTTP(domain)
	}
}

func (gp *gurlParams) adjustHeaders(contentLengthJson bool) {
	if contentLengthJson {
		gp.headers = append(gp.headers, "Content-Type: application/json")
	}
}

func newGurlParams() gurlParams {
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
	gp := gurlParams{
		headers: make([]string, 0),
		verbose: *verbose,
	}

	gp.parseDomain()
	gp.mustParseMethod(*methodPtr)
	gp.adjustHeaders(*ctJsonPtr)

	return gp
}
