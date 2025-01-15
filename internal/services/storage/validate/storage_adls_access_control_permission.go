// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/jackofallops/giovanni/storage/accesscontrol"
)

func ADLSAccessControlPermissions(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}
	if err := accesscontrol.ValidateACEPermissions(v); err != nil {
		errors = append(errors, fmt.Errorf("value of %s not valid: %s", k, err))
		return warnings, errors
	}
	return warnings, errors
}
