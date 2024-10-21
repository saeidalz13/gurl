package api

import (
	"crypto/tls"
	"fmt"
	"strings"
)

/*
HTTP response includes 3 main segments
 1. Status line
 2. Headers (variable number of lines)
 3. Body (Separated from the rest with \r\n)
*/
type HTTPResponse struct {
	version    string
	statusCode string
	statusMsg  string
	headers    []string
	body       string
}

func newHTTPResponse(response string) HTTPResponse {
	responseSegments := strings.Split(response, "\r\n")
	httpResp := HTTPResponse{
		headers: make([]string, 0, 3),
	}

	var bodyIdx int
	for i, segment := range responseSegments {
		// First line is always the status line
		if i == 0 {
			statusLineSegments := strings.Split(segment, " ")
			httpResp.version = statusLineSegments[0]
			httpResp.statusCode = statusLineSegments[1]
			httpResp.statusMsg = statusLineSegments[2]
			continue
		}

		// When separating by \r\n, the empty line would
		// would be the separator between header and body
		if segment == "" {
			bodyIdx = i + 1
			break
		}

		httpResp.headers = append(httpResp.headers, segment)
	}

	sb := strings.Builder{}
	sb.Grow(10) // 10 is min length of body

	// We found the idx that starts the body
	// aggregate again to have body string
	for _, seg := range responseSegments[bodyIdx:] {
		sb.WriteString(seg)
	}

	httpResp.body = sb.String()

	return httpResp
}

func execHTTPReq(tlsConn *tls.Conn, httpRequest string, gp gurlParams) {
	if err := writeToTLSConn(tlsConn, []byte(httpRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		return
	}

	readBytes := readFromTLSConn(tlsConn)

	if gp.pretty {
		httpResp := newHTTPResponse(string(readBytes))
		fmt.Printf("%+v\n", httpResp)
	} else {
		fmt.Println(string(readBytes))
	}
}
