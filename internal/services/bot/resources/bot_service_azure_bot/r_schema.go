// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bot_service_azure_bot

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Description:  "The SKU of the Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(botservice.PossibleSkuNameValues(), false),
		},

		"microsoft_app_id": {
			Description:  "The Microsoft Application ID for the Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"display_name": {
			Description:  "The name that the Azure Bot Service will be displayed as.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint": {
			Description:  "The Azure Bot Service endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_key": {
			Description:  "The Application Insight Key to associate with this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"developer_app_insights_api_key": {
			Description:  "The Application Insights API Key to associate with this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"developer_app_insights_application_id": {
			Description:  "The resource ID of the Application Insights instance to associate with this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cmk_key_vault_key_url": {
			Description:  "The CMK Key Vault Key URL that will be used to encrypt the Bot with the Customer Managed Encryption Key.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"microsoft_app_msi_id": {
			Description:  "The ID of the Microsoft App Managed Identity for this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"microsoft_app_tenant_id": {
			Description:  "The Tenant ID of the Microsoft App for this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"microsoft_app_type": {
			Description:  "The Microsoft App Type for this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(botservice.PossibleMsaAppTypeValues(), false),
		},

		"local_authentication_enabled": {
			Description: "Is local authentication enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"luis_app_ids": {
			Description: "A list of LUIS App IDs to associate with this Azure Bot Service.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"luis_key": {
			Description:  "The LUIS key to associate with this Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"streaming_endpoint_enabled": {
			Description: "Is the streaming endpoint enabled for this Azure Bot Service?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"icon_url": {
			Description:  "The Icon URL of the Azure Bot Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "https://docs.botframework.com/static/devportal/client/images/bot-framework-default.png",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
