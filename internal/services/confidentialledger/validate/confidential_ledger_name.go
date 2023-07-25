// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ConfidentialLedgerName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 32 {
		errors = append(errors, fmt.Errorf("%q may not exceed 32 characters in length", k))
	}

	if strings.HasPrefix(value, "-") {
		errors = append(errors, fmt.Errorf("%q may not start with a dash", k))
	}

	if strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q may not end with a dash", k))
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be blank", k))
	} else if matched := regexp.MustCompile(`^[^\-][A-Za-z0-9\-]{1,33}[^\-]$`).Match([]byte(value)); !matched {
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dash", k))
	}

	return warnings, errors
}
