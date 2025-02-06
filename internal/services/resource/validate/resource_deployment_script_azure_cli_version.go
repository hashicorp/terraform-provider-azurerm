// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ResourceDeploymentScriptAzureCliVersion(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	var errors []error
	if matched := regexp.MustCompile(`^\d+\.\d+\.\d+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q should be in the format `X.Y.Z` (e.g. `2.30.0`)", k))
	}

	return nil, errors
}
