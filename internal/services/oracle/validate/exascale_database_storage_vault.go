// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ExascaleDatabaseStorageVaultName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if strings.Contains(v, "--") || !regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must begin with a letter or underscore (_), contain only letters, numbers, underscores (_) and any consecutive hyphers (--)", k))
	}

	return
}
