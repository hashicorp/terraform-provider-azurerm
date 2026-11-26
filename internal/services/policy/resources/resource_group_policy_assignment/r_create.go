// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_policy_assignment

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return r.base.CreateFunc(r.ResourceType(), "resource_group_id")
}
