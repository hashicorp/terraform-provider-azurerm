// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net/mail"
)

func Email(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if _, err := mail.ParseAddress(v); err != nil {
		return []string{}, append(errors, fmt.Errorf("%v must be a valid email address", k))
	}

	return []string{}, []error{}
}
