// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageEncryptionScopeRead(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.EncryptionScopeName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.State != nil && *props.State == encryptionscopes.EncryptionScopeStateDisabled {
				log.Printf("[INFO] %s was not configured - removing from state", id)
				d.SetId("")
				return nil
			}

			d.Set("infrastructure_encryption_required", props.RequireInfrastructureEncryption)
			d.Set("source", pointer.FromEnum(props.Source))

			keyVaultKeyUri := ""
			if props.KeyVaultProperties != nil {
				keyVaultKeyUri = pointer.From(props.KeyVaultProperties.KeyUri)
			}
			d.Set("key_vault_key_id", keyVaultKeyUri)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
