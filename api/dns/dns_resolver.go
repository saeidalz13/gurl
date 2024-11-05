package dns

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/saeidalz13/gurl/internal/errutils"
)

type DNSResolver struct {
	domainSegments []string
}

func NewDNSResolver(domainSegments []string) DNSResolver {
	return DNSResolver{
		domainSegments: domainSegments,
	}
}

// Fetch the domain IPv4 from 8.8.8.8 (Google server).
// Average time is 25 ms.
func (dr DNSResolver) MustResolveIP() net.IP {
	query := NewDNSQueryManager().CreateQuery()

	// dnsQuery, err := dr.createDNSQuery()
	// errutils.CheckErr(err)

	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.IPv4(8, 8, 8, 8)})
	errutils.CheckErr(err)
	defer udpConn.Close()

	_, err = udpConn.Write(query)
	errutils.CheckErr(err)

	// DNS responses are small, 128 bytes is enough
	// Response share the same structure of request with an additional Answers section
	response := make([]byte, 128)
	_, _, err = udpConn.ReadFrom(response)
	errutils.CheckErr(err)
	/*
	   DNS Response:
	   * Header
	   id 2 bytes
	   flags 2 bytes
	   question count 2 bytes
	   answer count 2 bytes
	   authority RRs 2 bytes
	   additional RRs 2 bytes

	   * Query
	   Name (the domain name) variable
	   Type (DNS record type (e.g., A, CNAME, and MX)) 2 bytes
	   Class (allows domain names to be used for arbitrary objects) 2 bytes

	   * Answer
	   Name (variable)
	   Type (2 bytes)
	   Class (2 bytes)
	   TTL (4 bytes)
	   Data Length (2 bytes)
	   Data (addr, Cname) (variable)

	   The last 4 bytes are the IPv4 digits (each a single byte)
	   In `createDNSQuery` we asked for type A host, so we get IPv4 here
	*/

	// This shows if Answer section exists
	anCount := binary.BigEndian.Uint16(response[6:8])
	if anCount == 0 {
		log.Fatalln("no answer received from 8.8.8.8 server")
	}

	pos := 12
	endQNameValue := 0
	for int(response[pos]) != endQNameValue {
		// Move by length of each label + 1 for the length byte
		// Each label prefixed by a length byte, followed by the label itself
		// e.g: A domain name consists of labels (e.g., www, example, com in www.example.com).
		lenghtOfBytes := int(response[pos]) // value of length byte
		pos += lenghtOfBytes + 1
	}
	pos += 1 // End of domain 0x00
	pos += 4 // QType 2 + QClass 2

	// Now starting Answer section
	// Name 2, Type 2, Class 2, TTL 4, Data length 2
	pos += 12
	dataLen := binary.BigEndian.Uint16(response[pos-2 : pos])

	if dataLen != 4 {
		log.Fatalln("invalid ip address received from dns - must be ipv4")
	}

	ip := net.IPv4(response[pos], response[pos+1], response[pos+2], response[pos+3])
	return ip
}
