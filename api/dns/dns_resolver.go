package dns

import (
	"fmt"
	"net"
	"os"

	"github.com/saeidalz13/gurl/internal/errutils"
)

// Fetch the domain IPv4 from 8.8.8.8 (Google server).
// Average time is 25 ms.
func MustResolveIP(domainSegments []string) (net.IP, uint8) {
	ipType := IpTypeV4
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.IPv4(8, 8, 8, 8)})
	errutils.CheckErr(err)
	defer udpConn.Close()

	dqm := NewDNSQueryManager(domainSegments, ipType)
	dqm.prepareQuery()

dnsLoop:
	for {
		_, err = udpConn.Write(dqm.Query())
		errutils.CheckErr(err)

		// DNS responses are small, 256 bytes is enough. Especially that
		// I only have one question per request.
		response := make([]byte, 256)
		_, _, err = udpConn.ReadFrom(response)
		errutils.CheckErr(err)

		ip, err := NewDNSResponseParser(response, ipType).Parse()
		if err != nil {
			switch err.Error() {
			case "no ipv4":
				ipType = IpTypeV6
				dqm.toggleQuestionType(ipType)

				fmt.Println("ipv4 could not fetched. attempting for ipv6...")
				continue dnsLoop

			case "no ipv6":
				fmt.Println("could not fetch ip from DNS")
				os.Exit(1)

			default:
				fmt.Println(err)
				os.Exit(1)
			}
		}

		return ip, ipType
	}
}
