// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func StorageSyncName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[0-9a-zA-Z-_. ]*[0-9a-zA-Z-_]$").MatchString(input) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of letters, numbers, spaces, and any of the following characters: '.-_' and that does not end with characters: '. '", input))
	}

	return warnings, errors
}
