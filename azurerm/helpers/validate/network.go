package validate

import (
	"fmt"
	"net"
	"regexp"
)

func Ip4Address(i interface{}, k string) (_ []string, errors []error) {
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

func MacAddress(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q is not a valid MAC address: %q", k, i))
	}

	return
}
