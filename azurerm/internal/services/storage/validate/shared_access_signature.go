package validate

import (
	"fmt"
	"net"
	"strings"
)

func SharedAccessSignatureIP(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if net.ParseIP(value) != nil {
		return warnings, errors
	}

	ipRange := strings.Split(value, "-")

	if len(ipRange) != 2 || net.ParseIP(ipRange[0]) == nil || net.ParseIP(ipRange[1]) == nil {
		errors = append(errors, fmt.Errorf("%q must be a valid ipv4 address or a range of ipv4 addresses separated by a hyphen", k))
		return warnings, errors
	}

	ip1 := ipRange[0]
	ip2 := ipRange[1]

	if ip1 == ip2 {
		errors = append(errors, fmt.Errorf("IP addresses in a range for %q must be not be identical", k))
		return warnings, errors
	}

	return warnings, errors
}
