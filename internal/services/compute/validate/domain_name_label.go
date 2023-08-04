// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func OrchestratedDomainNameLabel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// Domain labels can be 63 characters long per the Network API, the compute team adds a dash and a UUID when deploying to multiple
	// Zones which causes a validation error in the RP. Updating the validation code to be artificially constrictive to account for the
	// RPs behavior...
	if matched := regexp.MustCompile(`^[a-z][a-z0-9-]{1,24}[a-z0-9]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 - 26 characters long, start with a lower case letter, end with a lower case letter or number and contains only a-z, 0-9 and hyphens", k))
	}
	return
}
