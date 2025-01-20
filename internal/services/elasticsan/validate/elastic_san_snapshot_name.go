// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ElasticSanSnapshotName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,61}[a-z0-9]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 63 characters. It can contain only lowercase letters, numbers, underscores (_) and hyphens (-). It must start and end with a lowercase letter or number", k))
	}

	if matched := regexp.MustCompile(`[_-][_-]`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q must have hyphens and underscores be surrounded by alphanumeric character", k))
	}

	return warnings, errors
}
