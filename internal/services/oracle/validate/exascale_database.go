// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ExascaleDatabaseResourceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The value must start with a letter or an underscore (_).
	// The value can only include letters, numbers, underscores (_), and hyphens (-).
	// The value cannot contain consecutive hyphens (--).
	if !regexp.MustCompile(`^[a-zA-Z_]([a-zA-Z0-9_]*(-[a-zA-Z0-9_]+)*-?)?$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must begin with a letter or underscore (_), contain only letters, numbers, underscores (_), and hyphens (-), and cannot contain consecutive hyphens (--)", k))
	}

	return
}
