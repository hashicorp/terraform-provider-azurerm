package validate

import (
	"bytes"
	"fmt"
	"net"
)

func IPAddressInAzureReservedAPIPARange(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	ip := net.ParseIP(v)
	if four := ip.To4(); four == nil {
		errors = append(errors, fmt.Errorf("expected %s to contain a valid IPv4 address, got: %s", k, v))
	}

	// See: https://docs.microsoft.com/en-us/azure/vpn-gateway/bgp-howto#2-create-the-vpn-gateway-for-testvnet1-with-bgp-parameters
	azureAPIPAStart := net.ParseIP("169.254.21.0")
	azureAPIPAEnd := net.ParseIP("169.254.22.255")

	if !(bytes.Compare(ip, azureAPIPAStart) >= 0 && bytes.Compare(ip, azureAPIPAEnd) <= 0) {
		errors = append(errors, fmt.Errorf("%s is not within Azure reserved APIPA range: [%s, %s]", ip, azureAPIPAStart, azureAPIPAEnd))
		return warnings, errors
	}

	return nil, nil
}
