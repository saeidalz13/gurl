package domainparser

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ProtocolHTTP uint8 = iota + 1
	ProtocolHTTPS
	ProtocolWS
)

type DomainParser struct {
	IsLocalHost   bool
	Protocol      uint8
	Domain        string
	Path          string
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
		return nil
	}

	after, found = strings.CutPrefix(d.Domain, "wss://")
	if found {
		d.Domain = after
		return nil
	}

	return errors.New("if websocket, input domain must contain protocol: ws:// or wss://")
}

func (d *DomainParser) determineProtocol() {
	if strings.HasPrefix(d.Domain, "ws://") || strings.HasPrefix(d.Domain, "wss://") {
		d.Protocol = ProtocolWS
	} else if strings.HasPrefix(d.Domain, "http://") {
		d.Protocol = ProtocolHTTP
	} else {
		d.Protocol = ProtocolHTTPS
	}
}

func (d *DomainParser) determineIfLocalhost() {
	if strings.Contains(d.Domain, "localhost") || strings.Contains(d.Domain, "127.0.0.1") {
		d.IsLocalHost = true
	}
}

func (d *DomainParser) Parse() error {
	d.Domain = strings.TrimSpace(d.Domain)
	d.determineProtocol()
	d.determineIfLocalhost()

	switch d.Protocol {
	case ProtocolWS:
		if err := d.trimProtocolFromWebSocketDomain(); err != nil {
			return err
		}

	case ProtocolHTTP, ProtocolHTTPS:
		if err := d.trimProtocolFromHTTPDomain(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("wrong protocol. must never reach here")
	}

	d.separateDomainAndPath()
	d.splitDomainIntoSegments()
	return nil
}
