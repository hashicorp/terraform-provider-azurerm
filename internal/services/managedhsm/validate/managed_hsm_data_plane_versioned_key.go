// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
)

func ManagedHSMDataPlaneVersionedKeyID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return warnings, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if _, err := parse.ManagedHSMDataPlaneVersionedKeyID(v, nil); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a Managed HSM Data Plane Versioned Key ID: %+v", v, err))
	}

	return warnings, errors
}
