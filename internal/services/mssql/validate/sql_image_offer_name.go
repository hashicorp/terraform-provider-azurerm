// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func SqlImageOfferName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^SQL[A-Za-z0-9]*-WS[A-Za-z0-9]*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q should be in the form SQL<SQLversion>-WS<OSversion>, for example SQL2019-WS2019: %q", k, value))
	}

	return warnings, errors
}
