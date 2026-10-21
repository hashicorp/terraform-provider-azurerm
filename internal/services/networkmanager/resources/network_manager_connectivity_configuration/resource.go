// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_network_manager_connectivity_configuration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectivityconfigurations.ValidateConnectivityConfigurationID
}
