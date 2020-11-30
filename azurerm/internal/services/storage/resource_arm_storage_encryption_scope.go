package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageEncryptionScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageEncryptionScopeCreate,
		Read:   resourceArmStorageEncryptionScopeRead,
		Update: resourceArmStorageEncryptionScopeUpdate,
		Delete: resourceArmStorageEncryptionScopeDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EncryptionScopeID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageEncryptionScopeName,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"source": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyVault),
					string(storage.MicrosoftStorage),
				}, false),
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: storageValidate.KeyVaultChildId,
			},
		},
	}
}

func resourceArmStorageEncryptionScopeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	resourceId := parse.NewEncryptionScopeId(*accountId, name).ID("")
	existing, err := client.Get(ctx, accountId.ResourceGroup, accountId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
		}
	}
	if existing.EncryptionScopeProperties != nil && strings.EqualFold(string(existing.EncryptionScopeProperties.State), string(storage.Enabled)) {
		return tf.ImportAsExistsError("azurerm_storage_encryption_scope", resourceId)
	}

	if d.Get("source").(string) == string(storage.MicrosoftKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(storage.KeySourceMicrosoftKeyvault))
		}
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			Source: storage.EncryptionScopeSource(d.Get("source").(string)),
			State:  storage.Enabled,
			KeyVaultProperties: &storage.EncryptionScopeKeyVaultProperties{
				KeyURI: utils.String(d.Get("key_vault_key_id").(string)),
			},
		},
	}
	if _, err := client.Put(ctx, accountId.ResourceGroup, accountId.Name, name, props); err != nil {
		return fmt.Errorf("creating Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceArmStorageEncryptionScopeRead(d, meta)
}

func resourceArmStorageEncryptionScopeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	if d.Get("source").(string) == string(storage.MicrosoftKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is required when source is `%s`", string(storage.KeySourceMicrosoftKeyvault))
		}
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			Source: storage.EncryptionScopeSource(d.Get("source").(string)),
			State:  storage.Enabled,
			KeyVaultProperties: &storage.EncryptionScopeKeyVaultProperties{
				KeyURI: utils.String(d.Get("key_vault_key_id").(string)),
			},
		},
	}
	if _, err := client.Patch(ctx, id.ResourceGroup, id.StorageAccountName, id.Name, props); err != nil {
		return fmt.Errorf("updating Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", id.Name, id.StorageAccountName, id.ResourceGroup, err)
	}

	return resourceArmStorageEncryptionScopeRead(d, meta)
}

func resourceArmStorageEncryptionScopeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionScopeID(d.Id())
	if err != nil {
		return err
	}
	accountId := parse.NewAccountId(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", id.Name, id.StorageAccountName, id.ResourceGroup, err)
	}

	if resp.EncryptionScopeProperties == nil {
		return fmt.Errorf("retrieving Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): `properties` was nil", id.Name, id.StorageAccountName, id.ResourceGroup)
	}

	props := *resp.EncryptionScopeProperties
	if strings.EqualFold(string(props.State), string(storage.Disabled)) {
		log.Printf("[INFO] Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q) does not exist - removing from state", id.Name, id.StorageAccountName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Name)
	d.Set("storage_account_id", accountId.ID(""))
	if props := resp.EncryptionScopeProperties; props != nil {
		d.Set("source", flattenEncryptionScopeSource(props.Source))
		var keyId string
		if kv := props.KeyVaultProperties; kv != nil {
			if kv.KeyURI != nil {
				keyId = *kv.KeyURI
			}
		}
		d.Set("key_vault_key_id", keyId)
	}

	return nil
}

func resourceArmStorageEncryptionScopeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			State: storage.Disabled,
		},
	}

	if _, err = client.Put(ctx, id.ResourceGroup, id.StorageAccountName, id.Name, props); err != nil {
		return fmt.Errorf("disabling Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", id.Name, id.StorageAccountName, id.ResourceGroup, err)
	}

	return nil
}

func flattenEncryptionScopeSource(input storage.EncryptionScopeSource) string {
	// TODO: file a bug
	// the Storage API differs from every other API in Azure in that these Enum's can be returned case-insensitively
	if strings.EqualFold(string(input), string(storage.MicrosoftKeyVault)) {
		return string(storage.MicrosoftKeyVault)
	}

	return string(storage.MicrosoftStorage)
}
