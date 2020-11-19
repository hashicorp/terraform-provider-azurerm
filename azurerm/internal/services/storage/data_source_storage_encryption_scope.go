package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStorageEncryptionScope() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStorageEncryptionScopeRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageEncryptionScopeName,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_vault_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmStorageEncryptionScopeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.EncryptionScopesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageAccountID, err := parse.AccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Encryption Scope %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Encryption Scope %q (Storage Account Name %q / Resource Group %q): %+v", name, storageAccountID.Name, storageAccountID.ResourceGroup, err)
	}

	d.SetId(parse.NewEncryptionScopeId(*storageAccountID, name).ID(subscriptionId))

	d.Set("name", resp.Name)
	d.Set("storage_account_id", storageAccountID.ID(subscriptionId))
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
