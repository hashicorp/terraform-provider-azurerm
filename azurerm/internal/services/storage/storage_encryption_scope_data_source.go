package storage

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceStorageEncryptionScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageEncryptionScopeRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageEncryptionScopeName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"source": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_vault_key_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStorageEncryptionScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, accountId.ResourceGroup, accountId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q) was not found", name, accountId.Name, accountId.ResourceGroup)
		}

		return fmt.Errorf("retrieving Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, accountId.Name, accountId.ResourceGroup, err)
	}

	d.SetId(parse.NewEncryptionScopeID(accountId.SubscriptionId, accountId.ResourceGroup, accountId.Name, name).ID())

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
