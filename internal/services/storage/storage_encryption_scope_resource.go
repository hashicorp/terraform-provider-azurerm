// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageEncryptionScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageEncryptionScopeCreate,
		Read:   resourceStorageEncryptionScopeRead,
		Update: resourceStorageEncryptionScopeUpdate,
		Delete: resourceStorageEncryptionScopeDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := encryptionscopes.ParseEncryptionScopeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageEncryptionScopeName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"source": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.EncryptionScopeSourceMicrosoftKeyVault),
					string(storage.EncryptionScopeSourceMicrosoftStorage),
				}, false),
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildIDWithOptionalVersion,
			},

			"infrastructure_encryption_required": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceStorageEncryptionScopeCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := encryptionscopes.NewEncryptionScopeID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroupName, id.StorageAccountName, id.EncryptionScopeName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
	}
	if existing.EncryptionScopeProperties != nil && strings.EqualFold(string(existing.EncryptionScopeProperties.State), string(storage.EncryptionScopeStateEnabled)) {
		return tf.ImportAsExistsError("azurerm_storage_encryption_scope", id.ID())
	}

	if d.Get("source").(string) == string(storage.EncryptionScopeSourceMicrosoftKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(storage.KeySourceMicrosoftKeyvault))
		}
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			Source: storage.EncryptionScopeSource(d.Get("source").(string)),
			State:  storage.EncryptionScopeStateEnabled,
			KeyVaultProperties: &storage.EncryptionScopeKeyVaultProperties{
				KeyURI: utils.String(d.Get("key_vault_key_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("infrastructure_encryption_required"); ok {
		props.EncryptionScopeProperties.RequireInfrastructureEncryption = utils.Bool(v.(bool))
	}

	if _, err := client.Put(ctx, id.ResourceGroupName, id.StorageAccountName, id.EncryptionScopeName, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageEncryptionScopeRead(d, meta)
}

func resourceStorageEncryptionScopeUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	if d.Get("source").(string) == string(storage.EncryptionScopeSourceMicrosoftKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(storage.KeySourceMicrosoftKeyvault))
		}
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			Source: storage.EncryptionScopeSource(d.Get("source").(string)),
			State:  storage.EncryptionScopeStateEnabled,
			KeyVaultProperties: &storage.EncryptionScopeKeyVaultProperties{
				KeyURI: utils.String(d.Get("key_vault_key_id").(string)),
			},
		},
	}
	if _, err := client.Patch(ctx, id.ResourceGroupName, id.StorageAccountName, id.EncryptionScopeName, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStorageEncryptionScopeRead(d, meta)
}

func resourceStorageEncryptionScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.StorageAccountName, id.EncryptionScopeName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.EncryptionScopeProperties == nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	props := *resp.EncryptionScopeProperties
	if strings.EqualFold(string(props.State), string(storage.EncryptionScopeStateDisabled)) {
		log.Printf("[INFO] %s was not configured - removing from state", id)
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Name)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())
	if props := resp.EncryptionScopeProperties; props != nil {
		d.Set("source", flattenEncryptionScopeSource(props.Source))
		var keyId string
		if kv := props.KeyVaultProperties; kv != nil {
			if kv.KeyURI != nil {
				keyId = *kv.KeyURI
			}
		}
		d.Set("key_vault_key_id", keyId)
		if props.RequireInfrastructureEncryption != nil {
			d.Set("infrastructure_encryption_required", props.RequireInfrastructureEncryption)
		}
	}

	return nil
}

func resourceStorageEncryptionScopeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			State: storage.EncryptionScopeStateDisabled,
		},
	}

	if _, err = client.Put(ctx, id.ResourceGroupName, id.StorageAccountName, id.EncryptionScopeName, props); err != nil {
		return fmt.Errorf("disabling %s: %+v", id, err)
	}

	return nil
}

func flattenEncryptionScopeSource(input storage.EncryptionScopeSource) string {
	// TODO: remove this logic when migrated to hashicorp/go-azure-sdk and the new base layer
	// the Storage API differs from every other API in Azure in that these Enum's can be returned case-insensitively
	if strings.EqualFold(string(input), string(storage.EncryptionScopeSourceMicrosoftKeyVault)) {
		return string(storage.EncryptionScopeSourceMicrosoftKeyVault)
	}

	return string(storage.EncryptionScopeSourceMicrosoftStorage)
}
