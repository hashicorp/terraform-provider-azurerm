// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageAccountCustomerManagedKey() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateKeyVaultID,
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri"},
			},

			"key_vault_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri"},
				Computed:     true,
			},

			"managed_hsm_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(validate.ManagedHSMDataPlaneVersionedKeyID, validate.ManagedHSMDataPlaneVersionlessKeyID),
				ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri"},
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

	if !features.FourPointOhBeta() {
		resource.Schema["key_vault_id"].ValidateFunc = validation.Any(
			commonids.ValidateKeyVaultID,
			managedhsms.ValidateManagedHSMID,
		)
	}

	return resource
}

func resourceStorageAccountCustomerManagedKeyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
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
		// whilst this looks superfluous given encryption is enabled by default, due to the way
		// the Azure API works this technically can be nil
		if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Encryption != nil && existing.Model.Properties.Encryption.KeySource != nil {
			if *existing.Model.Properties.Encryption.KeySource == storageaccounts.KeySourceMicrosoftPointKeyvault {
				return tf.ImportAsExistsError("azurerm_storage_account_customer_managed_key", id.ID())
			}
		}
	}

	keyName := ""
	keyVersion := ""
	keyVaultURI := ""
	if keyVaultURIRaw := d.Get("key_vault_uri").(string); keyVaultURIRaw != "" {
		keyName = d.Get("key_name").(string)
		keyVersion = d.Get("key_version").(string)
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

		keyName = d.Get("key_name").(string)
		keyVersion = d.Get("key_version").(string)
		keyVaultURI = *keyVaultBaseURL
	} else if managedHSMKeyId, ok := d.GetOk("managed_hsm_key_id"); ok {
		if keyId, err := parse.ManagedHSMDataPlaneVersionedKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = keyId.KeyName
			keyVersion = keyId.KeyVersion
			keyVaultURI = keyId.BaseUri()
		} else if keyId, err := parse.ManagedHSMDataPlaneVersionlessKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = keyId.KeyName
			keyVersion = ""
			keyVaultURI = keyId.BaseUri()
		} else {
			return fmt.Errorf("Failed to parse '%s' as HSM key ID", managedHSMKeyId)
		}
	}

	userAssignedIdentity := d.Get("user_assigned_identity_id").(string)
	federatedIdentityClientID := d.Get("federated_identity_client_id").(string)

	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			Encryption: &storageaccounts.Encryption{
				Services: &storageaccounts.EncryptionServices{
					Blob: &storageaccounts.EncryptionService{
						Enabled: utils.Bool(true),
					},
					File: &storageaccounts.EncryptionService{
						Enabled: utils.Bool(true),
					},
				},
				Identity: &storageaccounts.EncryptionIdentity{
					UserAssignedIdentity: utils.String(userAssignedIdentity),
				},
				KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointKeyvault),
				Keyvaultproperties: &storageaccounts.KeyVaultProperties{
					Keyname:     utils.String(keyName),
					Keyversion:  utils.String(keyVersion),
					Keyvaulturi: utils.String(keyVaultURI),
				},
			},
		},
	}

	if federatedIdentityClientID != "" {
		payload.Properties.Encryption.Identity.FederatedIdentityClientId = utils.String(federatedIdentityClientID)
	}
	if _, err = storageClient.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating Customer Managed Key for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageAccountCustomerManagedKeyRead(d, meta)
}

func resourceStorageAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	keyVaultsClient := meta.(*clients.Client).KeyVault
	env := meta.(*clients.Client).Account.Environment
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

	d.Set("storage_account_id", id.ID())

	enabled := false
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if encryption := props.Encryption; encryption != nil && encryption.KeySource != nil && *encryption.KeySource == storageaccounts.KeySourceMicrosoftPointKeyvault {
				enabled = true

				customerManagedKey := flattenCustomerManagedKey(encryption.Keyvaultproperties, env.KeyVault, env.ManagedHSM)
				d.Set("key_name", customerManagedKey.keyName)
				d.Set("key_version", customerManagedKey.keyVersion)
				d.Set("key_vault_uri", customerManagedKey.keyVaultBaseUrl)
				d.Set("managed_hsm_key_id", customerManagedKey.managedHsmKeyUri)

				federatedIdentityClientID := ""
				userAssignedIdentity := ""
				if identityProps := encryption.Identity; identityProps != nil {
					federatedIdentityClientID = pointer.From(identityProps.FederatedIdentityClientId)
					userAssignedIdentity = pointer.From(identityProps.UserAssignedIdentity)
				}
				// now we have the key vault uri we can look up the ID
				// we can't look up the ID when using federated identity as the key will be under different tenant
				keyVaultID := ""
				if federatedIdentityClientID == "" && customerManagedKey.keyVaultBaseUrl != "" {
					subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
					tmpKeyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, customerManagedKey.keyVaultBaseUrl)
					if err != nil {
						return fmt.Errorf("retrieving Key Vault ID from the Base URI %q: %+v", customerManagedKey.keyVaultBaseUrl, err)
					}
					keyVaultID = pointer.From(tmpKeyVaultID)
				}
				d.Set("key_vault_id", keyVaultID)

				d.Set("user_assigned_identity_id", userAssignedIdentity)
				d.Set("federated_identity_client_id", federatedIdentityClientID)
			}
		}
	}

	if !enabled {
		log.Printf("[DEBUG] Customer Managed Key was not defined for %s - removing from state!", id)
		d.SetId("")
		return nil
	}

	return nil
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
						Enabled: utils.Bool(true),
					},
					File: &storageaccounts.EncryptionService{
						Enabled: utils.Bool(true),
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
