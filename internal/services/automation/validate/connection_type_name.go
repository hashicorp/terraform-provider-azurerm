// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ConnectionTypeName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if !regexp.MustCompile(`^[\w\-]{1,128}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s contain only letters, numbers hyphens and underscore. The value must be between 1 and 128 characters long", k))
	}

	return nil, errors
}
