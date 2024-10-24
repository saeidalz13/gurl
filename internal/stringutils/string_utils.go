package stringutils

import (
	"errors"
	"strings"

	"github.com/saeidalz13/gurl/internal/wsutils"
)

func SplitDomainIntoSegments(domain string) ([]string, error) {
	domainSegments := strings.Split(domain, ".")
	if len(domainSegments) < 2 {
		return nil, errors.New("invalid domain - must be a string delimited by dot")
	}

	return domainSegments, nil
}

func TrimDomainPrefix(domain string) (string, error) {
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

func ExtractWsProtocol(domain string) (uint8, string, error) {
	after, found := strings.CutPrefix(domain, "ws://")
	if found {
		return wsutils.ProtocolWS, after, nil
	}

	after, found = strings.CutPrefix(domain, "wss://")
	if found {
		return wsutils.ProtocolWSS, after, nil
	}

	return 255, "", errors.New("if websocket, input domain must contain protocol: ws:// or wss://")
}

func TrimDomainSpace(domain string) string {
	return strings.TrimSpace(domain)
}

func IsDomainForWebsocket(domain string) bool {
	return strings.HasPrefix(domain, "ws://") || strings.HasPrefix(domain, "wss://")
}
