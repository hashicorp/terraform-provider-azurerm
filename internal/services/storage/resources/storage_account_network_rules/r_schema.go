// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageAccountNetworkRulesSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// lintignore: S013
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"bypass": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			// Note: O+C because this defaults to storageaccounts.BypassAzureServices in the API, but schema does not support defaults for lists/sets.
			Computed: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForBypass(), false),
			},
		},

		"ip_rules": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.StorageAccountIpRule,
			},
		},

		"virtual_network_subnet_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"default_action": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForDefaultAction(), false),
		},

		"private_link_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"endpoint_resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"endpoint_tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
	}
}
