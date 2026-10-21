// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_security_admin_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securityadminconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_network_manager_security_admin_configuration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securityadminconfigurations.ValidateSecurityAdminConfigurationID
}
