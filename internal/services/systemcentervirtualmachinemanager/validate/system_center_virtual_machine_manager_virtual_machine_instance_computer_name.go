// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^[a-zA-Z0-9]{1,}$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must only contain alphanumeric characters", k))
	}

	return warnings, errors
}
