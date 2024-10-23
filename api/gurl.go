package api

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/saeidalz13/gurl/internal/appconstants"
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

	// Server info
	serverIP net.IP
	port     int

	isConnTls  bool
	ipCacheDir string
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
	ipCacheDir := appconstants.MustMakeIpCacheDir()

	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	ctJsonPtr := domainCmd.Bool("json", false, "Add HTTP Request Header -> Content-Type: application/json")
	verbose := domainCmd.Bool("v", false, "Prints metadata and steps of request")
	isWs := domainCmd.Bool("ws", false, "Makes a WebSocket request")
	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}
	domainCmd.Parse(os.Args[2:])

	gp := gurlParams{
		headers:    make([]string, 0),
		isWs:       *isWs,
		verbose:    *verbose,
		ipCacheDir: ipCacheDir,
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
		ip, err := fetchCachedIp(ipCacheDir, gp.domain)
		if err != nil {
			ip = mustFetchDomainIp(gp.domain)
			if err := cacheDomainIp(ipCacheDir, gp.domain, ip.String()); err != nil {
				fmt.Println("failed to cache ip")
			}
		}

		gp.serverIP = ip
		gp.isConnTls = true
	}

	return gp
}
