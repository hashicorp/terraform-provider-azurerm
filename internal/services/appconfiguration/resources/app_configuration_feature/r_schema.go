// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (k Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Description: "Specifies the id of the App Configuration.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			// User-specified segments are lowercased in the API response
			// tracked in https://github.com/Azure/azure-rest-api-specs/issues/24337
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     configurationstores.ValidateConfigurationStoreID,
		},
		"description": {
			Description: "The description of the App Configuration Feature.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
		"enabled": {
			Description: "The status of the App Configuration Feature. By default, this is set to false.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
		"key": {
			Description: "The key of the App Configuration Feature. The value for `name` will be used if this is unspecified.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C We generate a value for this if it's omitted so this should be kept
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppConfigurationFeatureKey,
		},
		"name": {
			Description:  "The name of the App Configuration Feature.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppConfigurationFeatureName,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			// NOTE: O+C The value of this is updated anytime the resource changes
			Optional: true,
		},
		"label": {
			Description: "The label of the App Configuration Feature.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"locked": {
			Description: "Should this App Configuration Feature be Locked to prevent changes?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"percentage_filter_value": {
			Description:  "A number representing the value of the percentage required to enable this feature.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"targeting_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"default_rollout_percentage": {
						Description:  "A number representing the percentage of the entire user base.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Description: "The name of the group.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"rollout_percentage": {
									Description:  "Rollout percentage of the group.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 100),
								},
							},
						},
					},
					"users": {
						Description: "A list of users to target for this feature.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
		"timewindow_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"start": {
						Description:  "The earliest timestamp the feature is enabled. The timestamp must be in RFC3339 format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"end": {
						Description:  "The latest timestamp the feature is enabled. The timestamp must be in RFC3339 format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
		"custom_filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description:  "The name of the custom filter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
		"tags": commonschema.Tags(),
	}
}

func (k Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
