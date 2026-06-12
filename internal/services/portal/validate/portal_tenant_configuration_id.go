// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/tenantconfigurations"
)

func PortalTenantConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := tenantconfigurations.ParseTenantConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
