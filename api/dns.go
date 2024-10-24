package api

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/stringutils"
)

type DNSResolver struct {
	domain string
}

func newDNSResolver(domain string) DNSResolver {
	return DNSResolver{domain}
}

/*
Header sections

  - ID
  - Flags
  - Questions
  - Answer RRs (RR stands for Resource Record)
  - Authority RRs
  - Additional RRs

Queries section:
  - Whatever you want to put in the query
*/
func (dr DNSResolver) createDNSQuery(domainSegments []string) ([]byte, error) {
	id := uint16(rand.Intn(65535)) // To be within 8 bytes

	minBytesNeeded := 16 // header 12 + QTYPE 2 + QCLASS 2
	query := make([]byte, 0, minBytesNeeded)

	// * Header Section (12 bytes)
	// Transaction ID
	// We want the generated into 2 bytes as requested by DNS standard
	// Extracts the **most significant byte (MSB)** (shifts 16-bit, 8 bits to the right)
	// This means number / (2^8) -> This drops the remainder
	query = append(query, byte(id>>8))
	// bitwise AND that extracts the **least significant byte (LSB)**
	query = append(query, byte(id&0xff))

	// Flags
	// This means this is a standard query
	// Every bit of this is for a certain flag.
	// This just means standard query
	query = append(query, 0x01, 0x00)

	// QDCOUNT -> number of entries in question section
	query = append(query, 0x00, 0x01)

	// ANCOUNT -> number of RR in Answer section which should be
	// zero for the request
	query = append(query, 0x00, 0x00)

	// NSCOUNT -> number of RR in Authority section which
	// we do not have. Only meaningful in responses, so zero
	query = append(query, 0x00, 0x00)

	// ARCOUNT -> number of RR in the Additional section.
	// No additional records follow
	query = append(query, 0x00, 0x00)

	// * Question Section (Variable len)
	// QNAME (Domain section)
	for _, part := range domainSegments {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			return nil, fmt.Errorf("invalid input domain")
		}
		// Each label prefixed by a length byte, followed by the label itself
		query = append(query, byte(len(part)))
		query = append(query, []byte(part)...)
	}
	query = append(query, 0x00) // To show that this is end of the domain

	// QTYPE -> Type A (host address) - 2 bytes (A, AAAA, MX, etc.)
	query = append(query, 0x00, 0x01)

	// QCLASS -> Class IN (Internet) - 2 bytes
	query = append(query, 0x00, 0x01)

	return query, nil
}

// Fetch the domain IPv4 from 8.8.8.8 (Google server).
// Average time is 25 ms.
func (dr DNSResolver) mustResolveIP() net.IP {
	domainSegments, err := stringutils.SplitDomainIntoSegments(dr.domain)
	errutils.CheckErr(err)

	dnsQuery, err := dr.createDNSQuery(domainSegments)
	errutils.CheckErr(err)

	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.IPv4(8, 8, 8, 8)})
	errutils.CheckErr(err)
	defer udpConn.Close()

	_, err = udpConn.Write(dnsQuery)
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
