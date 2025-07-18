// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

const autorizationAAD = "Authorization=AAD"

func ApplicationInsightsAuthenticationString(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)

	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return warnings, errors
	}

	var matched bool = false

	if !strings.HasPrefix(v, autorizationAAD) {
		errors = append(errors, fmt.Errorf("%q must always begin with %q, got: %q", key, autorizationAAD, v))
		return warnings, errors
	} else if v == autorizationAAD {
		matched = true
	} else {
		parts := strings.Split(v, ";")
		if len(parts) == 2 {
			clientIDString := parts[1]
			matched = regexp.MustCompile(`^ClientId=([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`).Match([]byte(clientIDString))
		}
	}

	if !matched {
		errors = append(errors, fmt.Errorf("%q must be in the format \"Authorization=AAD;ClientId=<GUID>\" when used with user-assigned managed identity, got: %q", key, v))
	}

	return warnings, errors
}
