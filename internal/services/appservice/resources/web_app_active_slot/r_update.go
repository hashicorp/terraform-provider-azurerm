// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

// Note: `Update` re-uses `Create` as there is no actual resource being managed, this meta-resource simply triggers a
// swap operations between the named slot and `Production`. Without this changing which slot is `Active` would result in
// Terraform deleting and recreating this resource, which may cause concern that the operation is somehow destructive.

func (r Resource) Update() sdk.ResourceFunc {
	return r.Create()
}
