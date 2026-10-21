// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_static_member

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticmembers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_network_manager_static_member"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticmembers.ValidateStaticMemberID
}
