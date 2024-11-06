package http

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
type HTTPResponseParser struct {
	version    string
	statusCode string
	statusMsg  string
	headers    []string
	body       string

	responseSegments []string
}

func NewHTTPResponseParser(respBytes []byte) HTTPResponseParser {
	responseSegments := strings.Split(string(respBytes), "\r\n")
	httpResp := HTTPResponseParser{
		headers:          make([]string, 0, 3),
		responseSegments: responseSegments,
	}
	return httpResp
}

func (hr HTTPResponseParser) determineStatusCodeBashColor() string {
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

// In case of Transfer-Encoding: chunked, there will be
// two integers at the beginning and end of JSON resp.
// these 2 should be removed for aesthetics
func (hr HTTPResponseParser) trimJsonResp(body string) string {
	for _, header := range hr.headers {
		headerSegments := strings.Split(header, ":")

		if strings.TrimSpace(headerSegments[1]) == "chunked" {
			// the first characters until reaching "{"
			// are certainly ASCII
			var jsonStartIdx int
		strStartLoop:
			for i, char := range body {
				if char == '{' {
					jsonStartIdx = i
					break strStartLoop
				}
			}

			// '0' will exist at the end of the
			// resp so we should remove that too.
			// This indicates the end of the resp
			// as a part of HTTP standards.
			extra0CharIdx := len(body) - 1
			return body[jsonStartIdx:extra0CharIdx]
		}
	}

	return body
}

func (hr HTTPResponseParser) Parse() HTTPResponseParser {
	var bodyIdx int
	statusLineNum := 0

	for i, segment := range hr.responseSegments {
		// First line is always the status line
		if i == statusLineNum {
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

	hr.body = hr.trimJsonResp(sb.String())

	return hr
}

func (hr HTTPResponseParser) Print(verbose bool) {
	if verbose {
		fmt.Printf("\n%sStatus%s\n", terminalutils.BoldYellow, terminalutils.FormatReset)
		fmt.Println("---------------------")
		fmt.Printf("%sHTTP Version%s   | %s \n", terminalutils.RegularYellow, terminalutils.FormatReset, hr.version)
		fmt.Printf("%sStatus Code    | %s%s\n", terminalutils.RegularYellow, hr.determineStatusCodeBashColor(), hr.statusCode)
		fmt.Printf("%sStatus Message%s | %s \n", terminalutils.RegularYellow, terminalutils.FormatReset, hr.statusMsg)

		fmt.Printf("\n%sHeaders%s\n", terminalutils.BoldCyan, terminalutils.FormatReset)
		fmt.Println("---------------------")
		for _, header := range hr.headers {
			headerSegments := strings.Split(header, ":")
			fmt.Printf("%s%s%s: %s\n", terminalutils.RegularCyan, headerSegments[0], terminalutils.FormatReset, headerSegments[1])
		}
	}

	fmt.Printf("\n%sBody%s\n", terminalutils.BoldGreen, terminalutils.FormatReset)
	fmt.Println("---------------------")
	fmt.Println(hr.body)
}
