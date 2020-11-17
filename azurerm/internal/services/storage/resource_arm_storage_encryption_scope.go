package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
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
		Create: resourceArmStorageEncryptionScopeCreateUpdate,
		Read:   resourceArmStorageEncryptionScopeRead,
		Update: resourceArmStorageEncryptionScopeCreateUpdate,
		Delete: resourceArmStorageEncryptionScopeDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.StorageEncryptionScopeID(id)
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
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MicrosoftKeyVault),
					string(storage.MicrosoftStorage),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"key_vault_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: storageValidate.KeyVaultChildId,
			},
		},
	}
}

func resourceArmStorageEncryptionScopeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageAccountIDRaw := d.Get("storage_account_id").(string)
	storageAccountID, _ := parse.ParseAccountID(storageAccountIDRaw)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, storageAccountID.Name, storageAccountID.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" && existing.EncryptionScopeProperties != nil && existing.EncryptionScopeProperties.State == storage.Enabled {
			return tf.ImportAsExistsError("azurerm_storage_encryption_scope", *existing.ID)
		}
	}

	if d.Get("source").(string) == string(storage.MicrosoftKeyVault) {
		if _, ok := d.GetOk("key_vault_key_id"); !ok {
			return fmt.Errorf("`key_vault_key_id` is necessary when the source is `Microsoft.KeyVault`")
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

	if _, err := client.Put(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, name, props); err != nil {
		return fmt.Errorf("creating Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("storage Encryption Scope %q (Storage Account Name %q / Resource Group %q) ID is empty or nil", name, storageAccountID.Name, storageAccountID.ResourceGroup)
	}
	if props := resp.EncryptionScopeProperties; props != nil {
		if props.State == storage.Disabled {
			return fmt.Errorf("storage Encryption Scope %q (Storage Account Name %q / Resource Group %q) ID is not enabled", name, storageAccountID.Name, storageAccountID.ResourceGroup)
		}
	}

	d.SetId(*resp.ID)
	return resourceArmStorageEncryptionScopeRead(d, meta)
}

func resourceArmStorageEncryptionScopeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StorageAccName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", id.Name, id.StorageAccName, id.ResourceGroup, err)
	}
	d.Set("name", resp.Name)
	d.Set("storage_account_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", client.SubscriptionID, id.ResourceGroup, id.StorageAccName))
	if props := resp.EncryptionScopeProperties; props != nil {
		d.Set("source", string(props.Source))
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

	id, err := parse.StorageEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	props := storage.EncryptionScope{
		EncryptionScopeProperties: &storage.EncryptionScopeProperties{
			State: storage.Disabled,
		},
	}

	_, err = client.Put(ctx, id.ResourceGroup, id.StorageAccName, id.Name, props)
	if err != nil {
		return fmt.Errorf("disable Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", id.Name, id.StorageAccName, id.ResourceGroup, err)
	}

	return nil
}
