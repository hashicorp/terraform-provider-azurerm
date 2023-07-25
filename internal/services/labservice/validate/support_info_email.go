// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SupportInfoEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q is not a valid email", k))
	}
	return warnings, errors
}
