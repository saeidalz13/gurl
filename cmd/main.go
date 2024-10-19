package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/saeidalz13/gurl/api"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/stringutils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		os.Exit(1)
	}

	domainCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	_ = domainCmd.String("method", "GET", "HTTP method")
	domainCmd.Parse(os.Args[2:])

	domain := os.Args[1]
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)

	destIp := api.MustFetchDomainIp(domain)
	api.ExecGetHttpReq(destIp, domain)
}
