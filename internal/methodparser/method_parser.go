package methodparser

import (
	"fmt"
	"strings"

	"github.com/saeidalz13/gurl/internal/httpconstants"
)

func ParseMethod(rawMethod string) (string, error) {
	method := strings.TrimSpace(rawMethod)
	method = strings.ToUpper(method)

	_, ok := httpconstants.ValidHttpMethods[method]
	if !ok {
		return "", fmt.Errorf("invalid method: %s", method)
	}

	return method, nil
}
