package api

import (
	"strings"
)

type HTTPRequest struct {
	domain           string
	path             string
	method           string
	additonalHeaders []string
}

func newHTTPRequestCreator(domain, path, method string) HTTPRequest {
	return HTTPRequest{
		domain:           domain,
		path:             path,
		method:           method,
		additonalHeaders: make([]string, 0),
	}
}

func (h *HTTPRequest) AddContentTypeJson() {
	h.additonalHeaders = append(h.additonalHeaders, "application/json")
}

func (h HTTPRequest) Create() string {
	sb := strings.Builder{}
	sb.Grow(50)

	// Method
	sb.WriteString(h.method)
	sb.WriteString(" ")

	// Path
	sb.WriteString(h.path)
	sb.WriteString(" ")

	// Protocol and version
	sb.WriteString("HTTP/1.1\r\n")

	// Host
	sb.WriteString("Host: ")
	sb.WriteString(h.domain)
	sb.WriteString("\r\n")

	// User-Agent
	sb.WriteString("User-Agent: gurl/1.0.0\r\n")

	// Accept headers
	sb.WriteString("Accept: */*\r\n")

	// User Headers
	for _, header := range h.additonalHeaders {
		sb.WriteString(header)
		sb.WriteString("\r\n")
	}

	// Ending of request based on HTTP
	sb.WriteString("\r\n")

	return sb.String()
}
