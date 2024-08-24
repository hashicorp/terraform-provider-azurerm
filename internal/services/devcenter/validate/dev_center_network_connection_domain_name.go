// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DevCenterNetworkConnectionDomainName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if !regexp.MustCompile("^[a-zA-Z0-9-]([a-zA-Z0-9-.]{0,253}[a-zA-Z0-9-])?$").MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must start or end with an alphanumeric character or dashes, may contain alphanumeric characters, dashes or periods and must be between 1 and 255 characters long", k))
	}

	return warnings, errors
}
