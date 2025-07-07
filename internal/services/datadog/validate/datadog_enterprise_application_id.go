// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DatadogEnterpriseApplicationID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/25904
	if !regexp.MustCompile(`[a-z0-9-]$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("expected value of %s not match regular expression, got %v", k, v))
		return
	}
	if len(v) > 42 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 42))
		return
	}

	return
}
