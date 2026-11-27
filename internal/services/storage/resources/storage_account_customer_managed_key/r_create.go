// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	existing, err := storageClient.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			if existing.Model.Properties.Encryption != nil && pointer.From(existing.Model.Properties.Encryption.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
				return tf.ImportAsExistsError(storageAccountCustomerManagedKeyResourceName, id.ID())
			}
		}
	}

	keyID, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			Encryption: &storageaccounts.Encryption{
				Services: &storageaccounts.EncryptionServices{
					Blob: &storageaccounts.EncryptionService{
						Enabled: pointer.To(true),
					},
					File: &storageaccounts.EncryptionService{
						Enabled: pointer.To(true),
					},
				},
				Identity: &storageaccounts.EncryptionIdentity{
					UserAssignedIdentity: pointer.To(d.Get("user_assigned_identity_id").(string)),
				},
				KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointKeyvault),
				Keyvaultproperties: &storageaccounts.KeyVaultProperties{
					Keyname:     pointer.To(keyID.Name),
					Keyversion:  pointer.To(keyID.Version),
					Keyvaulturi: pointer.To(keyID.KeyVaultBaseURL),
				},
			},
		},
	}

	if fID := d.Get("federated_identity_client_id").(string); fID != "" {
		payload.Properties.Encryption.Identity.FederatedIdentityClientId = pointer.To(fID)
	}

	if _, err = storageClient.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return resourceStorageAccountCustomerManagedKeyRead(d, meta)
}
