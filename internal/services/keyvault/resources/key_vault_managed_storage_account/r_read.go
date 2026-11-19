// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultManagedStorageAccountRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeAny, keyvault.NestedItemTypeStorage)
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseURL)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", id.KeyVaultBaseURL, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if Key Vault %q for Managed Storage Account %q exists: %v", *keyVaultId, id.Name, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account %q (Key Vault %q) was not found - removing from state", id.Name, *keyVaultId)
		d.SetId("")
		return nil
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Storage.Clone(id.KeyVaultBaseURL)
	resp, err := client.GetStorageAccount(ctx, storage.NewStorageID(id.KeyVaultBaseURL, id.Name))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Managed Storage Account %q (Key Vault %q) was not found - removing from state", id.Name, *keyVaultId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Managed Storage Account %q (Key Vault %q): %+v", id.Name, *keyVaultId, err)
	}

	d.Set("name", id.Name)
	d.Set("key_vault_id", keyVaultId.ID())

	if model := resp.Model; model != nil {
		d.Set("storage_account_id", model.ResourceId)
		d.Set("storage_account_key", model.ActiveKeyName)
		d.Set("regenerate_key_automatically", model.AutoRegenerateKey)
		d.Set("regeneration_period", model.RegenerationPeriod)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
