package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/saeidalz13/gurl/internal/httpconstants"
)

type cliParams struct {
	Verbose  bool
	DataType uint8
	Data     string
	Domain   string
	Method   string
	Cookies  string
}

func mustDetermineDataInfo(jsonPtr, textPtr *string) (uint8, string) {
	var dataType uint8
	var data string

	if *jsonPtr != "" && *textPtr != "" {
		fmt.Println("only one body should be selected")
		os.Exit(1)
	}

	if *jsonPtr != "" {
		return httpconstants.DataTypeJson, *jsonPtr
	}

	if *textPtr != "" {
		return httpconstants.DataTypeText, *textPtr
	}

	return dataType, data
}

func InitCli() cliParams {
	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	methodPtr := domainCmd.String("method", "GET", "HTTP method")
	jsonPtr := domainCmd.String("json", "", "Add json data to body")
	textPtr := domainCmd.String("text", "", "Add plain text to body")
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

	dataType, data := mustDetermineDataInfo(jsonPtr, textPtr)

	return cliParams{
		Domain:   os.Args[1],
		Method:   *methodPtr,
		Verbose:  *verbose,
		Data:     data,
		DataType: dataType,
		Cookies:  *cookies,
	}
}
