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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_account_customer_managed_key -service-package-name storage -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id"

func resourceStorageAccountCustomerManagedKey() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
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

	if !features.FivePointOh() {
		resource.Schema["key_vault_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri", "key_vault_key_id"},
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		}

		resource.Schema["key_vault_uri"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri", "key_vault_key_id"},
			Deprecated:   "`key_vault_uri` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}

		resource.Schema["key_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Deprecated:   "`key_name` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}

		resource.Schema["key_version"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Deprecated:   "`key_version` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}

		resource.Schema["managed_hsm_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri", "key_vault_key_id"},
			Deprecated:   "`managed_hsm_key_id` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}

		resource.Schema["key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			ExactlyOneOf: []string{"managed_hsm_key_id", "key_vault_id", "key_vault_uri", "key_vault_key_id"},
			Deprecated:   "`key_vault_id` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}
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
		if existing.Model.Properties.Encryption != nil && pointer.From(existing.Model.Properties.Encryption.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
			return tf.ImportAsExistsError("azurerm_storage_account_customer_managed_key", id.ID())
		}
	}

	if features.FivePointOh() {
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
	} else {
		keyName := ""
		keyVersion := ""
		keyVaultURI := ""
		if !pluginsdk.IsExplicitlyNullInConfig(d, "key_vault_uri") {
			keyName = d.Get("key_name").(string)
			keyVersion = d.Get("key_version").(string)
			keyVaultURI = d.Get("key_vault_uri").(string)
		} else if !pluginsdk.IsExplicitlyNullInConfig(d, "key_vault_id") {
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
				return fmt.Errorf("%s must be configured for both Purge Protection and Soft Delete", keyVaultID)
			}

			keyVaultBaseURL, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
			if err != nil {
				return fmt.Errorf("looking up Key Vault URI from %s: %+v", *keyVaultID, err)
			}

			keyName = d.Get("key_name").(string)
			keyVersion = d.Get("key_version").(string)
			keyVaultURI = *keyVaultBaseURL
		} else if !pluginsdk.IsExplicitlyNullInConfig(d, "managed_hsm_key_id") {
			keyID, err := keyvault.ParseNestedItemID(d.Get("managed_hsm_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			keyName = keyID.Name
			keyVersion = keyID.Version
			keyVaultURI = keyID.KeyVaultBaseURL
		} else if !pluginsdk.IsExplicitlyNullInConfig(d, "key_vault_key_id") {
			keyID, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			keyName = keyID.Name
			keyVersion = keyID.Version
			keyVaultURI = keyID.KeyVaultBaseURL
		}

		userAssignedIdentity := d.Get("user_assigned_identity_id").(string)
		federatedIdentityClientID := d.Get("federated_identity_client_id").(string)

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
						UserAssignedIdentity: pointer.To(userAssignedIdentity),
					},
					KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointKeyvault),
					Keyvaultproperties: &storageaccounts.KeyVaultProperties{
						Keyname:     pointer.To(keyName),
						Keyversion:  pointer.To(keyVersion),
						Keyvaulturi: pointer.To(keyVaultURI),
					},
				},
			},
		}

		if federatedIdentityClientID != "" {
			payload.Properties.Encryption.Identity.FederatedIdentityClientId = pointer.To(federatedIdentityClientID)
		}

		if _, err = storageClient.Update(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating Customer Managed Key for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return resourceStorageAccountCustomerManagedKeyRead(d, meta)
}

func resourceStorageAccountCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	keyVaultsClient := meta.(*clients.Client).KeyVault

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
			if encryption := props.Encryption; encryption != nil && pointer.From(encryption.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
				enabled = true

				if features.FivePointOh() {
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
				} else {
					var keyID *keyvault.NestedItemID
					if kvProps := encryption.Keyvaultproperties; kvProps != nil {
						keyId, err := keyvault.NewNestedItemID(pointer.From(kvProps.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(kvProps.Keyname), pointer.From(kvProps.Keyversion))
						if err != nil {
							return err
						}

						d.Set("key_vault_key_id", keyId.ID())
						d.Set("key_name", kvProps.Keyname)
						d.Set("key_version", kvProps.Keyversion)
						d.Set("key_vault_uri", kvProps.Keyvaulturi)

						if keyId.IsManagedHSM() {
							d.Set("managed_hsm_key_id", keyId.ID())
						}

						keyID = keyId
					}

					federatedIdentityClientID := ""
					userAssignedIdentity := ""
					if identityProps := encryption.Identity; identityProps != nil {
						federatedIdentityClientID = pointer.From(identityProps.FederatedIdentityClientId)
						userAssignedIdentity = pointer.From(identityProps.UserAssignedIdentity)
					}
					// now we have the key vault uri we can look up the ID
					// we can't look up the ID when using federated identity as the key will be under different tenant
					keyVaultID := ""
					if federatedIdentityClientID == "" && keyID != nil && !keyID.IsManagedHSM() {
						subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
						tmpKeyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyID.KeyVaultBaseURL)
						if err != nil {
							return fmt.Errorf("retrieving Key Vault ID from the Base URI %q: %+v", keyID.KeyVaultBaseURL, err)
						}
						keyVaultID = pointer.From(tmpKeyVaultID)
					}
					d.Set("key_vault_id", keyVaultID)
					d.Set("user_assigned_identity_id", userAssignedIdentity)
					d.Set("federated_identity_client_id", federatedIdentityClientID)
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
