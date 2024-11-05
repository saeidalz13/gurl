package api

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
)

const (
	// Minimum capacity needed for query slice
	// sent to DNS
	// header 12 + QTYPE 2 + QCLASS 2
	minRequiredCap = 16
)

type DNSResolver struct {
	domainSegments []string
}

func NewDNSResolver(domainSegments []string) DNSResolver {
	return DNSResolver{
		domainSegments: domainSegments,
	}
}

/*
Prepares the query needs to be sent to DNS server.

It is comprised of header and query section.
12 bytes for header and variable length for
the query section.
*/
func (dr DNSResolver) createDNSQuery() ([]byte, error) {
	id := uint16(rand.Intn(65535)) // To be within 8 bytes
	query := make([]byte, 0, minRequiredCap)

	/*
		Header Section (12 bytes)
	*/
	// ** Transaction ID (2 bytes)
	// We want the generated into 2 bytes as requested by DNS standard
	// Extracts the **most significant byte (MSB)** (shifts 16-bit, 8 bits to the right)
	// This means number / (2^8) -> This drops the remainder
	query = append(query, byte(id>>8))
	// bitwise AND that extracts the **least significant byte (LSB)**
	query = append(query, byte(id&0xff))

	// ** Flags (2 bytes)
	// This means this is a standard query
	// Every bit of this is for a certain flag.
	// 0x01 + 0x00 means standard query
	// Byte 1: QR (bit 0), OPCODE (bits 1 to 4), AA (bit 5), TC (bit 6), RD (bit 7)
	// Byte 2: RA (bit 0), Z (bit 1 to 3), RCODE (bit 4 to 7)
	query = append(query, 0b00000001, 0b00000000)

	// ** QDCOUNT -> number of entries in question section (2 bytes -> 16 bit integer).
	// query = append(query, 0x00, 0x01)
	query = append(query, 0b00000000, 0b00000001)

	// ** ANCOUNT -> number of RR in Answer section which should be (2 bytes -> 16 bit integer).
	// zero for the request
	query = append(query, 0b00000000, 0b00000000)

	// ** NSCOUNT -> number of RR in Authority section which (2 bytes -> 16 bit integer).
	// we do not have. Only meaningful in responses, so zero
	query = append(query, 0b00000000, 0b00000000)

	// ** ARCOUNT -> number of RR in the Additional section (2 bytes -> 16 bit integer).
	// No additional records follow
	query = append(query, 0b00000000, 0b00000000)

	/*
		Query Section (Variable length)
	*/
	// ** QNAME (Domain section) (Variable len)
	for _, part := range dr.domainSegments {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			return nil, fmt.Errorf("invalid input domain")
		}
		// Each label prefixed by a length byte, followed by the label itself
		query = append(query, byte(len(part)))
		query = append(query, []byte(part)...)
	}
	query = append(query, 0b00000000) // To show that this is end of the domain

	// ** QTYPE -> Type A (host address) - 2 bytes (A, AAAA, MX, etc.)
	query = append(query, 0b00000000, 0b00000001)
	// Query for IPv6 (AAAA record)
	// query := append(query, 0b00000000, 0b00011100)

	// ** QCLASS -> Class IN (Internet) - 2 bytes
	query = append(query, 0b00000000, 0b00000001)

	return query, nil
}

// Fetch the domain IPv4 from 8.8.8.8 (Google server).
// Average time is 25 ms.
func (dr DNSResolver) MustResolveIP() net.IP {
	dnsQuery, err := dr.createDNSQuery()
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
