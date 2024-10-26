package domainparser

import (
	"errors"
	"strings"
)

type DomainParser struct {
	IsLocalHost   bool
	IsWebSocket   bool
	IsTLS         bool
	Domain        string
	Path          string
	Port          int
	DomainSegment []string
}

func NewDomainParser(domain string) DomainParser {
	return DomainParser{Domain: domain}
}

func (d *DomainParser) splitDomainIntoSegments() error {
	domainSegments := strings.Split(d.Domain, ".")
	if len(domainSegments) < 2 {
		return errors.New("invalid domain - must be a string delimited by dot")
	}

	d.DomainSegment = domainSegments
	return nil
}

func (d *DomainParser) trimProtocolFromHTTPDomain() error {
	d.Domain = strings.TrimPrefix(d.Domain, "http://")
	d.Domain = strings.TrimPrefix(d.Domain, "https://")

	if d.Domain == "" {
		return errors.New("no domain provided")
	}

	return nil
}

func (d *DomainParser) separateDomainAndPath() {
	segments := strings.SplitN(d.Domain, "/", 2)
	d.Domain = segments[0]

	if len(segments) == 1 || segments[1] == "" {
		d.Path = "/"
	} else {
		d.Path = "/" + segments[1]
	}
}

func (d *DomainParser) trimProtocolFromWebSocketDomain() error {
	after, found := strings.CutPrefix(d.Domain, "ws://")
	if found {
		d.Domain = after
		d.IsTLS = false
		return nil
	}

	after, found = strings.CutPrefix(d.Domain, "wss://")
	if found {
		d.Domain = after
		d.IsTLS = true
		return nil
	}

	return errors.New("if websocket, input domain must contain protocol: ws:// or wss://")
}

func (d *DomainParser) determineIfWebSocket() {
	if strings.HasPrefix(d.Domain, "ws://") || strings.HasPrefix(d.Domain, "wss://") {
		d.IsWebSocket = true
	}
}

func (d *DomainParser) determineIfLocalhost() {
	if strings.Contains(d.Domain, "localhost") || strings.Contains(d.Domain, "127.0.0.1") {
		d.IsLocalHost = true
	}
}

func (d *DomainParser) Parse() error {
	d.Domain = strings.TrimSpace(d.Domain)
	d.determineIfWebSocket()
	d.determineIfLocalhost()

	// WebSocket protocol
	if d.IsWebSocket {
		if err := d.trimProtocolFromWebSocketDomain(); err != nil {
			return err
		}
	} else {
		if err := d.trimProtocolFromHTTPDomain(); err != nil {
			return err
		}
	}

	d.separateDomainAndPath()
	d.splitDomainIntoSegments()
	return nil
}
