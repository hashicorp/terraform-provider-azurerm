// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupVaultCustomerManagedKeyResource struct{}

type DataProtectionBackupVaultCustomerManagedKeyModel struct {
	DataProtectionBackupVaultID string `tfschema:"data_protection_backup_vault_id"`
	KeyVaultKeyID               string `tfschema:"key_vault_key_id"`
}

var _ sdk.ResourceWithUpdate = DataProtectionBackupVaultCustomerManagedKeyResource{}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) ModelObject() interface{} {
	return &DataProtectionBackupVaultCustomerManagedKeyResource{}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) ResourceType() string {
	return "azurerm_data_protection_backup_vault_customer_managed_key"
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaults.ValidateBackupVaultID
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_protection_backup_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupvaults.ValidateBackupVaultID,
		},

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
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

			id, err := backupvaults.ParseBackupVaultID(cmk.DataProtectionBackupVaultID)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			resp, err := client.Get(ctx, *id)
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

			keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(cmk.KeyVaultKeyID)
			if err != nil {
				return err
			}

			payload.Properties.SecuritySettings.EncryptionSettings = &backupvaults.EncryptionSettings{
				State: pointer.To(backupvaults.EncryptionStateEnabled),
				KeyVaultProperties: &backupvaults.CmkKeyVaultProperties{
					KeyUri: pointer.To(keyId.ID()),
				},
			}

			payload.Properties.SecuritySettings.EncryptionSettings.KekIdentity = &backupvaults.CmkKekIdentity{
				IdentityType: pointer.To(backupvaults.IdentityTypeSystemAssigned),
			}

			err = client.CreateOrUpdateThenPoll(ctx, *id, *payload, backupvaults.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataProtectionBackupVaultCustomerManagedKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupVaultClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
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

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			resp, err := client.Get(ctx, *id)
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
				keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(cmk.KeyVaultKeyID)
				if err != nil {
					return err
				}
				payload.Properties.SecuritySettings.EncryptionSettings.KeyVaultProperties = &backupvaults.CmkKeyVaultProperties{
					KeyUri: pointer.To(keyId.ID()),
				}
			}

			err = client.CreateOrUpdateThenPoll(ctx, *id, *payload, backupvaults.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
