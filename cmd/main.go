package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/saeidalz13/gurl/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("must provide domain name")
		os.Exit(1)
	}

	domainCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	_ = domainCmd.String("method", "GET", "HTTP method")
	domainCmd.Parse(os.Args[2:])

	// os.Args[1] will be domain string
	destIp := api.MustFetchDomainIp(os.Args[1])
	fmt.Println(destIp.String())

	// tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer tcpConn.Close()

	// for {
	// 	_, err := tcpConn.Write([]byte(`GET / HTTP/1.1\r\nHost: google.com\r\nUser-Agent: Client\r\n\r\n`))
	// 	if err != nil {
	// 		log.Printf("error write: %v\n", err)
	// 		break
	// 	}

	// 	buf := make([]byte, 2048)
	// 	n, err := tcpConn.Read(buf)
	// 	if err != nil {
	// 		if err.Error() == "EOF" {
	// 			os.Exit(0)
	// 		}

	// 		log.Printf("error read: %v\n", err)
	// 		continue
	// 	}

	// 	fmt.Println(string(buf[:n]))
	// }

	// 2. Make a TCP connection with the server
	// This might involve TLS (Almost always it does)
}
