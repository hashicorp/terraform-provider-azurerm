// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

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

func resourceKeyVaultSecretDelete(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
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
			log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseURL, err)
	}

	shouldPurge := meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedSecretsOnDestroy
	if shouldPurge && kv.Model != nil && pointer.From(kv.Model.Properties.EnablePurgeProtection) {
		log.Printf("[DEBUG] cannot purge secret %q because %s has purge protection enabled", id.Name, *keyVaultId)
		shouldPurge = false
	}

	description := fmt.Sprintf("Secret %q (Key Vault %q)", id.Name, id.KeyVaultBaseURL)
	deleter := deleteAndPurgeSecret{
		client:      client,
		keyVaultUri: id.KeyVaultBaseURL,
		name:        id.Name,
	}
	if err := helpers.DeleteAndOptionallyPurge(ctx, description, shouldPurge, deleter); err != nil {
		return err
	}

	return nil
}
