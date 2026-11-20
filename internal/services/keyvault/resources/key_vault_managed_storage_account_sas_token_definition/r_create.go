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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedstorage"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
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

	client := keyVaultsClient.DataPlaneKeyVaultClient.Storage.Clone(*keyVaultBaseUri)
	sasId := storage.NewSaID(*keyVaultBaseUri, storageAccount.Name, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetSasDefinition(ctx, sasId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
				}
			}

			if model := existing.Model; model != nil && pointer.From(model.Id) != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account_sas_token_definition", *model.Id)
			}
		}
	}

	t := d.Get("tags").(map[string]any)
	parameters := storage.SasDefinitionCreateParameters{
		TemplateUri:    d.Get("sas_template_uri").(string),
		SasType:        storage.SasTokenType(d.Get("sas_type").(string)),
		ValidityPeriod: d.Get("validity_period").(string),
		Attributes: &storage.SasDefinitionAttributes{
			Enabled: pointer.To(true),
		},
		Tags: tags.Expand(t),
	}

	if resp, err := client.SetSasDefinition(ctx, sasId, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && response.WasConflict(resp.HttpResponse) {
			deletedStorageClient := keyVaultsClient.DataPlaneKeyVaultClient.DeletedStorage.Clone(*keyVaultBaseUri)
			recoveredStorageAccount, err := deletedStorageClient.RecoverDeletedSasDefinition(ctx, deletedstorage.NewDeletedstorageSaID(*keyVaultBaseUri, storageAccount.Name, name))
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account Sas Definition %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
			// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
			if model := recoveredStorageAccount.Model; model != nil && model.Id != nil {
				secret := model.Id
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

				if _, err := client.SetSasDefinition(ctx, sasId, parameters); err != nil {
					return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creation of Managed Storage Account SAS Definition %q (Storage Account %q, Key Vault %q): %+v", name, storageAccount.Name, *keyVaultId, err)
		}
	}

	read, err := client.GetSasDefinition(ctx, sasId)
	if err != nil {
		return err
	}

	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("cannot read Managed Storage Account Sas Definition  %q (Storage Account %q, Key Vault %q)", name, storageAccount.Name, *keyVaultId)
	}

	parsedSasId, err := parse.SasDefinitionID(*read.Model.Id)
	if err != nil {
		return err
	}

	d.SetId(parsedSasId.ID())

	return resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead(d, meta)
}
