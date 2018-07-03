package validate

import (
	"fmt"
	"net"
)

func IP4Address(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
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
