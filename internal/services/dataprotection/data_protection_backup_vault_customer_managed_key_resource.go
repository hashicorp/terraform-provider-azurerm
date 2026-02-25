// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_vault_customer_managed_key -service-package-name dataprotection -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,name:data_protection_backup_vault_id" -test-name complete

type DataProtectionBackupVaultCustomerManagedKeyResource struct{}

type DataProtectionBackupVaultCustomerManagedKeyModel struct {
	DataProtectionBackupVaultID string `tfschema:"data_protection_backup_vault_id"`
	KeyVaultKeyID               string `tfschema:"key_vault_key_id"`
}

var (
	_ sdk.ResourceWithUpdate   = DataProtectionBackupVaultCustomerManagedKeyResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupVaultCustomerManagedKeyResource{}
)

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Identity() resourceids.ResourceId {
	return &backupvaultresources.BackupVaultId{}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) ModelObject() interface{} {
	return &DataProtectionBackupVaultCustomerManagedKeyResource{}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) ResourceType() string {
	return "azurerm_data_protection_backup_vault_customer_managed_key"
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaultresources.ValidateBackupVaultID
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_protection_backup_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupvaultresources.ValidateBackupVaultID,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},
	}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Create() sdk.ResourceFunc {
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

			if resp.Model.Properties.SecuritySettings != nil && resp.Model.Properties.SecuritySettings.EncryptionSettings != nil {
				return metadata.ResourceRequiresImport(r.ResourceType(), *id)
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

			err = client.BackupVaultsCreateOrUpdateThenPoll(ctx, *id, *payload, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions())
			if err != nil {
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

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupVaultClient

			id, err := backupvaultresources.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.BackupVaultsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state DataProtectionBackupVaultCustomerManagedKeyModel
			state.DataProtectionBackupVaultID = id.ID()

			if model := existing.Model; model != nil {
				props := model.Properties
				if props.SecuritySettings != nil && props.SecuritySettings.EncryptionSettings != nil {
					if props.SecuritySettings.EncryptionSettings.KeyVaultProperties != nil {
						state.KeyVaultKeyID = pointer.From(props.SecuritySettings.EncryptionSettings.KeyVaultProperties.KeyUri)
					}
				}
			}
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			log.Printf(`[INFO] Customer Managed Keys cannot be removed from Data Protection Backup Vaults once added. To remove the Customer Managed Key, delete and recreate the parent Data Protection Backup Vault`)
			return nil
		},
	}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Update() sdk.ResourceFunc {
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

			err = client.BackupVaultsCreateOrUpdateThenPoll(ctx, *id, *payload, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
