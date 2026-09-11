// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -parent-id "storage_account_id"

var storageAccountCustomerManagedKeyResourceName = "azurerm_storage_account_customer_managed_key"

func resourceStorageAccountCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceStorageAccountCustomerManagedKeyRead,
		Update: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceStorageAccountCustomerManagedKeyDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
			},

			"user_assigned_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},

			"federated_identity_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"user_assigned_identity_id"},
			},
		},
	}
}

func resourceStorageAccountCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceStorageAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := storageClient.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceStorageAccountCustomerManagedKeyFlatten(d, id, resp.Model)
}

func resourceStorageAccountCustomerManagedKeyFlatten(d *pluginsdk.ResourceData, id *commonids.StorageAccountId, storageAccount *storageaccounts.StorageAccount) error {
	d.Set("storage_account_id", id.ID())

	enabled := false
	if storageAccount != nil {
		if props := storageAccount.Properties; props != nil {
			if encryption := props.Encryption; encryption != nil && pointer.From(encryption.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
				enabled = true

				if kvProps := encryption.Keyvaultproperties; kvProps != nil {
					keyID, err := keyvault.NewNestedItemID(pointer.From(kvProps.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(kvProps.Keyname), pointer.From(kvProps.Keyversion))
					if err != nil {
						return err
					}
					d.Set("key_vault_key_id", keyID.ID())
				}

				if identityProps := encryption.Identity; identityProps != nil {
					d.Set("user_assigned_identity_id", identityProps.UserAssignedIdentity)
					d.Set("federated_identity_client_id", identityProps.FederatedIdentityClientId)
				}
			}
		}
	}

	if !enabled {
		log.Printf("[DEBUG] Customer Managed Key was not defined for %s - removing from state!", id)
		d.SetId("")
		return nil
	}

	return pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual)
}

func resourceStorageAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	// confirm it still exists prior to trying to update it, else we'll get an error
	storageAccount, err := storageClient.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(storageAccount.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Storage Account I am
	// making it reset the Storage Account to its default state
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
				KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointStorage),
			},
		},
	}

	if _, err = storageClient.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing Customer Managed Key for %s: %+v", *id, err)
	}

	return nil
}
