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

func newHTTPResponse(tcpResponse string) HTTPResponse {
	responseSegments := strings.Split(tcpResponse, "\r\n")
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

func printPretty(httpResp HTTPResponse) {
	fmt.Println("\n\033[1;37mStatus\033[0m")
	fmt.Println("---------------------")
	fmt.Printf("\033[0;33mHTTP Version\033[0m   | %s \n", httpResp.version)
	fmt.Printf("\033[0;33mStatus Code\033[0m    | %s \n", httpResp.statusCode)
	fmt.Printf("\033[0;33mStatus Message\033[0m | %s \n", httpResp.statusMsg)

	fmt.Println("\n\033[1;37mHeaders\033[0m")
	fmt.Println("---------------------")
	for _, header := range httpResp.headers {
		headerSegments := strings.Split(header, ":")
		fmt.Printf("\033[0;36m%s\033[0m: %s\n", headerSegments[0], headerSegments[1])
		// fmt.Println("")
	}

	fmt.Println("\n\033[1;37mBody\033[0m")
	fmt.Println("---------------------")
	fmt.Println(httpResp.body)
}

func execHTTPReq(tlsConn *tls.Conn, httpRequest string) {
	if err := writeToTLSConn(tlsConn, []byte(httpRequest)); err != nil {
		fmt.Printf("write tcp read: %v\n", err)
		return
	}

	tcpRespBytes := readFromTLSConn(tlsConn)
	tcpResponse := string(tcpRespBytes)

	httpResp := newHTTPResponse(tcpResponse)
	printPretty(httpResp)
}
