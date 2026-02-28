// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// StorageDiscoveryWorkspaceRoot validates that the value is a valid Subscription ID or Resource Group ID.
// Per the portal behaviour, only subscription and resource group scope IDs are allowed.
func StorageDiscoveryWorkspaceRoot(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if _, err := commonids.ParseSubscriptionID(input); err == nil {
		return warnings, errors
	}

	if _, err := commonids.ParseResourceGroupID(input); err == nil {
		return warnings, errors
	}

	errors = append(errors, fmt.Errorf("each value in `workspace_root` must be a valid Subscription ID (e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`) or Resource Group ID (e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg`)"))
	return warnings, errors
}
