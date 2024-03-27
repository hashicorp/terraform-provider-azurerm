// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageAccountCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Read:   resourceStorageAccountCustomerManagedKeyRead,
		Update: resourceStorageAccountCustomerManagedKeyCreateUpdate,
		Delete: resourceStorageAccountCustomerManagedKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

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

			"key_vault_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.Any(
					// TODO 4.0: revert to only accepting key vault IDs as there is an explicit attribute for managed HSMs
					commonids.ValidateKeyVaultID,
					managedhsms.ValidateManagedHSMID,
				),
				ExactlyOneOf: []string{"managed_hsm_uri", "key_vault_id", "key_vault_uri"},
			},

			"key_vault_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
				ExactlyOneOf: []string{"managed_hsm_uri", "key_vault_id", "key_vault_uri"},
				Computed:     true,
			},

			"managed_hsm_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
				ExactlyOneOf: []string{"managed_hsm_uri", "key_vault_id", "key_vault_uri"},
			},

			"key_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	vaultsClient := keyVaultsClient.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	storageAccount, err := storageClient.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if storageAccount.AccountProperties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if d.IsNewResource() {
		// whilst this looks superfluous given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		if storageAccount.AccountProperties.Encryption != nil {
			if storageAccount.AccountProperties.Encryption.KeySource == storage.KeySourceMicrosoftKeyvault {
				return tf.ImportAsExistsError("azurerm_storage_account_customer_managed_key", id.ID())
			}
		}
	}

	keyVaultURI := ""
	if keyVaultURIRaw := d.Get("key_vault_uri").(string); keyVaultURIRaw != "" {
		keyVaultURI = keyVaultURIRaw
	} else if _, ok := d.GetOk("key_vault_id"); ok {
		keyVaultID, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
		if err != nil {
			return err
		}

		keyVault, err := vaultsClient.Get(ctx, *keyVaultID)
		if err != nil {
			return fmt.Errorf("retrieving Key Vault %q (Resource Group %q): %+v", keyVaultID.VaultName, keyVaultID.ResourceGroupName, err)
		}

		softDeleteEnabled := false
		purgeProtectionEnabled := false
		if model := keyVault.Model; model != nil {
			if esd := model.Properties.EnableSoftDelete; esd != nil {
				softDeleteEnabled = *esd
			}
			if epp := model.Properties.EnablePurgeProtection; epp != nil {
				purgeProtectionEnabled = *epp
			}
		}
		if !softDeleteEnabled || !purgeProtectionEnabled {
			return fmt.Errorf("Key Vault %q (Resource Group %q) must be configured for both Purge Protection and Soft Delete", keyVaultID.VaultName, keyVaultID.ResourceGroupName)
		}

		keyVaultBaseURL, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
		if err != nil {
			return fmt.Errorf("looking up Key Vault URI from %s: %+v", *keyVaultID, err)
		}

		keyVaultURI = *keyVaultBaseURL
	} else if _, ok := d.GetOk("managed_hsm_uri"); ok {
		keyVaultURI = d.Get("managed_hsm_uri").(string)
	}

	keyName := d.Get("key_name").(string)
	keyVersion := d.Get("key_version").(string)
	userAssignedIdentity := d.Get("user_assigned_identity_id").(string)
	federatedIdentityClientID := d.Get("federated_identity_client_id").(string)

	props := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
				},
				EncryptionIdentity: &storage.EncryptionIdentity{
					EncryptionUserAssignedIdentity: utils.String(userAssignedIdentity),
				},
				KeySource: storage.KeySourceMicrosoftKeyvault,
				KeyVaultProperties: &storage.KeyVaultProperties{
					KeyName:     utils.String(keyName),
					KeyVersion:  utils.String(keyVersion),
					KeyVaultURI: utils.String(keyVaultURI),
				},
			},
		},
	}

	if federatedIdentityClientID != "" {
		props.Encryption.EncryptionIdentity.EncryptionFederatedIdentityClientID = utils.String(federatedIdentityClientID)
	}

	if _, err = storageClient.Update(ctx, id.ResourceGroupName, id.StorageAccountName, props); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageAccountCustomerManagedKeyRead(d, meta)
}

func resourceStorageAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	env := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	storageAccount, err := storageClient.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if storageAccount.AccountProperties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if storageAccount.AccountProperties.Encryption == nil || storageAccount.AccountProperties.Encryption.KeySource != storage.KeySourceMicrosoftKeyvault {
		log.Printf("[DEBUG] Customer Managed Key was not defined for %s - removing from state!", id)
		d.SetId("")
		return nil
	}

	encryption := *storageAccount.AccountProperties.Encryption

	keyName := ""
	keyVaultURI := ""
	keyVersion := ""
	if props := encryption.KeyVaultProperties; props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultURI != nil {
			keyVaultURI = *props.KeyVaultURI
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
	}

	userAssignedIdentity := ""
	federatedIdentityClientID := ""
	if props := encryption.EncryptionIdentity; props != nil {
		if props.EncryptionUserAssignedIdentity != nil {
			userAssignedIdentity = *props.EncryptionUserAssignedIdentity
		}
		if props.EncryptionFederatedIdentityClientID != nil {
			federatedIdentityClientID = *props.EncryptionFederatedIdentityClientID
		}
	}

	if keyVaultURI == "" {
		return fmt.Errorf("retrieving %s: `properties.encryption.keyVaultProperties.keyVaultURI` was nil", id)
	}

	// now we have the key vault uri we can look up the ID

	// we can't look up the ID when using federated identity as the key will be under different tenant
	if federatedIdentityClientID == "" {
		isHSMURI, err := managedHsmParse.IsManagedHSMURI(keyVaultURI, &env)
		switch {
		case err != nil:
			{
				return fmt.Errorf("parsing Base Key Vault URI %q: %+v", keyVaultURI, err)
			}
		case isHSMURI:
			{
				d.Set("managed_hsm_uri", keyVaultURI)
			}
		case !isHSMURI:
			{
				d.Set("key_vault_uri", keyVaultURI)
				subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
				tmpKeyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultURI)
				if err != nil {
					return fmt.Errorf("retrieving Key Vault ID from the Base URI %q: %+v", keyVaultURI, err)
				}
				d.Set("key_vault_id", pointer.From(tmpKeyVaultID))
			}
		}
	}

	d.Set("storage_account_id", id.ID())
	d.Set("key_name", keyName)
	d.Set("key_version", keyVersion)
	d.Set("user_assigned_identity_id", userAssignedIdentity)
	d.Set("federated_identity_client_id", federatedIdentityClientID)

	return nil
}

func resourceStorageAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	// confirm it still exists prior to trying to update it, else we'll get an error
	storageAccount, err := storageClient.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// Since this isn't a real object, just modifying an existing object
	// "Delete" doesn't really make sense it should really be a "Revert to Default"
	// So instead of the Delete func actually deleting the Storage Account I am
	// making it reset the Storage Account to its default state
	props := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			Encryption: &storage.Encryption{
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
					File: &storage.EncryptionService{
						Enabled: utils.Bool(true),
					},
				},
				KeySource: storage.KeySourceMicrosoftStorage,
			},
		},
	}

	if _, err = storageClient.Update(ctx, id.ResourceGroupName, id.StorageAccountName, props); err != nil {
		return fmt.Errorf("removing Customer Managed Key for %s: %+v", id, err)
	}

	return nil
}
