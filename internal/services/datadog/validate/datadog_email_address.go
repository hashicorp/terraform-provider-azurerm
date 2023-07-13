// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DatadogMonitorsEmailAddress(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if !regexp.MustCompile(`^[A-Za-z0-9._%+-]+@(?:[A-Za-z0-9-]+\.)+[A-Za-z]{2,}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("expected value of %s not match regular expression, got %v", k, v))
		return
	}
	return
}
