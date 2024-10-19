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

	for _, vm := range httpconstants.ValidHttpMethods {
		if vm == method {
			gp.method = method
			return
		}
	}

	fmt.Printf("invalid method: %s\n", method)
	os.Exit(1)
}

func newGurlParams(rawMethod string, ctJson bool) gurlParams {
	gp := gurlParams{headers: make([]string, 0)}

	if ctJson {
		gp.headers = append(gp.headers, "Content-Type: application/json")
	}

	gp.mustParseDomain()
	gp.mustParseMethod(rawMethod)

	return gp
}
