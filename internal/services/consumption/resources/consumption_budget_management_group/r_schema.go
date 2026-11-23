// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_management_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Management Group Consumption Budget.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"management_group_id": {
			Description:  "The ID of the Management Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagementGroupID,
		},

		// Consumption Budgets for Management Groups have a different notification schema,
		// here we override the notification schema in the base resource
		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Should the notification be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"threshold": {
						Description:  "Threshold value associated with a notification. Notification is sent when the cost exceeds the threshold. It is always percent and has to be between 0 and 1000.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"threshold_type": {
						Description: "The type of threshold for the notification. This determines whether the notification is triggered by forecasted costs or actual costs.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(budgets.ThresholdTypeActual),
						ValidateFunc: validation.StringInSlice([]string{
							string(budgets.ThresholdTypeActual),
							"Forecasted",
						}, false),
					},
					"operator": {
						Description:  "The comparison operator for the notification.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(budgets.PossibleValuesForOperatorType(), false),
					},

					"contact_emails": {
						Description: "A list of email addresses to send the budget notification to when the threshold is exceeded.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
	return r.base.Arguments(schema)
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.Attributes()
}
