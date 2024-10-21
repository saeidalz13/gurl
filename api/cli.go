package api

import (
	"flag"
	"fmt"
	"os"
)

type cliParams struct {
	method       string
	isHeaderJson bool
	isPretty     bool
}

func initGurlCli() gurlParams {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	ctJsonPtr := domainCmd.Bool("json", false, "Set Header -> Content-Type: application/json")
	isPrettyPtr := domainCmd.Bool("pretty", false, "Print the response in an organized manner")

	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}

	domainCmd.Parse(os.Args[2:])

	return newGurlParams(cliParams{
		method:       *methodPtr,
		isHeaderJson: *ctJsonPtr,
		isPretty:     *isPrettyPtr,
	})
}
