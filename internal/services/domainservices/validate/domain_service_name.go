// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func DomainServiceName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	p := regexp.MustCompile(`^(([0-9a-zA-Z])|(([0-9a-zA-Z][0-9a-zA-Z-]{0,28}[0-9a-zA-Z])))(\.[0-9a-zA-Z-]+)+$`)
	if !p.MatchString(v) {
		errors = append(errors, fmt.Errorf("domain_name must be a valid FQDN and the first element must be 15 or fewer characters"))
	}

	return
}
