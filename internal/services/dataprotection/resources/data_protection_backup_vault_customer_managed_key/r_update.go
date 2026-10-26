// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_vault_customer_managed_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupVaultClient
			var cmk DataProtectionBackupVaultCustomerManagedKeyModel

			if err := metadata.Decode(&cmk); err != nil {
				return err
			}

			id, err := backupvaultresources.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			resp, err := client.BackupVaultsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", *id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", *id)
			}

			if resp.Model.Properties.SecuritySettings == nil || resp.Model.Properties.SecuritySettings.EncryptionSettings == nil {
				return fmt.Errorf("retrieving %s: Customer Managed Key was not found", *id)
			}

			payload := resp.Model

			if metadata.ResourceData.HasChange("key_vault_key_id") {
				keyId, err := keyvault.ParseNestedItemID(cmk.KeyVaultKeyID, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}
				payload.Properties.SecuritySettings.EncryptionSettings.KeyVaultProperties = &backupvaultresources.CmkKeyVaultProperties{
					KeyUri: pointer.To(keyId.ID()),
				}
			}

			if err = client.BackupVaultsCreateOrUpdateThenPoll(ctx, *id, *payload, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
