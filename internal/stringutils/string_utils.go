package stringutils

import (
	"errors"
	"strings"
)

func SplitDomainIntoSegments(domain string) ([]string, error) {
	domainSegments := strings.Split(domain, ".")
	if len(domainSegments) < 2 {
		return nil, errors.New("invalid domain - must be a string delimited by dot")
	}

	return domainSegments, nil
}

func TrimDomainPrefix(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")

	if domain == "" {
		return "", errors.New("no domain provided")
	}

	return domain, nil
}

// Returns domain and path from the raw domain 
// provided by user.
func ExtractPathFromDomain(domain string) (string, string) {
	segments := strings.SplitN(domain, "/", 2)

	if len(segments) == 1 || segments[1] == "" {
		return segments[0], "/"
	}

	return segments[0], "/" + segments[1]
}
