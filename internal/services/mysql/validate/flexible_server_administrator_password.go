// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerAdministratorPassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 8 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 8, v))
		return
	}

	if len(v) > 128 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 128, v))
		return
	}

	switch {
	case regexp.MustCompile(`^.*[a-z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[A-Z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[0-9]+.*$`).MatchString(v):
		return
	case regexp.MustCompile(`^.*[a-z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[A-Z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[\W]+.*$`).MatchString(v):
		return
	case regexp.MustCompile(`^.*[a-z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[\W]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[0-9]+.*$`).MatchString(v):
		return
	case regexp.MustCompile(`^.*[A-Z]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[\W]+.*$`).MatchString(v) && regexp.MustCompile(`^.*[0-9]+.*$`).MatchString(v):
		return
	default:
		errors = append(errors, fmt.Errorf("%q must contain characters from three of the categories – uppercase letters, lowercase letters, numbers and non-alphanumeric characters, got %v", k, v))
		return
	}
}
