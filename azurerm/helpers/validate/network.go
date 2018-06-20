package validate

import (
	"fmt"
	"regexp"
)

func Ip4Address(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`((^[0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q is not a valid IP4 address: '%q`", k, i))
	}

	return
}

func MacAddress(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`((^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q is not a valid MAC address: '%q`", k, i))
	}

	return
}
