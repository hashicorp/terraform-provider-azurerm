// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_channels_registration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func botChannelsRegistrationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Description:  "The SKU of the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(botservice.PossibleSkuNameValues(), false),
		},

		"microsoft_app_id": {
			Description:  "The Microsoft Application ID for the Bot Channels Registration.",
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

		"cmk_key_vault_url": {
			Description:  "The CMK Key Vault Key URL to encrypt the Bot Channels Registration with the Customer Managed Encryption Key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"description": {
			Description:  "The description of the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 512),
		},

		"display_name": {
			Description:  "The name of the Bot Channels Registration will be displayed as.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint": {
			Description:  "The Bot Channels Registration endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_key": {
			Description:  "The Application Insights Key to associate with the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"developer_app_insights_api_key": {
			Description:  "The Application Insights API Key to associate with the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_application_id": {
			Description:  "The Application Insights Application ID to associate with the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"icon_url": {
			Description:  "The icon URL to visually identify the Bot Channels Registration.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png",
			ValidateFunc: validate.BotChannelRegistrationIconUrl,
		},

		"streaming_endpoint_enabled": {
			Description: "Whether the streaming endpoint is enabled for the Bot Channels Registration.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"public_network_access_enabled": {
			Description: "Whether the Bot Channels Registration is in an isolated network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tags": commonschema.Tags(),
	}
}
