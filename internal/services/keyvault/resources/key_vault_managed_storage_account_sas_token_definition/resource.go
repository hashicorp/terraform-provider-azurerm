// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account_sas_token_definition

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate,
		Read:   resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead,
		Update: resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate,
		Delete: resourceKeyVaultManagedStorageAccountSasTokenDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SasDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultManagedStorageAccountSasTokenDefinitionSchema(),
	}
}
