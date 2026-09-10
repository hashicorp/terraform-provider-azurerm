// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func SkuProfileVMSizeName(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected `%s` to be a string", key))
		return warnings, errors
	}

	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("property `%s` cannot be empty", key))
		return warnings, errors
	}

	if !strings.HasPrefix(v, "Standard_") {
		errors = append(errors, fmt.Errorf("property `%s` must begin with `Standard_`, got `%s`", key, v))
		return warnings, errors
	}

	family := strings.TrimPrefix(v, "Standard_")

	if strings.HasPrefix(family, "DC") || strings.HasPrefix(family, "EC") {
		errors = append(errors, fmt.Errorf("property `%s` must use a supported instance mix VM size family (`A`, `B`, `D`, `E`, or `F`), got `%s`", key, v))
		return warnings, errors
	}

	if strings.HasPrefix(family, "A") || strings.HasPrefix(family, "B") || strings.HasPrefix(family, "D") || strings.HasPrefix(family, "E") || strings.HasPrefix(family, "F") {
		return warnings, errors
	}

	errors = append(errors, fmt.Errorf("property `%s` must use a supported instance mix VM size family (`A`, `B`, `D`, `E`, or `F`), got `%s`", key, v))

	return warnings, errors
}
