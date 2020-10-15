package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVaultKeyDecrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultKeyDecryptRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildId,
			},

			"payload": {
				Type:     schema.TypeString,
				Required: true,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.RSA15),
					string(keyvault.RSAOAEP),
					string(keyvault.RSAOAEP256),
				}, false),
			},

			"plaintext": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmKeyVaultKeyDecryptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	payload := d.Get("payload").(string)
	keyVaultKeyIdRaw := d.Get("key_vault_key_id").(string)
	keyVaultKeyId, err := azure.ParseKeyVaultChildID(keyVaultKeyIdRaw)
	if err != nil {
		return err
	}

	params := keyvault.KeyOperationsParameters{
		Algorithm: keyvault.JSONWebKeyEncryptionAlgorithm(d.Get("algorithm").(string)),
		Value:     utils.String(payload),
	}
	result, err := client.Decrypt(ctx, keyVaultKeyId.KeyVaultBaseUrl, keyVaultKeyId.Name, keyVaultKeyId.Version, params)
	if err != nil {
		return fmt.Errorf("failed to decrypt '%s' using key %s: %+v", payload, keyVaultKeyIdRaw, err)
	}

	if result.Result == nil {
		return fmt.Errorf("nil decrypt result of '%s' using key %s: %+v", payload, keyVaultKeyIdRaw, err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("plaintext", result.Result)

	return nil
}
