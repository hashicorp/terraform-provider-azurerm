// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
)

func DatadogMonitorsPhoneNumber(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if len(v) > 40 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 40))
		return
	}
	return
}
