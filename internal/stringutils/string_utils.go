package stringutils

import (
	"errors"
	"strings"
)

func PrepareDomainSegments(domain string) ([]string, error) {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	if domain == "" {
		return nil, errors.New("no domain provided")
	}

	domainSegments := strings.Split(domain, ".")
	if len(domainSegments) < 2 {
		return nil, errors.New("invalid domain - must be a string delimited by dot")
	}

	return domainSegments, nil
}
