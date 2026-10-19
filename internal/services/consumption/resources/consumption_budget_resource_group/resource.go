// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_resource_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct {
	base helpers.ConsumptionBudgetBaseResource
}

var (
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_consumption_budget_resource_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return r.base.ImporterFunc()
}
