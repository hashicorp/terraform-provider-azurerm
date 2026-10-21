// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule_collection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_network_manager_admin_rule_collection"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return adminrulecollections.ValidateSecurityAdminConfigurationRuleCollectionID
}
