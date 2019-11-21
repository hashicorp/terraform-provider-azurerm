package netapp

import (
	"fmt"
	"net"
	"regexp"
)

func ValidateNetAppAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-_\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{2,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-_\da-zA-Z]{0,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length and start with letters and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeCreationToken(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{0,79}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length and start with letters and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeAllowedClients(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must start with IPV4 address and/or slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8. Got %q.", k, value))
	}

	ip := net.ParseIP(value)
	if four := ip.To4(); four == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 address: %q", k, v))
	}

	if len(errors) == 2 {
		return warnings, errors
	} else {
		return warnings, nil
	}
}
