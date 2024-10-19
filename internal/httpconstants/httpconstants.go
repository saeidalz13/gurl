package httpconstants

const (
	MethodGET    string = "GET"
	MethodPOST   string = "POST"
	MethodPATCH  string = "PATCH"
	MethodDELETE string = "DELETE"
	MethodPUT    string = "PUT"
)

var ValidHttpMethods = []string{
	MethodGET,
	MethodPOST,
	MethodPATCH,
	MethodDELETE,
	MethodPUT,
}

const (
	PortHTTPS = "443"
	PortHTTP  = "80"
)
