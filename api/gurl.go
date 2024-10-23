package api

import (
	"flag"
	"fmt"
	"net"
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

	// Is this a WebSocket request
	isWs       bool
	wsProtocol uint8

	// Aesthetics
	pretty bool

	// Server info
	serverIP net.IP
	port     int

	//
	isConnTls bool
}

func (gp *gurlParams) mustParseDomainHTTP() {
	domain := os.Args[1]
	domain = stringutils.TrimDomainSpace(domain)
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)
	gp.domain, gp.path = stringutils.ExtractPathFromDomain(domain)
}

func (gp *gurlParams) mustParseDomainWS() {
	domain := os.Args[1]
	domain = stringutils.TrimDomainSpace(domain)
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

func newGurlParams() gurlParams {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	ctJsonPtr := domainCmd.Bool("json", false, "Set Header -> Content-Type: application/json")
	isPrettyPtr := domainCmd.Bool("pretty", false, "Print the response in an organized manner")
	isWs := domainCmd.Bool("ws", false, "Makes a WebSocket request")
	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}
	domainCmd.Parse(os.Args[2:])

	gp := gurlParams{
		headers: make([]string, 0),
		isWs:    *isWs,
		pretty:  *isPrettyPtr,
	}

	if *ctJsonPtr {
		gp.headers = append(gp.headers, "Content-Type: application/json")
	}

	if gp.isWs {
		gp.mustParseDomainWS()
	} else {
		gp.mustParseDomainHTTP()
	}

	gp.mustParseMethod(*methodPtr)

	if stringutils.IsDomainLocalHost(gp.domain) {
		gp.serverIP = net.IPv4(127, 0, 0, 1)
		port, err := stringutils.ExtractPort(gp.domain)
		errutils.CheckErr(err)
		gp.port = port
		gp.isConnTls = false
	} else {
		ip := mustFetchDomainIp(gp.domain)
		gp.serverIP = ip
		gp.isConnTls = true
	}

	return gp
}
