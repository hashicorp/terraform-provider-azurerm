// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	validateResourceGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ResourceGroupConsumptionBudget struct {
	base consumptionBudgetBaseResource
}

var (
	_ sdk.Resource                   = ResourceGroupConsumptionBudget{}
	_ sdk.ResourceWithCustomImporter = ResourceGroupConsumptionBudget{}
)

func (r ResourceGroupConsumptionBudget) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateResourceGroup.ResourceGroupID,
		},
	}
	return r.base.arguments(schema)
}

func (r ResourceGroupConsumptionBudget) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ResourceGroupConsumptionBudget) ModelObject() interface{} {
	return nil
}

func (r ResourceGroupConsumptionBudget) ResourceType() string {
	return "azurerm_consumption_budget_resource_group"
}

func (r ResourceGroupConsumptionBudget) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r ResourceGroupConsumptionBudget) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "resource_group_id")
}

func (r ResourceGroupConsumptionBudget) Read() sdk.ResourceFunc {
	return r.base.readFunc("resource_group_id")
}

func (r ResourceGroupConsumptionBudget) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ResourceGroupConsumptionBudget) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r ResourceGroupConsumptionBudget) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc()
}
