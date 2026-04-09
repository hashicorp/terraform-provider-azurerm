// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import (
	"strings"
)

// IsAzureStack returns whether the current Environment is an Azure Stack Environment
// this can be contextually useful since the Azure Stack implementation differs slightly
// from the other Azure Environments, particularly around Authentication.
func (e *Environment) IsAzureStack() bool {
	if strings.EqualFold(e.Name, "AzureStackCloud") {
		return true
	}

	if !strings.EqualFold(e.Authorization.IdentityProvider, "AAD") || !strings.EqualFold(e.Authorization.Tenant, "common") {
		return true
	}

	return false
}
