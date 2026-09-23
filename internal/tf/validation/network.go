// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// IsCIDRIPv4 is a SchemaValidateFunc which tests if the provided value is a valid IPv4 CIDR.
// IsCIDR accepts IPv6 as well, so this exists for the properties that only take IPv4.
func IsCIDRIPv4(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$`), "must start with IPV4 address and/or slash, number of bits (0-32) as prefix. Example: 127.0.0.1/8")(i, k)
}
