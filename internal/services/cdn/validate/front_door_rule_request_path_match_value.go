// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func FrontDoorRuleRequestPathMatchValue(i interface{}, k string) (_ []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected `%s` to be a string", k)}
	}

	if value == "" {
		return nil, []error{fmt.Errorf("`%s` must not be empty", k)}
	}

	if strings.HasPrefix(value, "/") && value != "/" {
		return nil, []error{fmt.Errorf("`%s` must not begin with `/` unless the value is `/`, got `%s`", k, value)}
	}

	return nil, nil
}
