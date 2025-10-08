// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ManagedRedisDatabaseGeoreplicationGroupName(val interface{}, argName string) ([]string, []error) {
	v, ok := val.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", argName)}
	}

	// Check length between 1-63
	if len(v) < 1 || len(v) > 63 {
		return nil, []error{fmt.Errorf("%q must be between 1 and 63 characters long", argName)}
	}

	// Check if it contains only letters, numbers, and hyphens
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(v) {
		return nil, []error{fmt.Errorf("%q can only contain letters, numbers, and hyphens", argName)}
	}

	// Check if first and last characters are letters or numbers
	firstChar := v[0]
	lastChar := v[len(v)-1]
	if (firstChar < 'a' || firstChar > 'z') && (firstChar < 'A' || firstChar > 'Z') && (firstChar < '0' || firstChar > '9') {
		return nil, []error{fmt.Errorf("%q must start with a letter or number", argName)}
	}
	if (lastChar < 'a' || lastChar > 'z') && (lastChar < 'A' || lastChar > 'Z') && (lastChar < '0' || lastChar > '9') {
		return nil, []error{fmt.Errorf("%q must end with a letter or number", argName)}
	}

	// Check for consecutive hyphens
	if strings.Contains(v, "--") {
		return nil, []error{fmt.Errorf("%q cannot contain consecutive hyphens", argName)}
	}

	return nil, nil
}
