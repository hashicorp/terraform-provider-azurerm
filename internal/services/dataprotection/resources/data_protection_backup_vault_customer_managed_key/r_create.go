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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupVaultClient
			var cmk DataProtectionBackupVaultCustomerManagedKeyModel

			if err := metadata.Decode(&cmk); err != nil {
				return err
			}

			id, err := backupvaultresources.ParseBackupVaultID(cmk.DataProtectionBackupVaultID)
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if resp.Model.Properties.SecuritySettings != nil && resp.Model.Properties.SecuritySettings.EncryptionSettings != nil {
					return metadata.ResourceRequiresImport(r.ResourceType(), *id)
				}
			}

			payload := resp.Model

			keyId, err := keyvault.ParseNestedItemID(cmk.KeyVaultKeyID, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			payload.Properties.SecuritySettings.EncryptionSettings = &backupvaultresources.EncryptionSettings{
				State: pointer.To(backupvaultresources.EncryptionStateEnabled),
				KeyVaultProperties: &backupvaultresources.CmkKeyVaultProperties{
					KeyUri: pointer.To(keyId.ID()),
				},
			}

			payload.Properties.SecuritySettings.EncryptionSettings.KekIdentity = &backupvaultresources.CmkKekIdentity{
				IdentityType: pointer.To(backupvaultresources.IdentityTypeSystemAssigned),
			}

			if err = client.BackupVaultsCreateOrUpdateCallbackThenPoll(ctx, *id, *payload, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(id)); err != nil {
				return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return nil
		},
	}
}
