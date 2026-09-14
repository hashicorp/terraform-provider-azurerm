// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// CIDR is a SchemaValidateFunc which tests if the provided value is a valid IPv4 CIDR
func CIDR(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$`), "must start with IPV4 address and/or slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8")(i, k)
}

func IPv4Address(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	ip := net.ParseIP(v)
	if four := ip.To4(); four == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 address: %q", k, v))
	}

	return warnings, errors
}

var (
	PortNumber       = validation.IsPortNumber
	PortNumberOrZero = validation.IsPortNumberOrZero
)
