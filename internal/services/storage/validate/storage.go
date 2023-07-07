// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func StorageShareDirectoryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	// Per: https://learn.microsoft.com/en-us/rest/api/storageservices/naming-and-referencing-shares--directories--files--and-metadata#directory-and-file-names
	if regexp.MustCompile(`^\.+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%s must not only contain dots`, k))
	}
	// Note that we didn't forbid the forward slash in the non-head segment here as it seems to be allowed as the directory name for constructing directory hierarchy.
	if !regexp.MustCompile(`^[^"/\:|<>*?]+(/[^"\:|<>*?]+)*$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%s must not contain following characters: "\:|<>*?`, k))
	}

	return warnings, errors
}
