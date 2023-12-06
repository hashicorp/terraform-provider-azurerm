// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func StorageTargetName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	p := regexp.MustCompile(`^[-0-9a-zA-Z_]{1,31}$`)
	if !p.MatchString(v) {
		errors = append(errors, fmt.Errorf("%q can contain alphanumeric characters, dashes and underscores and has to be between 1 and 31 characters", k))
	}

	return warnings, errors
}
