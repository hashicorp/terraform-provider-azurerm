// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// CustomDataWriteOnly validates that the value is Base64 encoded. Unlike `validation.StringIsBase64`
// the error omits the value, as a write-only value is a secret that must not reach the plan output.
func CustomDataWriteOnly(i any, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	v, _ := i.(string)
	if _, err := base64.StdEncoding.DecodeString(v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a base64 string", k))
	}

	return warnings, errors
}
