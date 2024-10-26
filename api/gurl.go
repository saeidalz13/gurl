package api

import (
	"flag"
	"fmt"
	"os"
)

type cliParams struct {
	domain  string
	method  string
	verbose bool
	ctJson  bool

	// path    string
	// headers []string
	// isWs       bool
	// wsProtocol uint8
}

// func (hwp *cliParams) mustParseDomainHTTP(domain string) {
// 	domain, err := stringutils.TrimDomainPrefix(domain)
// 	errutils.CheckErr(err)
// 	hwp.domain, hwp.path = stringutils.ExtractPathFromDomain(domain)
// }

// func (hwp *cliParams) mustParseDomainWS(domain string) {
// 	wsProtocol, trimmedDomain, err := stringutils.ExtractWsProtocol(domain)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	hwp.wsProtocol = wsProtocol
// 	hwp.domain, hwp.path = stringutils.ExtractPathFromDomain(trimmedDomain)
// }

// func (hwp *cliParams) mustParseMethod(rawMethod string) {
// 	method := strings.TrimSpace(rawMethod)
// 	method = strings.ToUpper(method)

// 	_, ok := httpconstants.ValidHttpMethods[method]
// 	if !ok {
// 		fmt.Printf("invalid method: %s\n", method)
// 		os.Exit(1)
// 	}

// 	hwp.method = method
// }

// func (hwp *cliParams) parseDomain() {
// 	domain := os.Args[1]
// 	domain = stringutils.TrimDomainSpace(domain)

// 	// Decide if the input wants a websocket req.
// 	// User MUST include the procotol in the input
// 	// domain for websocket requests.
// 	// (ws:// or wss://)
// 	hwp.isWs = stringutils.IsDomainForWebsocket(domain)

// 	if hwp.isWs {
// 		hwp.mustParseDomainWS(domain)
// 	} else {
// 		hwp.mustParseDomainHTTP(domain)
// 	}
// }

// func (hwp *cliParams) adjustHeaders(contentLengthJson bool) {
// 	if contentLengthJson {
// 		hwp.headers = append(hwp.headers, "Content-Type: application/json")
// 	}
// }

func initCli() cliParams {
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

	return cliParams{
		domain:  os.Args[1],
		method:  *methodPtr,
		verbose: *verbose,
		ctJson:  *ctJsonPtr,
	}
}
