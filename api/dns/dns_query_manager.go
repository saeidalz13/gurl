package dns

import (
	"log"
	"math/rand"
	"strings"
)

const (
	// Minimum capacity needed for query slice
	// sent to DNS
	// header 12 + QTYPE 2 + QCLASS 2
	minQueryCap = 16
)

type DNSQueryManager struct {
	domainSegments []string
	query          []byte
}

func NewDNSQueryManager(domainSegments []string) *DNSQueryManager {
	return &DNSQueryManager{
		query:          make([]byte, 0, minQueryCap),
		domainSegments: domainSegments,
	}
}

// Transaction ID length is 2 bytes.
func (d *DNSQueryManager) setTransactionId() {
	id := uint16(rand.Intn(65535)) // To be within 8 bytes

	// Extracts the **most significant byte (MSB)** (shifts 16-bit, 8 bits to the right)
	// This means number / (2^8) -> This drops the remainder
	d.query = append(d.query, byte(id>>8))

	// bitwise AND that extracts the **least significant byte (LSB)**
	d.query = append(d.query, byte(id&0xff))
}

// Adding a standard flag which means this is
// just a standard query for fetching IP.
func (d *DNSQueryManager) setStandardFlags() {
	// Byte 1: QR (bit 0), OPCODE (bits 1 to 4), AA (bit 5), TC (bit 6), RD (bit 7)
	// Byte 2: RA (bit 0), Z (bit 1 to 3), RCODE (bit 4 to 7)
	d.query = append(d.query, 0b00000000, 0b00000000)
}

// By default we set this to 1 since we always
// ask for one IP.
func (d *DNSQueryManager) setNumOfQuestions() {
	d.query = append(d.query, 0b00000000, 0b00000001)
}

// Has to be included. This has to be zero to show
// this message is a query
func (d *DNSQueryManager) setNumOfAnswers() {
	d.query = append(d.query, 0b00000000, 0b00000000)
}

// Number of auth resource records. We set the 16-bit
// to zero.
func (d *DNSQueryManager) setNumOfAuthorityRRs() {
	d.query = append(d.query, 0b00000000, 0b00000000)
}

// number of additional resource records. We set the
// 16-bit to zero.
func (d *DNSQueryManager) setNumOfAdditionalRRs() {
	d.query = append(d.query, 0b00000000, 0b00000000)
}

func (d *DNSQueryManager) setHeader() {
	d.setTransactionId()
	d.setStandardFlags()
	d.setNumOfQuestions()
	d.setNumOfAnswers()
	d.setNumOfAuthorityRRs()
	d.setNumOfAdditionalRRs()
}

func (d *DNSQueryManager) setQuestionName() {
	for _, part := range d.domainSegments {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			log.Fatalln("invalid input domain")
		}
		// Each label prefixed by a length byte, followed by the label itself
		d.query = append(d.query, byte(len(part)))
		d.query = append(d.query, []byte(part)...)
	}
	d.query = append(d.query, 0b00000000) // To show that this is end of the domain
}

func (d *DNSQueryManager) setQuestionType() {
	// ** QTYPE -> Type A (host address) - 2 bytes (A, AAAA, MX, etc.)
	d.query = append(d.query, 0b00000000, 0b00000001)
	// Query for IPv6 (AAAA record)
	// query := append(query, 0b00000000, 0b00011100)
}

func (d *DNSQueryManager) setQuestionClass() {
	// ** QCLASS -> Class IN (Internet) - 2 bytes
	d.query = append(d.query, 0b00000000, 0b00000001)
}

func (d *DNSQueryManager) setQuestion() {
	d.setQuestionName()
	d.setQuestionType()
	d.setQuestionClass()
}

func (d *DNSQueryManager) CreateQuery() []byte {
	d.setHeader()
	d.setQuestion()
	return d.query
}
