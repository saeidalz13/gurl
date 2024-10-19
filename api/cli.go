package api

import (
	"flag"
	"fmt"
	"os"
)

func initGurlCli() gurlParams {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	ctJsonPtr := domainCmd.Bool("json", false, "Set Header -> Content-Type: application/json")

	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		domainCmd.Usage()
		os.Exit(1)
	}

	domainCmd.Parse(os.Args[2:])

	return newGurlParams(*methodPtr, *ctJsonPtr)
}
