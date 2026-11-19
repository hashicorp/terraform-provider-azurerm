// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "Specifies the name of the Automation Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Description:  "The SKU of the account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(automationaccount.PossibleValuesForSkuNameEnum(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"user_assigned_identity_id": {
						Description:  "The User Assigned Managed Identity ID to be used for accessing the Customer Managed Key for encryption.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},

					"key_vault_key_id": {
						Description:  "The ID of the Key Vault Key which should be used to Encrypt the data in this Automation Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},
				},
			},
		},

		"local_authentication_enabled": {
			Description: "Whether requests using non-AAD authentication are blocked.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),

		"dsc_server_endpoint": {
			Description: "The DSC Server Endpoint associated with this Automation Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"dsc_primary_access_key": {
			Description: "The Primary Access Key for the DSC Endpoint associated with this Automation Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"dsc_secondary_access_key": {
			Description: "The Secondary Access Key for the DSC Endpoint associated with this Automation Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for the automation account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"private_endpoint_connection": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"hybrid_service_url": {
			Description: "The URL of automation hybrid service which is used for hybrid worker on-boarding With this Automation Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
