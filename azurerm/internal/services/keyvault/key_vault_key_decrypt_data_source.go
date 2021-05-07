package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultKeyDecrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyVaultKeyDecryptRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NestedItemId,
			},

			"encrypted_base64url_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringIsRawBase64Url,
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

func dataSourceKeyVaultKeyDecryptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	data := d.Get("encrypted_base64url_data").(string)
	keyVaultKeyIdRaw := d.Get("key_vault_key_id").(string)
	keyVaultKeyId, err := parse.ParseNestedItemID(keyVaultKeyIdRaw)
	if err != nil {
		return err
	}

	params := keyvault.KeyOperationsParameters{
		Algorithm: keyvault.JSONWebKeyEncryptionAlgorithm(d.Get("algorithm").(string)),
		Value:     utils.String(data),
	}
	result, err := client.Decrypt(ctx, keyVaultKeyId.KeyVaultBaseUrl, keyVaultKeyId.Name, keyVaultKeyId.Version, params)
	if err != nil {
		return fmt.Errorf("failed to decrypt '%s' using key %s: %+v", data, keyVaultKeyIdRaw, err)
	}

	if result.Result == nil {
		return fmt.Errorf("nil decrypt result of '%s' using key %s: %+v", data, keyVaultKeyIdRaw, err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("plaintext", result.Result)

	return nil
}
