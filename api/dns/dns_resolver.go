package dns

import (
	"log"
	"net"

	"github.com/saeidalz13/gurl/internal/errutils"
)

// Fetch the domain IPv4 from 8.8.8.8 (Google server).
// Average time is 25 ms.
func MustResolveIP(domainSegments []string) net.IP {
	ipType := ipTypeV6

dnsLoop:
	for {
		query := NewDNSQueryManager(domainSegments, ipType).CreateQuery()

		udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.IPv4(8, 8, 8, 8)})
		errutils.CheckErr(err)
		defer udpConn.Close()

		_, err = udpConn.Write(query)
		errutils.CheckErr(err)

		// DNS responses are small, 128 bytes is enough
		// Response share the same structure of request with an additional Answers section
		response := make([]byte, 256)
		_, _, err = udpConn.ReadFrom(response)
		errutils.CheckErr(err)

		ip, err := NewDNSResponseParser(response, ipType).Parse()
		if err != nil {
			switch err.Error() {
			case "no ipv4":
				ipType = ipTypeV6
				continue dnsLoop
			case "no ipv6":
				log.Fatalln("could not fetch ip from DNS")
			default:
				log.Fatalln(err)
			}
		}

		return ip
	}
}
