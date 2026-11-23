// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_capability

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateChaosStudioCapabilityID
}

func (r Resource) ResourceType() string {
	return "azurerm_chaos_studio_capability"
}
