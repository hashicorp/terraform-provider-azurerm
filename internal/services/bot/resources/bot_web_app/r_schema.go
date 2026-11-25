// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_web_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func botWebAppSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Description:  "The SKU of the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(botservice.PossibleSkuNameValues(), false),
		},

		"microsoft_app_id": {
			Description:  "The Microsoft Application ID for the Web App Bot.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"microsoft_app_type": {
			Description:  "The Microsoft Application Type for the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(botservice.PossibleMsaAppTypeValues(), false),
		},

		"microsoft_app_tenant_id": {
			Description:  "The Microsoft Application Tenant ID for the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"microsoft_app_user_assigned_identity_id": {
			Description:  "The ID of Microsoft Application User Assigned Identity for the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"display_name": {
			Description:  "The name of the Web App Bot will be displayed as.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint": {
			Description:  "The Web App Bot endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_key": {
			Description:  "The Application Insights Key to associate with the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"developer_app_insights_api_key": {
			Description:  "The Application Insights API Key to associate with the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_application_id": {
			Description:  "The Application Insights Application ID to associate with the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"luis_app_ids": {
			Description: "A list of LUIS App IDs to associate with the Web App Bot.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"luis_key": {
			Description:  "The LUIS key to associate with the Web App Bot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}
