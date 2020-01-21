package validate

import (
	"fmt"
	"net"
	"strings"
)

// IPv4Address validation returns true if the given address is a valid IPv4 address
func IPv4Address(i interface{}, k string) (warnings []string, errors []error) {
	return validateIpv4Address(i, k, false)
}

// IPv4AddressOrEmpty validation returns true if the given address is a valid IPv4 address or empty
func IPv4AddressOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	return validateIpv4Address(i, k, true)
}

func validateIpv4Address(i interface{}, k string, allowEmpty bool) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" && allowEmpty {
		return
	}

	ip := net.ParseIP(v)
	if four := ip.To4(); four == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 address: %q", k, v))
	}

	return warnings, errors
}

// IPv6Address validation returns true if the given address is a valid IPv6 address
func IPv6Address(i interface{}, k string) (warnings []string, errors []error) {
	return validateIpv6Address(i, k, false)
}

// IPv6AddressOrEmpty validation returns true if the given address is a valid IPv6 address or empty
func IPv6AddressOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	return validateIpv6Address(i, k, true)
}

func validateIpv6Address(i interface{}, k string, allowEmpty bool) (warnings []string, errors []error) { // nolint: unparam
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" && allowEmpty {
		return
	}

	ip := net.ParseIP(v)
	if six := ip.To16(); six == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv6 address: %q", k, v))
	}

	return warnings, errors
}

// IPAddress validation returns true if the given address is a valid IPv4 or IPv6 address
func IPAddress(i interface{}, k string) (warnings []string, errors []error) {
	return validateIPAddress(i, k, false)
}

// IPAddressOrEmpty validation returns true if the given address is a valid IPv4 or IPv6 address, or empty
func IPAddressOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	return validateIPAddress(i, k, true)
}

func validateIPAddress(i interface{}, k string, allowEmpty bool) (warnings []string, errors []error) { // nolint: unparam
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if v == "" && allowEmpty {
		return
	}

	ip := net.ParseIP(v)
	if ip == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 or IPv6 address: %q", k, v))
	}

	return warnings, errors
}

// CIDR validations allow IPv4 and IPv6 CIDR notations as defined in RFC 4632 and RFC 4291
func CIDR(i interface{}, k string, allowNoSuffix bool) (warnings []string, errors []error) {
	cidr := i.(string)

	// A CIDR without the suffix is invalid, but for compatibility reasons we need
	// to be able to validate those as valid, too.
	if allowNoSuffix && !strings.Contains(cidr, "/") {
		cidr = cidr + "/0"
	}
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s must be a valid IPv4 or IPv6 address and slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8. Got %q", k, cidr))
	}

	return warnings, errors
}

// MACAddress validates the input to be a valid MAC address
// Valid are addresses according to IEEE 802 MAC-48, EUI-48, EUI-64, or a 20-octet IP over InfiniBand link-layer address
func MACAddress(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := net.ParseMAC(v); err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid MAC address: %q (%v)", k, i, err))
	}

	return warnings, errors
}

func PortNumber(i interface{}, k string) (warnings []string, errors []error) {
	return validatePortNumber(i, k, false)
}

func PortNumberOrZero(i interface{}, k string) (warnings []string, errors []error) {
	return validatePortNumber(i, k, true)
}

func validatePortNumber(i interface{}, k string, allowZero bool) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be int", k))
		return
	}

	if allowZero && v == 0 {
		return
	}

	if v < 1 || 65535 < v {
		errors = append(errors, fmt.Errorf("%q is not a valid port number: %d", k, v))
	}

	return warnings, errors
}
