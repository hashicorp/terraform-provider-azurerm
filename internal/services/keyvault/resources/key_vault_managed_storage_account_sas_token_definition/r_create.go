// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account_sas_token_definition

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv74 "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageAccount, err := keyvault.ParseNestedItemID(d.Get("managed_storage_account_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeStorage)
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, storageAccount.KeyVaultBaseURL)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", storageAccount.KeyVaultBaseURL, err)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Managed Storage Account Key Vault %s: %+v", *keyVaultId, err)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
				}
			}

			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account_sas_token_definition", *existing.ID)
			}
		}
	}

	t := d.Get("tags").(map[string]any)
	parameters := kv74.SasDefinitionCreateParameters{
		TemplateURI:    pointer.To(d.Get("sas_template_uri").(string)),
		SasType:        kv74.SasTokenType(d.Get("sas_type").(string)),
		ValidityPeriod: pointer.To(d.Get("validity_period").(string)),
		SasDefinitionAttributes: &kv74.SasDefinitionAttributes{
			Enabled: pointer.To(true),
		},
		Tags: tags.Expand(t),
	}

	if resp, err := client.SetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && response.WasConflict(resp.Response.Response) {
			recoveredStorageAccount, err := client.RecoverDeletedSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
			// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
			if secret := recoveredStorageAccount.ID; secret != nil {
				stateConf := &pluginsdk.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   helpers.KeyVaultChildItemRefreshFunc(ctx, *secret),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for Key Vault Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q) to become available: %s", name, storageAccount.Name, *keyVaultId, err)
				}
				log.Printf("[DEBUG] Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q) recovered", name, storageAccount.Name, *keyVaultId)

				if _, err := client.SetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name, parameters); err != nil {
					return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
		}
	}

	read, err := client.GetSasDefinition(ctx, *keyVaultBaseUri, storageAccount.Name, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Managed Storage Account Sas Definition  %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
	}

	sasId, err := parse.SasDefinitionID(*read.ID)
	if err != nil {
		return err
	}

	d.SetId(sasId.ID())

	return resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead(d, meta)
}
