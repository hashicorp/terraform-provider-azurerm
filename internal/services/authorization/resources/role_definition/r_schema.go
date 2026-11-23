// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Description:  "A unique UUID/GUID which identifies this role - one will be generated if not specified.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"name": {
			Description:  "The name of the Role Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Description:  "The scope at which the Role Definition applies to, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, `/providers/Microsoft.Management/managementGroups/0b1f6471-1bf0-4dda-aec3-111122223333`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`. It is recommended to use the first entry of the `assignable_scopes`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringStartsWithOneOf("/subscriptions/", "/providers/Microsoft.Management/managementGroups/"),
		},

		"description": {
			Description:  "A description of the Role Definition.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// lintignore:XS003
		"permissions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Description: "One or more Allowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager Resource Provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"not_actions": {
						Description: "One or more Disallowed Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager Resource Provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"data_actions": {
						Description: "One or more Allowed Data Actions, such as `*`, `Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read`. See ['Azure Resource Manager Resource Provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
					"not_data_actions": {
						Description: "One or more Disallowed Data Actions, such as `*`, `Microsoft.Resources/subscriptions/resourceGroups/read`. See ['Azure Resource Manager Resource Provider operations'](https://docs.microsoft.com/azure/role-based-access-control/resource-provider-operations) for details.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"assignable_scopes": {
			Description: "One or more assignable scopes for this Role Definition, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333`, `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`, `/providers/Microsoft.Management/managementGroups/0b1f6471-1bf0-4dda-aec3-111122223333`, or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup/providers/Microsoft.Compute/virtualMachines/myVM`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateScopeID,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_resource_id": {
			Description: "The Azure Resource Manager ID for the resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
