// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func ParameterNames(v any, _ string) (warnings []string, errors []error) {
	m := v.(map[string]any)
	for k := range m {
		if k != strings.ToLower(k) {
			errors = append(errors, fmt.Errorf("due to a bug in the implementation of Runbooks in Azure, the parameter names need to be specified in lowercase only. See: \"https://github.com/Azure/azure-sdk-for-go/issues/4780\" for more information"))
		}
	}

	return warnings, errors
}
