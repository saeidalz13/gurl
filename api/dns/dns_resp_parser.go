package dns

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	ipTypeV4 uint8 = iota
	ipTypeV6
)

const (
	startingQueryIdx       int = 12
	endOfQuestionIndicator int = 0
	byteForLength          int = 1
)

/*
DNS Response:
* Header
id 2 bytes
flags 2 bytes (second byte, in bits 4-7 contains probable error)
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
type DNSResponseParser struct {
	iptype   uint8
	pos      int
	response []byte
}

func NewDNSResponseParser(response []byte, iptype uint8) DNSResponseParser {
	return DNSResponseParser{
		response: response,
		pos:      startingQueryIdx,
		iptype:   iptype,
	}
}

func (drp *DNSResponseParser) isAnswerAvailable() bool {
	anCount := binary.BigEndian.Uint16(drp.response[6:8])
	return anCount != 0
	// if anCount == 0 {
	// 	log.Fatalln("no answer received from 8.8.8.8 server")
	// }
}

func (drp *DNSResponseParser) determineIdxAfterQuery() {
	for int(drp.response[drp.pos]) != endOfQuestionIndicator {
		// Move by length of each label + 1 for the length byte
		// Each label prefixed by a length byte, followed by the label itself
		// e.g: A domain name consists of labels (e.g., www, example, com in www.example.com).
		lenghtOfBytes := int(drp.response[drp.pos]) // value of length byte
		drp.pos += lenghtOfBytes + byteForLength
	}
	drp.pos += 1 // End of domain 0x00
	drp.pos += 4 // QType 2 + QClass 2
}

func (drp *DNSResponseParser) determineDataLength() uint16 {
	// Assuming name is a pointer because of name
	// compression.
	// TODO: consider other scenarios
	drp.pos += 2

	// Type 2, Class 2, TTL 4, Data length 2
	drp.pos += 10
	return binary.BigEndian.Uint16(drp.response[drp.pos-2 : drp.pos])
}

func (drp DNSResponseParser) Parse() (net.IP, error) {
	if !drp.isAnswerAvailable() {
		return nil, fmt.Errorf("no answer received from 8.8.8.8 server")
	}

	drp.determineIdxAfterQuery()
	dataLength := drp.determineDataLength()

	switch drp.iptype {
	case ipTypeV4:
		if dataLength != 4 {
			return nil, fmt.Errorf("dns server did not provide ipv4")
		}
		ip := net.IPv4(drp.response[drp.pos], drp.response[drp.pos+1], drp.response[drp.pos+2], drp.response[drp.pos+3])
		return ip, nil
	case ipTypeV6:
		if dataLength != 16 {
			return nil, fmt.Errorf("dns server did not provide ipv6")
		}
		return net.IPv6zero, nil
	default:
		return nil, fmt.Errorf("unsupported IP type")
	}
}
