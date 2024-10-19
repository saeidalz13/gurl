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

func mustParseDomain() string {
	domain := os.Args[1]
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)

	return domain
}

func mustParseMethod(rawMethod string) string {
	method := strings.TrimSpace(rawMethod)
	method = strings.ToUpper(method)

	for _, vm := range httpconstants.ValidHttpMethods {
		if vm == method {
			return method
		}
	}

	fmt.Printf("invalid method: %s\n", method)
	os.Exit(1)
	return "" // to silence the compiler
}

func initGurlCli() gurlCli {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
    
	if len(os.Args) < 2 {
        fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}
    
    domainCmd.Parse(os.Args[2:])

	return gurlCli{
		domain: mustParseDomain(),
		method: mustParseMethod(*methodPtr),
	}
}
