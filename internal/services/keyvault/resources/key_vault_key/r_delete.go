// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseURL)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseURL, err)
	}
	if keyVaultIdRaw == nil {
		return fmt.Errorf("unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseURL)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	kv, err := keyVaultsClient.VaultsClient.Get(ctx, *keyVaultId)
	if err != nil {
		if response.WasNotFound(kv.HttpResponse) {
			log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving key vault %q properties: %+v", *keyVaultId, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedKeysOnDestroy
	if shouldPurge && kv.Model != nil && pointer.From(kv.Model.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge key %q because vault %q has purge protection enabled", id.Name, keyVaultId.String())
		shouldPurge = false
	}

	description := fmt.Sprintf("Key %q (Key Vault %q)", id.Name, id.KeyVaultBaseURL)
	deleter := deleteAndPurgeKey{
		keysClient:        keyVaultsClient.DataPlaneKeyVaultClient.Keys.Clone(id.KeyVaultBaseURL),
		deletedKeysClient: keyVaultsClient.DataPlaneKeyVaultClient.DeletedKeys.Clone(id.KeyVaultBaseURL),
		keyVaultUri:       id.KeyVaultBaseURL,
		name:              id.Name,
	}
	if err := helpers.DeleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}
