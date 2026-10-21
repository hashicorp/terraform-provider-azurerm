// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_network_security_perimeter_access_rule"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networksecurityperimeteraccessrules.ValidateAccessRuleID
}
