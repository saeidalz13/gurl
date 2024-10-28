package api

import (
	"flag"
	"fmt"
	"os"
)

type cliParams struct {
	verbose  bool
	jsonData string
	domain   string
	method   string
	cookies  string
}

func initCli() cliParams {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	jsonPtr := domainCmd.String("json", "", "Add json data to body")
	verbose := domainCmd.Bool("v", false, "Verbose run")
	cookies := domainCmd.String("cookies", "", "Add cookie to request header; e.g. -cookies='name1=value1; name2=value2'")

	help := flag.Bool("h", false, "gURL usage")
	flag.Parse()

	if *help {
		domainCmd.Usage()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}

	domainCmd.Parse(os.Args[2:])

	return cliParams{
		domain:   os.Args[1],
		method:   *methodPtr,
		verbose:  *verbose,
		jsonData: *jsonPtr,
		cookies:  *cookies,
	}
}
