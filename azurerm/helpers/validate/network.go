package validate

import (
	"fmt"
	"net"
)

func IPv4Address(i interface{}, k string) (_ []string, errors []error) {
	return validateIpv4Address(i, k, false)
}

func IPv4AddressOrEmpty(i interface{}, k string) (_ []string, errors []error) {
	return validateIpv4Address(i, k, true)
}

func validateIpv4Address(i interface{}, k string, allowEmpty bool) (_ []string, errors []error) {
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
		errors = append(errors, fmt.Errorf("%q is not a valid IP4 address: %q", k, v))
	}

	return
}

func MACAddress(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := net.ParseMAC(v); err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid MAC address: %q (%v)", k, i, err))
	}

	return
}

func PortNumber(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be int", k))
		return
	}

	if v < 0 || 65535 < v {
		errors = append(errors, fmt.Errorf("%q is not a valid port number: %q", k, i))
	}

	return
}
