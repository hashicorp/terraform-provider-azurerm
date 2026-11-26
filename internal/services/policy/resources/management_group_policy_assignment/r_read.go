// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_assignment

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return r.base.ReadFunc("management_group_id")
}
