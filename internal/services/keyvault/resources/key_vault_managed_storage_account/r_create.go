// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv74 "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultManagedStorageAccountCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
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

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetStorageAccount(ctx, *keyVaultBaseUrl, name)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultId, err)
				}
			}

			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_key_vault_managed_storage_account", *existing.ID)
			}
		}
	}

	t := d.Get("tags").(map[string]any)

	parameters := kv74.StorageAccountCreateParameters{
		ResourceID:         pointer.To(d.Get("storage_account_id").(string)),
		ActiveKeyName:      pointer.To(d.Get("storage_account_key").(string)),
		AutoRegenerateKey:  pointer.To(d.Get("regenerate_key_automatically").(bool)),
		RegenerationPeriod: pointer.To(d.Get("regeneration_period").(string)),
		Tags:               tags.Expand(t),
	}

	if resp, err := client.SetStorageAccount(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
		// In the case that the Storage Account already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeyVaults && response.WasConflict(resp.Response.Response) {
			recoveredStorageAccount, err := client.RecoverDeletedStorageAccount(ctx, *keyVaultBaseUrl, name)
			if err != nil {
				return fmt.Errorf("recovery of Managed Storage Account %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
			}
			log.Printf("[DEBUG] Recovering Managed Storage Account %q (Key Vault %q)", name, *keyVaultBaseUrl)
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
					return fmt.Errorf("waiting for Managed Storage Account %q (Key Vault %q) to become available after recovery: %s", name, *keyVaultId, err)
				}
				log.Printf("[DEBUG] Managed Storage Account %q recovered with ID: %q", name, *recoveredStorageAccount.ID)

				if _, err := client.SetStorageAccount(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
					return fmt.Errorf("creating Managed Storage Account %s (Key Vault %q): %+v", name, *keyVaultId, err)
				}
			}
		} else {
			return fmt.Errorf("creating Managed Storage Account %s: %+v", name, err)
		}
	}

	read, err := client.GetStorageAccount(ctx, *keyVaultBaseUrl, name)
	if err != nil {
		return fmt.Errorf("cannot read Managed Storage Account %q (Key Vault %q): %+v", name, *keyVaultId, err)
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read Managed Storage Account ID %q (Key Vault %q)", name, *keyVaultId)
	}

	storageId, err := keyvault.ParseNestedItemID(*read.ID, keyvault.VersionTypeAny, keyvault.NestedItemTypeStorage)
	if err != nil {
		return err
	}

	d.SetId(storageId.ID())

	return resourceKeyVaultManagedStorageAccountRead(d, meta)
}
