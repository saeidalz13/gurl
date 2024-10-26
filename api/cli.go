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
}

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
