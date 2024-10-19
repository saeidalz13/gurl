package api

import (
	"fmt"

	"github.com/saeidalz13/gurl/internal/httpconstants"
)

type gurlCli struct {
	domain string
	method string
}

func ExecGurl() {
	gc := initGurlCli()
	ip := mustFetchDomainIp(gc.domain)

	tlsConn := initTlsConn(ip.String(), gc.domain)
	defer tlsConn.Close()

	switch gc.method {
	case httpconstants.MethodGET:
		execGetHttpReq(tlsConn, gc.domain)

	case httpconstants.MethodPOST:
		fmt.Println("Posting...")
	}
}
