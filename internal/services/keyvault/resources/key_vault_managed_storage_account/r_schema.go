// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultManagedStorageAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Key Vault Managed Storage Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"key_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KeyVaultId{}, "Key Vault"),

		"storage_account_key": {
			Description: "Which Storage Account access key that is managed by Key Vault.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"key1",
				"key2",
			}, false),
		},

		"storage_account_id": {
			Description:  "The ID of the Storage Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"regenerate_key_automatically": {
			Description:  "Whether Storage Account access key should be regenerated periodically.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			RequiredWith: []string{"regeneration_period"},
		},

		"regeneration_period": {
			Description:  "How often Storage Account access key should be regenerated. Value needs to be in [ISO 8601 duration format](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601Duration,
			RequiredWith: []string{"regenerate_key_automatically"},
		},

		"tags": commonschema.TagsForceNew(),
	}
}
