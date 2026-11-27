// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageEncryptionScopeUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	if d.Get("source").(string) == string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault))
		}
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	payload := existing.Model
	payload.Properties.State = pointer.To(encryptionscopes.EncryptionScopeStateEnabled)
	if d.HasChange("key_vault_key_id") {
		payload.Properties.KeyVaultProperties = &encryptionscopes.EncryptionScopeKeyVaultProperties{
			KeyUri: pointer.To(d.Get("key_vault_key_id").(string)),
		}
	}
	if d.HasChange("source") {
		payload.Properties.Source = pointer.ToEnum[encryptionscopes.EncryptionScopeSource](d.Get("source").(string))
	}

	if _, err := client.Patch(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStorageEncryptionScopeRead(d, meta)
}
