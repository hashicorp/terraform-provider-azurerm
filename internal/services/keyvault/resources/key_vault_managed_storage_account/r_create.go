// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultManagedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Managed Storage Account Key Vault %q in %s: %+v", name, *keyVaultId, err)
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Storage.Clone(*keyVaultBaseUrl)
	storageId := storage.NewStorageID(*keyVaultBaseUrl, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetStorageAccount(ctx, storageId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultId, err)
				}
			}

			if model := existing.Model; model != nil && pointer.From(model.Id) != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account", *model.Id)
			}
		}
	}

	t := d.Get("tags").(map[string]any)

	parameters := storage.StorageAccountCreateParameters{
		ResourceId:         d.Get("storage_account_id").(string),
		ActiveKeyName:      d.Get("storage_account_key").(string),
		AutoRegenerateKey:  d.Get("regenerate_key_automatically").(bool),
		RegenerationPeriod: pointer.To(d.Get("regeneration_period").(string)),
		Tags:               tags.Expand(t),
	}

	if resp, err := client.SetStorageAccount(ctx, storageId, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && response.WasConflict(resp.HttpResponse) {
			deletedStorageClient := keyVaultsClient.DataPlaneKeyVaultClient.DeletedStorage.Clone(*keyVaultBaseUrl)
			recoveredStorageAccount, err := deletedStorageClient.RecoverDeletedStorageAccount(ctx, deletedstorage.NewDeletedstorageID(*keyVaultBaseUrl, name))
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account %q (Key Vault %q)", name, *keyVaultBaseUrl)
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
					return fmt.Errorf("waiting for Managed Storage Account %q (Key Vault %q) to become available after recovery: %s", name, *keyVaultId, err)
				}
				log.Printf("[DEBUG] Managed Storage Account %q recovered with ID: %q", name, *secret)

				if _, err := client.SetStorageAccount(ctx, storageId, parameters); err != nil {
					return fmt.Errorf("creating Managed Storage Account %s (Key Vault %q): %+v", name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creating Managed Storage Account %s: %+v", name, err)
		}
	}

	read, err := client.GetStorageAccount(ctx, storageId)
	if err != nil {
		return fmt.Errorf("cannot read Managed Storage Account %q (Key Vault %q): %+v", name, *keyVaultId, err)
	}
	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("cannot read Managed Storage Account ID %q (Key Vault %q)", name, *keyVaultId)
	}

	nestedId, err := keyvault.ParseNestedItemID(*read.Model.Id, keyvault.VersionTypeAny, keyvault.NestedItemTypeStorage)
	if err != nil {
		return err
	}

	d.SetId(nestedId.ID())

	return resourceKeyVaultManagedStorageAccountRead(d, meta)
}
