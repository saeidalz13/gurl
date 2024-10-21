package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
	"github.com/saeidalz13/gurl/internal/stringutils"
)

type gurlParams struct {
	domain  string
	method  string
	path    string
	headers []string

	// Aesthetics
	pretty bool
}

func (gp *gurlParams) mustParseDomain() {
	domain := os.Args[1]
	domain, err := stringutils.TrimDomainPrefix(domain)
	errutils.CheckErr(err)
	gp.domain, gp.path = stringutils.ExtractPathFromDomain(domain)
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

func newGurlParams(cp cliParams) gurlParams {
	gp := gurlParams{headers: make([]string, 0)}

	if cp.isHeaderJson {
		gp.headers = append(gp.headers, "Content-Type: application/json")
	}

	gp.pretty = cp.isPretty

	gp.mustParseDomain()
	gp.mustParseMethod(cp.method)

	return gp
}
