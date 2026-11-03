// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account_sas_token_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func keyVaultManagedStorageAccountSasTokenDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"managed_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeStorage),
		},

		"sas_template_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sas_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"service",
				"account",
			}, false),
		},

		"secret_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"validity_period": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"tags": commonschema.TagsForceNew(),
	}
}
