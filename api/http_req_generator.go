package api

import (
	"fmt"
	"strings"

	"github.com/saeidalz13/gurl/internal/httpconstants"
)

type HTTPRequestGenerator struct {
	domain           string
	path             string
	cookies          string
	additonalHeaders []string

	sb *strings.Builder
}

func NewHTTPRequestGenerator(domain, path, cookies string) HTTPRequestGenerator {
	return HTTPRequestGenerator{
		domain:           domain,
		path:             path,
		cookies:          cookies,
		additonalHeaders: make([]string, 0, 3),
	}
}

func (h *HTTPRequestGenerator) adjustHeaderForData(contentType string, contentLength int) {
	h.additonalHeaders = append(h.additonalHeaders, fmt.Sprintf("Content-Type: %s", contentType))
	h.additonalHeaders = append(h.additonalHeaders, fmt.Sprintf("Content-Length: %d", contentLength))
}

func (h *HTTPRequestGenerator) addAdditionalHeaders() {
	for _, header := range h.additonalHeaders {
		h.sb.WriteString(header)
		h.sb.WriteString("\r\n")
	}
}

func (h *HTTPRequestGenerator) addCookie() {
	if h.cookies != "" {
		h.additonalHeaders = append(h.additonalHeaders, fmt.Sprintf("Cookie: %s", h.cookies))
	}
}

func (h *HTTPRequestGenerator) addGenericPartsHeader(method string) {
	sb := &strings.Builder{}
	sb.Grow(50)
	// Method
	sb.WriteString(method)
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

	h.sb = sb
}

func (h HTTPRequestGenerator) GenerateGETRequest() string {
	h.addGenericPartsHeader(httpconstants.MethodGET)
	h.addCookie()
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")
	return h.sb.String()
}

func (h HTTPRequestGenerator) GeneratePOSTRequest(data, contentType string) string {
	h.addGenericPartsHeader(httpconstants.MethodPOST)
	h.addCookie()
	h.adjustHeaderForData(contentType, len(data))
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")

	h.sb.WriteString(data)

	return h.sb.String()
}
