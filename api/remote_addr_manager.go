package api

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/saeidalz13/gurl/internal/errutils"
	"github.com/saeidalz13/gurl/internal/httpconstants"
)

type RemoteAddrManager struct {
	domain     string
	ipCacheDir string
}

func newRemoteAddrManager(ipCacheDir, domain string) RemoteAddrManager {
	return RemoteAddrManager{
		domain:     domain,
		ipCacheDir: ipCacheDir,
	}
}

func (ram RemoteAddrManager) isDomainLocalHost() bool {
	return strings.Contains(ram.domain, "localhost") || strings.Contains(ram.domain, "127.0.0.1")
}

func (ram RemoteAddrManager) extractPort() (int, error) {
	domainSegments := strings.Split(ram.domain, ":")

	if len(domainSegments) != 2 {
		return 0, fmt.Errorf("domain must be in format of ip:port")
	}

	return strconv.Atoi(domainSegments[1])
}

func (ram RemoteAddrManager) fetchCachedIp() (net.IP, error) {
	domainFile := filepath.Join(ram.ipCacheDir, ram.domain)
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

func (ram RemoteAddrManager) cacheDomainIp(ipStr string) error {
	domainFile := filepath.Join(ram.ipCacheDir, ram.domain)
	// 0o600 read and write permissions only for the owner
	f, err := os.OpenFile(domainFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}

	_, err = f.WriteString(ipStr)
	return err
}

// bool shows if the connection should be TLS
func (ram RemoteAddrManager) resolveConnectionInfo() (net.IP, int, bool) {
	if ram.isDomainLocalHost() {
		ip := net.IPv4(127, 0, 0, 1)
		port, err := ram.extractPort()
		errutils.CheckErr(err)
		return ip, port, false
	}

	// If not localhost, the IP needs to be fetched
	// from DNS server. We cache the data to prevent
	// unnecessary network I/O.
	ip, err := ram.fetchCachedIp()
	if err != nil {
		ip = newDNSResolver(ram.domain).mustResolveIP()
		if err := ram.cacheDomainIp(ip.String()); err != nil {
			// Should not stop the operation
			fmt.Printf("skipped ip caching: %v\n", err)
		}
	}
	return ip, httpconstants.PortHTTPS, true
}
