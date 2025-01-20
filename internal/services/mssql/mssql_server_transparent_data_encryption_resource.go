// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/encryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmHelpers "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	mhsmParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	mssqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlTransparentDataEncryption() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlTransparentDataEncryptionRead,
		Update: resourceMsSqlTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlTransparentDataEncryptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EncryptionProtectorID(id)

			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.MsSqlTransparentDataEncryptionV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: mssqlValidate.ServerID,
			},

			"key_vault_key_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  keyVaultValidate.NestedItemId,
				ConflictsWith: []string{"managed_hsm_key_id"},
			},

			"managed_hsm_key_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validate.ManagedHSMDataPlaneVersionedKeyID,
				ConflictsWith: []string{"key_vault_key_id"},
			},

			"auto_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceMsSqlTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	serverKeysClient := meta.(*clients.Client).MSSQL.ServerKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	var serverKey serverkeys.ServerKey

	// Default values for Service Managed keys. Will update to AKV values if key_vault_key_id references a key.
	serverKeyName := ""
	serverKeyType := serverkeys.ServerKeyTypeServiceManaged

	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKeyId := strings.TrimSpace(v.(string))
		// Update the server key type to AKV
		serverKeyType = serverkeys.ServerKeyTypeAzureKeyVault

		// Set the SQL Server Key properties
		serverKeyProperties := serverkeys.ServerKeyProperties{
			ServerKeyType:       serverKeyType,
			Uri:                 &keyVaultKeyId,
			AutoRotationEnabled: utils.Bool(d.Get("auto_rotation_enabled").(bool)),
		}
		serverKey.Properties = &serverKeyProperties

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
			serverKeyName = fmt.Sprintf("%s_%s_%s", vaultName, keyName, keyVersion)
		} else {
			return fmt.Errorf("key vault key id must be a reference to a key, but got: %s", keyId.NestedItemType)
		}
	}

	if v, ok := d.GetOk("managed_hsm_key_id"); ok {
		mhsmKeyId := strings.TrimSpace(v.(string))
		// Update the server key type to AKV
		serverKeyType = serverkeys.ServerKeyTypeAzureKeyVault

		// Set the SQL Server Key properties z
		serverKeyProperties := serverkeys.ServerKeyProperties{
			ServerKeyType:       serverKeyType,
			Uri:                 &mhsmKeyId,
			AutoRotationEnabled: utils.Bool(d.Get("auto_rotation_enabled").(bool)),
		}
		serverKey.Properties = &serverKeyProperties

		// Make sure it's a key, if not, throw an error
		keyId, err := mhsmParser.ManagedHSMDataPlaneVersionedKeyID(mhsmKeyId, nil)
		if err != nil {
			return fmt.Errorf("failed to parse '%s' as HSM key ID", mhsmKeyId)
		}

		// Extract the vault name from the keyvault base url
		idURL, err := url.ParseRequestURI(keyId.BaseUri())
		if err != nil {
			return fmt.Errorf("unable to parse key vault hostname: %s", keyId.BaseUri())
		}

		hostParts := strings.Split(idURL.Host, ".")
		vaultName := hostParts[0]

		// Create the key path for the Encryption Protector. Format is: {vaultname}_{key}_{key_version}
		serverKeyName = fmt.Sprintf("%s_%s_%s", vaultName, keyId.KeyName, keyId.KeyVersion)
	}

	keyType := encryptionprotectors.ServerKeyTypeServiceManaged
	if serverKeyType == serverkeys.ServerKeyTypeAzureKeyVault {
		keyType = encryptionprotectors.ServerKeyTypeAzureKeyVault
	}

	// Service managed doesn't require a key name
	encryptionProtectorProperties := encryptionprotectors.EncryptionProtectorProperties{
		ServerKeyType:       keyType,
		ServerKeyName:       &serverKeyName,
		AutoRotationEnabled: utils.Bool(d.Get("auto_rotation_enabled").(bool)),
	}

	keyId := serverkeys.NewKeyID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, serverKeyName)

	// Only create a server key if the properties have been set
	if serverKey.Properties != nil {
		// Create a key on the server
		err = serverKeysClient.CreateOrUpdateThenPoll(ctx, keyId, serverKey)
		if err != nil {
			return fmt.Errorf("creating/updating server key for %s: %+v", serverId, err)
		}
	}

	encryptionProtectorObject := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionProtectorProperties,
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewEncryptionProtectorID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, "current")

	err = encryptionProtectorClient.CreateOrUpdateThenPoll(ctx, *serverId, encryptionProtectorObject)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlTransparentDataEncryptionRead(d, meta)
}

func resourceMsSqlTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	env := meta.(*clients.Client).Account.Environment

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}
	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := encryptionProtectorClient.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request for %s: %v", id, err)
	}

	d.Set("server_id", serverId.ID())

	log.Printf("[INFO] Encryption protector key type is %s", resp.Model.Properties.ServerKeyType)

	keyId := ""
	autoRotationEnabled := false
	// Only set the key type if it's an AKV key. For service managed, we can omit the setting the key_vault_key_id
	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.ServerKeyType == encryptionprotectors.ServerKeyTypeAzureKeyVault {
		log.Printf("[INFO] Setting Key Vault URI to %s", *resp.Model.Properties.Uri)

		keyId = *resp.Model.Properties.Uri

		// autoRotation is only for AKV keys
		if resp.Model.Properties.AutoRotationEnabled != nil {
			autoRotationEnabled = *resp.Model.Properties.AutoRotationEnabled
		}
	}

	hsmKey := ""
	keyVaultKeyId := ""
	if keyId != "" {
		isHSMURI, err, _, _ := managedHsmHelpers.IsManagedHSMURI(env, keyId)
		if err != nil {
			return err
		}

		if isHSMURI {
			hsmKey = keyId
		} else {
			keyVaultKeyId = keyId
		}
	}

	if err := d.Set("managed_hsm_key_id", hsmKey); err != nil {
		return fmt.Errorf("setting `managed_hsm_key_id`: %+v", err)
	}

	if err := d.Set("key_vault_key_id", keyVaultKeyId); err != nil {
		return fmt.Errorf("setting `key_vault_key_id`: %+v", err)
	}

	if err := d.Set("auto_rotation_enabled", autoRotationEnabled); err != nil {
		return fmt.Errorf("setting `auto_rotation_enabled`: %+v", err)
	}

	return nil
}

func resourceMsSqlTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	encryptionProtectorClient := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	serverKeyName := ""

	// Service managed doesn't require a key name
	encryptionProtector := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionprotectors.EncryptionProtectorProperties{
			ServerKeyType: encryptionprotectors.ServerKeyTypeServiceManaged,
			ServerKeyName: &serverKeyName,
		},
	}

	err = encryptionProtectorClient.CreateOrUpdateThenPoll(ctx, serverId, encryptionProtector)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return nil
}
