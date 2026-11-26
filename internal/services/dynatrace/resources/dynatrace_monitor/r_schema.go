// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_monitor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "Name of the Dynatrace monitor.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"monitoring_enabled": {
			Description: "Flag specifying if the resource monitoring is enabled or disabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},

		"marketplace_subscription": {
			Description:  "Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(monitors.PossibleValuesForMarketplaceSubscriptionStatus(), false),
		},

		"identity": commonschema.SystemAssignedIdentityRequired(),

		"plan": {
			Description: "Billing plan information.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// Issue: https://github.com/Azure/azure-rest-api-specs/issues/31284
					// This should be an enum.
					"billing_cycle": {
						Description: "Different billing cycles.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"MONTHLY",
							"WEEKLY",
							"YEARLY",
						}, false),
					},

					"plan": {
						Description:  "Plan ID as published by Dynatrace.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					// Issue: https://github.com/Azure/azure-rest-api-specs/issues/31284
					// This should be an enum.
					"usage_type": {
						Description: "Different usage type.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"PAYG",
							"COMMITTED",
						}, false),
					},

					"effective_date": {
						Description: "Date when plan was applied.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"user": {
			Description: "User's information.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Description:  "Email of the user used by Dynatrace for contacting them if needed.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsEmailAddress,
					},

					"first_name": {
						Description:  "First name of the user.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Description:  "Last name of the user.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone_number": {
						Description:  "Phone number of the user by Dynatrace for contacting them if needed.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"country": {
						Description:  "Country of the user.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"environment_properties": {
			Description: "Properties of the Dynatrace environment.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"environment_info": {
						Description: "Information about the Dynatrace environment.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"environment_id": {
									Description: "The ID of the Dynatrace environment to be created.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
