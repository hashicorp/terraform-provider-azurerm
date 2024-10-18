// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consumption

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	validateManagementGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupConsumptionBudget struct {
	base consumptionBudgetBaseResource
}

var (
	_ sdk.Resource                   = ManagementGroupConsumptionBudget{}
	_ sdk.ResourceWithCustomImporter = ManagementGroupConsumptionBudget{}
)

func (r ManagementGroupConsumptionBudget) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateManagementGroup.ManagementGroupID,
		},

		// Consumption Budgets for Management Groups have a different notification schema,
		// here we override the notification schema in the base resource
		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"threshold": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"threshold_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(budgets.ThresholdTypeActual),
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.ThresholdTypeActual),
							"Forecasted",
						}, false),
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.OperatorTypeEqualTo),
							string(budgets.OperatorTypeGreaterThan),
							string(budgets.OperatorTypeGreaterThanOrEqualTo),
						}, false),
					},

					"contact_emails": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
	return r.base.arguments(schema)
}

func (r ManagementGroupConsumptionBudget) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r ManagementGroupConsumptionBudget) ModelObject() interface{} {
	return nil
}

func (r ManagementGroupConsumptionBudget) ResourceType() string {
	return "azurerm_consumption_budget_management_group"
}

func (r ManagementGroupConsumptionBudget) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return budgets.ValidateScopedBudgetID
}

func (r ManagementGroupConsumptionBudget) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "management_group_id")
}

func (r ManagementGroupConsumptionBudget) Read() sdk.ResourceFunc {
	return r.base.readFunc("management_group_id")
}

func (r ManagementGroupConsumptionBudget) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r ManagementGroupConsumptionBudget) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r ManagementGroupConsumptionBudget) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc()
}
