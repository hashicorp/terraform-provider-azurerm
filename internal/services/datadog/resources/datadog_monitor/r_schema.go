// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func datadogMonitorSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatadogMonitorsName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"datadog_organization": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"api_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						ForceNew:  true,
						Sensitive: true,
					},

					"application_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						ForceNew:  true,
						Sensitive: true,
					},

					"enterprise_app_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"linking_auth_code": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
					},

					"linking_client_id": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						ForceNew:  true,
						Sensitive: true,
					},

					"redirect_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"sku_name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: SkuNameDiffSuppress,
		},

		"user": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 50),
					},

					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.DatadogMonitorsEmailAddress,
					},

					"phone_number": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(0, 40),
					},
				},
			},
		},

		"monitoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"marketplace_subscription_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
