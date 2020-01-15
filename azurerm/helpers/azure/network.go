package azure

import (
	"fmt"
	"net"
)

// CompressIPv6Address creates a canonicalized IPv6 address according to RFC5952
func CompressIPv6Address(ipv6address string) (string, error) {
	// Validate, if IPv6 address
	ip := net.ParseIP(ipv6address)
	if six := ip.To16(); six == nil {
		return "", fmt.Errorf("%q is not a valid IPv6 address", ipv6address)
	}

	return ip.String(), nil
}
