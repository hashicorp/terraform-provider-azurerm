// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func BatchAccountIpRange(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	re := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|30))?$`)
	if re != nil && !re.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must start with IPV4 address and/or slash, number of bits (0-30) as prefix. Example: 127.0.0.1/8. Got %q.", k, value))
	}

	return warnings, errors
}
