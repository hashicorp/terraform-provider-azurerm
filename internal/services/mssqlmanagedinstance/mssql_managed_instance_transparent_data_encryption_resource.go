// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstanceencryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstancekeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlManagedInstanceTransparentDataEncryption() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlManagedInstanceTransparentDataEncryptionRead,
		Update: resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlManagedInstanceTransparentDataEncryptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedInstanceEncryptionProtectorID(id)

			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"managed_instance_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagedInstanceID,
			},
			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},
			"auto_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient
	managedInstanceKeysClient := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedInstanceId, err := commonids.ParseSqlManagedInstanceID(d.Get("managed_instance_id").(string))
	if err != nil {
		return err
	}

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	var managedInstanceKey managedinstancekeys.ManagedInstanceKey

	// Default values for Service Managed keys. Will update to AKV values if key_vault_key_id references a key.
	managedInstanceKeyName := ""
	managedInstanceKeyType := managedinstancekeys.ServerKeyTypeServiceManaged

	keyVaultKeyId := strings.TrimSpace(d.Get("key_vault_key_id").(string))

	// If it has content, then we assume it's a key vault key id
	if keyVaultKeyId != "" {
		// Update the server key type to AKV
		managedInstanceKeyType = managedinstancekeys.ServerKeyTypeAzureKeyVault

		// Set the SQL Managed Instance Key properties
		managedInstanceKeyProperties := managedinstancekeys.ManagedInstanceKeyProperties{
			ServerKeyType:       managedInstanceKeyType,
			Uri:                 &keyVaultKeyId,
			AutoRotationEnabled: utils.Bool(d.Get("auto_rotation_enabled").(bool)),
		}
		managedInstanceKey.Properties = &managedInstanceKeyProperties

		// Set the encryption protector properties
		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		// Make sure it's a key, if not, throw an error
		if keyId.NestedItemType == keyVaultParser.NestedItemTypeKey {
			keyName := keyId.Name
			keyVersion := keyId.Version

			// Extract the vault name from the keyvault base url
			idURL, err := url.ParseRequestURI(keyId.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("unable to parse key vault hostname: %s", keyId.KeyVaultBaseUrl)
			}

			hostParts := strings.Split(idURL.Host, ".")
			vaultName := hostParts[0]

			// Create the key path for the Encryption Protector. Format is: {vaultname}_{key}_{key_version}
			managedInstanceKeyName = fmt.Sprintf("%s_%s_%s", vaultName, keyName, keyVersion)
		} else {
			return fmt.Errorf("key vault key id must be a reference to a key, but got: %s", keyId.NestedItemType)
		}
	}

	keyType := managedinstanceencryptionprotectors.ServerKeyTypeServiceManaged
	if managedInstanceKeyType == managedinstancekeys.ServerKeyTypeAzureKeyVault {
		keyType = managedinstanceencryptionprotectors.ServerKeyTypeAzureKeyVault
	}
	// Service managed doesn't require a key name
	encryptionProtectorProperties := managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtectorProperties{
		ServerKeyType:       keyType,
		ServerKeyName:       &managedInstanceKeyName,
		AutoRotationEnabled: utils.Bool(d.Get("auto_rotation_enabled").(bool)),
	}
	managedInstanceKeyId := managedinstancekeys.NewManagedInstanceKeyID(managedInstanceId.SubscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, managedInstanceKeyName)

	// Only create a managed instance key if the properties have been set
	if managedInstanceKey.Properties != nil {
		// Create a key on the managed instance
		err = managedInstanceKeysClient.CreateOrUpdateThenPoll(ctx, managedInstanceKeyId, managedInstanceKey)
		if err != nil {
			return fmt.Errorf("creating/updating managed instance key for %s: %+v", managedInstanceId, err)
		}
	}

	encryptionProtectorObject := managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtector{
		Properties: &encryptionProtectorProperties,
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewManagedInstanceEncryptionProtectorID(managedInstanceId.SubscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, "current")

	err = encryptionProtectorClient.CreateOrUpdateThenPoll(ctx, *managedInstanceId, encryptionProtectorObject)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceTransparentDataEncryptionRead(d, meta)
}

func resourceMsSqlManagedInstanceTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceEncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	resp, err := encryptionProtectorClient.Get(ctx, managedInstanceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request for %s: %v", id, err)
	}

	d.Set("managed_instance_id", managedInstanceId.ID())

	if resp.Model != nil && resp.Model.Properties != nil {
		log.Printf("[INFO] Encryption protector key type is %s", resp.Model.Properties.ServerKeyType)

		keyVaultKeyId := ""
		autoRotationEnabled := false

		// Only set the key type if it's an AKV key. For service managed, we can omit the setting the key_vault_key_id
		if resp.Model.Properties.ServerKeyType == managedinstanceencryptionprotectors.ServerKeyTypeAzureKeyVault {
			if resp.Model.Properties.Uri != nil {
				log.Printf("[INFO] Setting Key Vault URI to %s", *resp.Model.Properties.Uri)
				keyVaultKeyId = *resp.Model.Properties.Uri
			}

			// autoRotation is only for AKV keys
			if resp.Model.Properties.AutoRotationEnabled != nil {
				autoRotationEnabled = *resp.Model.Properties.AutoRotationEnabled
			}
		}

		if err := d.Set("key_vault_key_id", keyVaultKeyId); err != nil {
			return fmt.Errorf("setting `key_vault_key_id`: %+v", err)
		}

		if err := d.Set("auto_rotation_enabled", autoRotationEnabled); err != nil {
			return fmt.Errorf("setting `auto_rotation_enabled`: %+v", err)
		}
	}
	return nil
}

func resourceMsSqlManagedInstanceTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	encryptionProtectorClient := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceEncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	managedInstanceKeyName := ""

	// Service managed doesn't require a key name
	encryptionProtector := managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtector{
		Properties: &managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtectorProperties{
			ServerKeyType: managedinstanceencryptionprotectors.ServerKeyTypeServiceManaged,
			ServerKeyName: &managedInstanceKeyName,
		},
	}

	err = encryptionProtectorClient.CreateOrUpdateThenPoll(ctx, managedInstanceId, encryptionProtector)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return nil
}
