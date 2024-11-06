package dns

const (
	// Minimum capacity needed for query slice
	// sent to DNS
	// header 12 + QTYPE 2 + QCLASS 2
	minQueryCap = 16
)

const (
	IpTypeV4 uint8 = iota
	IpTypeV6
)

const (
	startingQueryIdx       int = 12
	endOfQuestionIndicator int = 0
	byteForLength          int = 1
)
