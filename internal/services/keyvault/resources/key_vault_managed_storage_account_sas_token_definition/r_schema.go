// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this SAS Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemName,
		},

		"managed_storage_account_id": {
			Description:  "The ID of the Managed Storage Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeStorage),
		},

		"sas_template_uri": {
			Description:  "The SAS definition token template signed with an arbitrary key. Tokens created according to the SAS definition will have the same properties as the template, but regenerated with a new validity period.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sas_type": {
			Description: "The type of SAS token the SAS definition will create.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"service",
				"account",
			}, false),
		},

		"secret_id": {
			Description: "The ID of the Secret that is created by Managed Storage Account SAS Definition.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"validity_period": {
			Description:  "Validity period of SAS token. Value needs to be in [ISO 8601 duration format](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"tags": commonschema.TagsForceNew(),
	}
}
