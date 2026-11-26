// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the user that will be associated with the Datadog Monitor.",
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
						Description: "API key associated to the Datadog organization.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						Sensitive:   true,
					},

					"application_key": {
						Description: "Application key associated to the Datadog organization.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						Sensitive:   true,
					},

					"enterprise_app_id": {
						Description: "The ID of the enterprise_app.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},

					"linking_auth_code": {
						Description: "The auth code used to linking to an existing Datadog organization.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Sensitive:   true,
					},

					"linking_client_id": {
						Description: "The ID of the linking_client.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
						Sensitive:   true,
					},

					"redirect_uri": {
						Description: "The redirect URI for linking.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
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
			Description:      "The name which should be used for this SKU.",
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
						Description:  "The name which should be used for this user_info.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 50),
					},

					"email": {
						Description:  "Email of the user used by Datadog for contacting them if needed.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.DatadogMonitorsEmailAddress,
					},

					"phone_number": {
						Description:  "Phone number of the user used by Datadog for contacting them if needed.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(0, 40),
					},
				},
			},
		},

		"monitoring_enabled": {
			Description: "Whether monitoring is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"marketplace_subscription_status": {
			Description: "Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
