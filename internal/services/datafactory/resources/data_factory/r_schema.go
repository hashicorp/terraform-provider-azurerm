// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactorySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"github_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"vsts_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Description:  "The GitHub account name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"branch_name": {
						Description:  "The branch of the repository to get code from.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"git_url": {
						Description: "The GitHub Enterprise host name. For example: `https://github.<mydomain>.com`. Use `https://github.com` for open source repositories.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"last_commit_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"repository_name": {
						Description:  "The name of the git repository.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"root_folder": {
						Description:  "The root folder within the repository. Set to `/` for the top level.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"publishing_enabled": {
						Description: "Is automated publishing enabled?",
						Type:        pluginsdk.TypeBool,
						Default:     true,
						Optional:    true,
					},
				},
			},
		},

		"vsts_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Description:  "The VSTS account name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"branch_name": {
						Description:  "The branch of the repository to get code from.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"project_name": {
						Description:  "The name of the VSTS project.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"repository_name": {
						Description:  "The name of the git repository.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"root_folder": {
						Description:  "The root folder within the repository. Set to `/` for the top level.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Description:  "The Tenant ID associated with the VSTS account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"publishing_enabled": {
						Description: "Is automated publishing enabled?",
						Type:        pluginsdk.TypeBool,
						Default:     true,
						Optional:    true,
					},
				},
			},
		},

		"global_parameter": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The global parameter name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"type": {
						Description:  "The global parameter type.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(factories.PossibleValuesForGlobalParameterType(), false),
					},

					"value": {
						Description:  "The global parameter value.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"managed_virtual_network_enabled": {
			Description: "Is Managed Virtual Network enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"public_network_enabled": {
			Description: "Is the Data Factory visible to the public network?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"purview_id": {
			Description:  "The ID of the Purview Account resource associated with the Data Factory.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: account.ValidateAccountID,
		},

		"customer_managed_key_id": {
			Description:  "The Azure Key Vault Key ID to be used as the Customer Managed Key (CMK) for double encryption.",
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"customer_managed_key_identity_id": {
			Description:  "The ID of the User Assigned Identity associated with the Customer Managed Key. Can be left empty once factory is created to use the System Assigned Identity. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/enable-customer-managed-key) for more information.",
			Type:         pluginsdk.TypeString,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"tags": commonschema.Tags(),
	}
}
