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
	method           string
	contentType      string
	data             string
	additonalHeaders []string

	sb *strings.Builder
}

func NewHTTPRequestGenerator(domain, path, cookies, method, contentType, data string) HTTPRequestGenerator {
	return HTTPRequestGenerator{
		domain:           domain,
		path:             path,
		cookies:          cookies,
		method:           method,
		contentType:      contentType,
		data:             data,
		additonalHeaders: make([]string, 0, 3),
	}
}

func (h *HTTPRequestGenerator) adjustHeaderForData() {
	h.additonalHeaders = append(h.additonalHeaders, fmt.Sprintf("Content-Type: %s", h.contentType))
	h.additonalHeaders = append(h.additonalHeaders, fmt.Sprintf("Content-Length: %d", len(h.data)))
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

func (h HTTPRequestGenerator) generateGETRequest() string {
	h.addGenericPartsHeader(httpconstants.MethodGET)
	h.addCookie()
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")
	return h.sb.String()
}

func (h HTTPRequestGenerator) generatePOSTRequest() string {
	h.addGenericPartsHeader(httpconstants.MethodPOST)
	h.addCookie()
	h.adjustHeaderForData()
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")

	h.sb.WriteString(h.data)

	return h.sb.String()
}

func (h HTTPRequestGenerator) generatePUTPATCHRequest() string {
	h.addGenericPartsHeader(httpconstants.MethodPUT)
	h.addCookie()
	h.adjustHeaderForData()
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")

	h.sb.WriteString(h.data)

	return h.sb.String()
}

func (h HTTPRequestGenerator) generateDELETERequest() string {
	h.addGenericPartsHeader(httpconstants.MethodDELETE)
	h.addCookie()
	h.addAdditionalHeaders()

	// separator between body and header
	h.sb.WriteString("\r\n")
	return h.sb.String()
}

func (h HTTPRequestGenerator) Generate() string {
	switch h.method {
	case httpconstants.MethodGET:
		return h.generateGETRequest()

	case httpconstants.MethodPOST:
		return h.generatePOSTRequest()

	case httpconstants.MethodPUT, httpconstants.MethodPATCH:
		return h.generatePUTPATCHRequest()

	case httpconstants.MethodDELETE:
		return h.generateDELETERequest()

	default:
		// The error catching has been taken care of
		// in the steps before
		return ""
	}
}
