package conninfo

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/saeidalz13/gurl/api/dns"
	"github.com/saeidalz13/gurl/internal/domainparser"
	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
	"github.com/saeidalz13/gurl/models"
)

type ConnInfoResolver struct {
	protocol       uint8
	domain         string
	ipCacheDir     string
	domainSegments []string
}

func NewConnInfoResolver(ipCacheDir, domain string, domainSegments []string, protocol uint8) ConnInfoResolver {
	return ConnInfoResolver{
		domain:         domain,
		ipCacheDir:     ipCacheDir,
		domainSegments: domainSegments,
		protocol:       protocol,
	}
}

func (c ConnInfoResolver) isDomainLocalHost() bool {
	return strings.Contains(c.domain, "localhost") || strings.Contains(c.domain, "127.0.0.1")
}

func (c ConnInfoResolver) extractPort() (int, error) {
	domainSegments := strings.Split(c.domain, ":")

	if len(domainSegments) != 2 {
		return 0, fmt.Errorf("domain must be in format of ip:port")
	}

	return strconv.Atoi(domainSegments[1])
}

func (c ConnInfoResolver) fetchCachedIp() (net.IP, error) {
	domainFile := filepath.Join(c.ipCacheDir, c.domain)
	f, err := os.OpenFile(domainFile, os.O_RDONLY, 0o600)
	if err != nil {
		return nil, err
	}

	// ip v4 string is xxx.xxx.xxx.xxx
	// max is 4*3 bytes + 3 dots (bytes) = 15 bytes
	buf := make([]byte, 15)
	n, err := f.Read(buf)
	if err != nil {
		return nil, err
	}

	ipStrSegments := strings.Split(string(buf[:n]), ".")
	ipBytes := make([]byte, 0, 4)

	for _, b := range ipStrSegments {
		n, err := strconv.Atoi(b)
		if err != nil {
			return nil, err
		}

		if n < 0 || n > 255 {
			return nil, fmt.Errorf("ip segment > 255 or < 0: %d", n)
		}

		ipBytes = append(ipBytes, byte(n))
	}

	if len(ipBytes) != 4 {
		return nil, fmt.Errorf("cached ip is incorrect")
	}

	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]), nil
}

func (c ConnInfoResolver) cacheDomainIp(ipStr string) error {
	domainFile := filepath.Join(c.ipCacheDir, c.domain)

	// TODO: Decide which one to use for writing
	// Method 1: (commented out for now)
	// 0o600 read and write permissions only for the owner.
	//
	// os.O_EXCL causes `OpenFile` to give error if file already exists.
	// https://man7.org/linux/man-pages/man2/open.2.html
	// f, err := os.OpenFile(domainFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	// if err != nil {
	// 	return err
	// }
	// _, err = f.WriteString(ipStr)
	// return err

	// Method 2:
	return os.WriteFile(domainFile, []byte(ipStr), 0o600)
}

// bool shows if the connection should be TLS
func (c ConnInfoResolver) Resolve() models.ConnInfo {
	if c.isDomainLocalHost() {
		ip := net.IPv4(127, 0, 0, 1)
		port, err := c.extractPort()
		errutils.CheckErr(err)
		return models.ConnInfo{
			IP:     ip,
			IPType: 0,
			Port:   port,
			IsTls:  false,
		}
	}

	// If not localhost, the IP needs to be fetched
	// from DNS server. We cache the data to prevent
	// unnecessary network I/O.
	var ipType uint8
	ip, err := c.fetchCachedIp()
	if err != nil {
		ip, ipType = dns.MustResolveIP(c.domainSegments)
		if err := c.cacheDomainIp(ip.String()); err != nil {
			// Should not stop the operation
			fmt.Printf("skipped ip caching: %v\n", err)
		}
	}
	return models.ConnInfo{
		IP:     ip,
		IPType: ipType,
		Port:   httpconstants.PortHTTPS,
		IsTls:  c.protocol == domainparser.ProtocolHTTPS,
	}
}
