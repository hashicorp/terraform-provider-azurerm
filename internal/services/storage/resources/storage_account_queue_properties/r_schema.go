// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_queue_properties

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"cors_rule": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     5,
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"exposed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"allowed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"allowed_methods": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"DELETE",
								"GET",
								"HEAD",
								"MERGE",
								"POST",
								"OPTIONS",
								"PUT",
							}, false),
						},
					},
					"max_age_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 2000000000),
					},
				},
			},
		},

		"hour_metrics": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			MaxItems:     1,
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"include_apis": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
		},
		"logging": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			MaxItems:     1,
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"delete": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"read": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"write": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
		},
		"minute_metrics": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			MaxItems:     1,
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"include_apis": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
		},
	}
}

func (s Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
