package validate

import (
	"fmt"
	"net"
	"regexp"
)

// CIDR is a SchemaValidateFunc which tests if the provided value is a valid IPv4 CIDR
func CIDR(i interface{}, k string) (warnings []string, errors []error) {
	cidr := i.(string)

	re := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$`)
	if re != nil && !re.MatchString(cidr) {
		errors = append(errors, fmt.Errorf("%s must start with IPV4 address and/or slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8. Got %q.", k, cidr))
	}

	return warnings, errors
}

func IPv4Address(i interface{}, k string) (warnings []string, errors []error) {
	return validateIpv4Address(i, k, false)
}

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
