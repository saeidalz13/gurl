package conninfo

import (
	"fmt"
	"strconv"
)

func convertIpStringToBytes(ipStrSegments []string, ip *[]byte) error {
	for _, b := range ipStrSegments {
		n, err := strconv.Atoi(b)
		if err != nil {
			return err
		}

		if n < 0 || n > 255 {
			return fmt.Errorf("ip segment > 255 or < 0: %d", n)
		}

		*ip = append(*ip, byte(n))
	}

	return nil
}
