package httpconstants

const (
	MethodGET    string = "GET"
	MethodPOST   string = "POST"
	MethodPATCH  string = "PATCH"
	MethodDELETE string = "DELETE"
	MethodPUT    string = "PUT"
)

var ValidHttpMethods = map[string]bool{
	MethodGET:    true,
	MethodPOST:   true,
	MethodPATCH:  true,
	MethodDELETE: true,
	MethodPUT:    true,
}

const (
	PortHTTPS = 443
	PortHTTP  = 80
)
