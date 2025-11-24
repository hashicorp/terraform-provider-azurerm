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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/encryptionscopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_encryption_scope -service-package-name storage -properties "name" -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id" -test-name "keyVaultKey"

func resourceStorageEncryptionScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageEncryptionScopeCreate,
		Read:   resourceStorageEncryptionScopeRead,
		Update: resourceStorageEncryptionScopeUpdate,
		Delete: resourceStorageEncryptionScopeDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&encryptionscopes.EncryptionScopeId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&encryptionscopes.EncryptionScopeId{}),
		},

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
					string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault),
					string(encryptionscopes.EncryptionScopeSourceMicrosoftPointStorage),
				}, false),
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
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
	client := meta.(*clients.Client).Storage.ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := encryptionscopes.NewEncryptionScopeID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
	}
	if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.State != nil {
		if *existing.Model.Properties.State == encryptionscopes.EncryptionScopeStateEnabled {
			return tf.ImportAsExistsError("azurerm_storage_encryption_scope", id.ID())
		}
	}

	if d.Get("source").(string) == string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault))
		}
	}

	payload := encryptionscopes.EncryptionScope{
		Properties: &encryptionscopes.EncryptionScopeProperties{
			Source: pointer.To(encryptionscopes.EncryptionScopeSource(d.Get("source").(string))),
			State:  pointer.To(encryptionscopes.EncryptionScopeStateEnabled),
			KeyVaultProperties: &encryptionscopes.EncryptionScopeKeyVaultProperties{
				KeyUri: pointer.To(d.Get("key_vault_key_id").(string)),
			},
		},
	}
	if v, ok := d.GetOk("infrastructure_encryption_required"); ok {
		payload.Properties.RequireInfrastructureEncryption = pointer.To(v.(bool))
	}

	if _, err := client.Put(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageEncryptionScopeRead(d, meta)
}

func resourceStorageEncryptionScopeUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	if d.Get("source").(string) == string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(encryptionscopes.EncryptionScopeSourceMicrosoftPointKeyVault))
		}
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	payload := existing.Model
	payload.Properties.State = pointer.To(encryptionscopes.EncryptionScopeStateEnabled)
	if d.HasChange("key_vault_key_id") {
		payload.Properties.KeyVaultProperties = &encryptionscopes.EncryptionScopeKeyVaultProperties{
			KeyUri: pointer.To(d.Get("key_vault_key_id").(string)),
		}
	}
	if d.HasChange("source") {
		payload.Properties.Source = pointer.To(encryptionscopes.EncryptionScopeSource(d.Get("source").(string)))
	}

	if _, err := client.Patch(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStorageEncryptionScopeRead(d, meta)
}

func resourceStorageEncryptionScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.EncryptionScopeName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.State != nil && *props.State == encryptionscopes.EncryptionScopeStateDisabled {
				log.Printf("[INFO] %s was not configured - removing from state", id)
				d.SetId("")
				return nil
			}

			d.Set("infrastructure_encryption_required", props.RequireInfrastructureEncryption)
			d.Set("source", string(pointer.From(props.Source)))

			keyVaultKeyUri := ""
			if props.KeyVaultProperties != nil {
				keyVaultKeyUri = pointer.From(props.KeyVaultProperties.KeyUri)
			}
			d.Set("key_vault_key_id", keyVaultKeyUri)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceStorageEncryptionScopeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	payload := encryptionscopes.EncryptionScope{
		Properties: &encryptionscopes.EncryptionScopeProperties{
			State: pointer.To(encryptionscopes.EncryptionScopeStateDisabled),
		},
	}

	if _, err = client.Put(ctx, *id, payload); err != nil {
		return fmt.Errorf("disabling %s: %+v", id, err)
	}

	return nil
}
