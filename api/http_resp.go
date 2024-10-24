package api

import (
	"fmt"
	"strings"

	"github.com/saeidalz13/gurl/internal/encodingutils"
	"github.com/saeidalz13/gurl/internal/terminalutils"
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

	responseSegments []string
}

func newHTTPResponse(respBytes []byte) HTTPResponse {
	responseSegments := strings.Split(string(respBytes), "\r\n")
	httpResp := HTTPResponse{
		headers:          make([]string, 0, 3),
		responseSegments: responseSegments,
	}
	return httpResp
}

func (hr HTTPResponse) determineStatusCodeBashColor() string {
	switch hr.statusCode[0] {
	case encodingutils.ASCII2:
		return terminalutils.BoldGreen

	case encodingutils.ASCII3:
		return terminalutils.BoldCyan

	case encodingutils.ASCII4:
		return terminalutils.BoldRed

	case encodingutils.ASCII5:
		return terminalutils.BoldPurple
	}

	return "\033[0;31m"
}

func (hr HTTPResponse) parse() HTTPResponse {
	var bodyIdx int
	for i, segment := range hr.responseSegments {
		// First line is always the status line
		if i == 0 {
			statusLineSegments := strings.Split(segment, " ")
			hr.version = statusLineSegments[0]
			hr.statusCode = statusLineSegments[1]
			hr.statusMsg = statusLineSegments[2]
			continue
		}

		// When separating by \r\n, the empty line would
		// would be the separator between header and body
		if segment == "" {
			bodyIdx = i + 1
			break
		}

		hr.headers = append(hr.headers, segment)
	}

	sb := strings.Builder{}
	sb.Grow(10) // 10 is min length of body

	// We found the idx that starts the body
	// aggregate again to have body string
	for _, seg := range hr.responseSegments[bodyIdx:] {
		sb.WriteString(seg)
	}

	hr.body = sb.String()

	return hr
}

func (hr HTTPResponse) printPretty(verbose bool) {
	if verbose {
		fmt.Println("\n\033[1;37mStatus\033[0m")
		fmt.Println("---------------------")
		fmt.Printf("\033[0;33mHTTP Version\033[0m   | %s \n", hr.version)
		fmt.Printf("\033[0;33mStatus Code    | %s%s\033[0m\n", hr.determineStatusCodeBashColor(), hr.statusCode)
		fmt.Printf("\033[0;33mStatus Message\033[0m | %s \n", hr.statusMsg)

		fmt.Println("\n\033[1;37mHeaders\033[0m")
		fmt.Println("---------------------")
		for _, header := range hr.headers {
			headerSegments := strings.Split(header, ":")
			fmt.Printf("\033[0;36m%s\033[0m: %s\n", headerSegments[0], headerSegments[1])
			// fmt.Println("")
		}
	}

	fmt.Println("\n\033[1;37mBody\033[0m")
	fmt.Println("---------------------")
	fmt.Println(hr.body)
}
