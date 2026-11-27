// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageEncryptionScopeCreate(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := encryptionscopes.NewEncryptionScopeID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
			}
		}
		if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.State != nil {
			if *existing.Model.Properties.State == encryptionscopes.EncryptionScopeStateEnabled {
				return tf.ImportAsExistsError("azurerm_storage_encryption_scope", id.ID())
			}
		}
	}

	if d.Get("source").(string) == string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault))
		}
	}

	payload := encryptionscopes.EncryptionScope{
		Properties: &encryptionscopes.EncryptionScopeProperties{
			Source: pointer.ToEnum[encryptionscopes.EncryptionScopeSource](d.Get("source").(string)),
			State:  pointer.To(encryptionscopes.EncryptionScopeStateEnabled),
			KeyVaultProperties: &encryptionscopes.EncryptionScopeKeyVaultProperties{
				KeyUri: pointer.To(d.Get("key_vault_key_id").(string)),
			},
		},
	}
	if v, ok := d.GetOk("infrastructure_encryption_required"); ok {
		payload.Properties.RequireInfrastructureEncryption = pointer.To(v.(bool))
	}

	if _, err := client.Put(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceStorageEncryptionScopeRead(d, meta)
}
