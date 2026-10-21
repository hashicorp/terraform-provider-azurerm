// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_strategy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleetupdatestrategies.ValidateUpdateStrategyID
}

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_fleet_update_strategy"
}
