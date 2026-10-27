// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_eventgrid_partner_configuration"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateResourceGroupID
}
