// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerAdministratorLogin(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("length should equal to or greater than %d, got %q", 1, v))
		return
	}

	if len(v) > 32 {
		errors = append(errors, fmt.Errorf("length should be equal to or less than %d, got %q", 32, v))
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must only contains characters, numbers or '_', got %v", k, v))
		return
	}

	if v == "azure_superuser" || v == "admin" || v == "administrator" || v == "root" || v == "guest" || v == "public" {
		errors = append(errors, fmt.Errorf("%q cannot be `azure_superuser`, `admin`, `administrator`, `root`, `guest` or `public`", k))
		return
	}

	return
}
