// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceMacAddress(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must be in format `00:00:00:00:00:00`", k))
	}

	return warnings, errors
}
