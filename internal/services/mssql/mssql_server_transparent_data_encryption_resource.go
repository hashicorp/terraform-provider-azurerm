// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/encryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverkeys"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	mssqlValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlTransparentDataEncryption() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
			},

			"auto_rotation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["key_vault_key_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			DiffSuppressFunc: func(_, oldValue, newValue string, d *schema.ResourceData) bool {
				if newValue == "" {
					// If using `managed_hsm_key_id`, `key_vault_key_id` will also be set
					// ignore diff if the 2 are equal.
					raw := d.GetRawConfig().AsValueMap()["managed_hsm_key_id"]
					if raw.IsKnown() && !raw.IsNull() {
						return raw.AsString() == oldValue
					}
				}

				return false
			},
			ValidateFunc:  keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
			ConflictsWith: []string{"managed_hsm_key_id"},
		}

		r.Schema["managed_hsm_key_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			DiffSuppressFunc: func(_, oldValue, newValue string, d *schema.ResourceData) bool {
				if newValue == "" {
					// If using `key_vault_key_id` with MHSM key, `managed_hsm_key_id` will also be set
					// ignore diff if the 2 are equal.
					raw := d.GetRawConfig().AsValueMap()["key_vault_key_id"]
					if raw.IsKnown() && !raw.IsNull() {
						return raw.AsString() == oldValue
					}
				}

				return false
			},
			ValidateFunc:  keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
			ConflictsWith: []string{"key_vault_key_id"},
			Deprecated:    "`managed_hsm_key_id` has been deprecated in favour of `key_vault_key_id` and will be removed in v5.0 of the AzureRM provider",
		}
	}

	return r
}

func resourceMsSqlTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	serverKeysClient := meta.(*clients.Client).MSSQL.ServerKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewEncryptionProtectorID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, "current")

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	payload := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionprotectors.EncryptionProtectorProperties{
			AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
			ServerKeyName:       pointer.To(""),
			ServerKeyType:       encryptionprotectors.ServerKeyTypeServiceManaged,
		},
	}

	var key *keyvault.NestedItemID
	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyId, err := keyvault.ParseNestedItemID(v.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}
		key = keyId
	}

	if !features.FivePointOh() {
		if !pluginsdk.IsExplicitlyNullInConfig(d, "managed_hsm_key_id") {
			keyId, err := keyvault.ParseNestedItemID(d.Get("managed_hsm_key_id").(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}
			key = keyId
		}
	}

	if key != nil {
		keyVaultName, err := resourceMsSqlTransparentDataEncryptionKeyVaultName(key.KeyVaultBaseURL)
		if err != nil {
			return err
		}

		serverKeyName := fmt.Sprintf("%s_%s_%s", keyVaultName, key.Name, key.Version)

		serverKeyId := serverkeys.NewKeyID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, serverKeyName)
		serverKeyPayload := serverkeys.ServerKey{
			Properties: &serverkeys.ServerKeyProperties{
				AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
				ServerKeyType:       serverkeys.ServerKeyTypeAzureKeyVault,
				Uri:                 pointer.To(key.ID()),
			},
		}

		if err := serverKeysClient.CreateOrUpdateThenPoll(ctx, serverKeyId, serverKeyPayload); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", serverKeyId, err)
		}

		// Update TDE properties to reflect usage of Server Key
		payload.Properties.ServerKeyName = pointer.To(serverKeyName)
		payload.Properties.ServerKeyType = encryptionprotectors.ServerKeyTypeAzureKeyVault
	}

	err = client.CreateOrUpdateThenPoll(ctx, *serverId, payload)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlTransparentDataEncryptionRead(d, meta)
}

func resourceMsSqlTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("server_id", serverId.ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			var key *keyvault.NestedItemID
			if props.ServerKeyType == encryptionprotectors.ServerKeyTypeAzureKeyVault && props.Uri != nil {
				key, err = keyvault.ParseNestedItemID(*props.Uri, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}
			}

			var hsmKeyId, keyVaultKeyId string
			if key != nil {
				keyVaultKeyId = key.ID()
				if !features.FivePointOh() && key.IsManagedHSM() {
					hsmKeyId = keyVaultKeyId
				}
			}

			if !features.FivePointOh() {
				if err := d.Set("managed_hsm_key_id", hsmKeyId); err != nil {
					return fmt.Errorf("setting `managed_hsm_key_id`: %+v", err)
				}
			}

			if err := d.Set("key_vault_key_id", keyVaultKeyId); err != nil {
				return fmt.Errorf("setting `key_vault_key_id`: %+v", err)
			}

			if err := d.Set("auto_rotation_enabled", pointer.From(props.AutoRotationEnabled)); err != nil {
				return fmt.Errorf("setting `auto_rotation_enabled`: %+v", err)
			}
		}
	}

	return nil
}

func resourceMsSqlTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	encryptionProtector := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionprotectors.EncryptionProtectorProperties{
			ServerKeyType: encryptionprotectors.ServerKeyTypeServiceManaged,
			ServerKeyName: pointer.To(""),
		},
	}

	err = client.CreateOrUpdateThenPoll(ctx, serverId, encryptionProtector)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func resourceMsSqlTransparentDataEncryptionKeyVaultName(keyVaultURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(keyVaultURL)
	if err != nil {
		return "", fmt.Errorf("parsing Key Vault URL (%s): %+v", keyVaultURL, err)
	}

	return strings.Split(parsedURL.Host, ".")[0], nil
}
