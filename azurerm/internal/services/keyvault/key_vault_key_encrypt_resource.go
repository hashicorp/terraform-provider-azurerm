package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultKeyEncrypt() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultKeyEncryptCreate,
		Read:   schema.Noop,
		Delete: schema.Noop,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			return fmt.Errorf("resource `azurerm_key_vault_key_encrypt` does not support `import` operation")
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_key_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildId,
			},

			"plaintext": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.RSA15),
					string(keyvault.RSAOAEP),
					string(keyvault.RSAOAEP256),
				}, false),
			},

			"cipher_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmKeyVaultKeyEncryptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	plaintext := d.Get("plaintext").(string)
	keyVaultKeyIdRaw := d.Get("key_vault_key_id").(string)
	keyVaultKeyId, err := azure.ParseKeyVaultChildID(keyVaultKeyIdRaw)
	if err != nil {
		return err
	}

	params := keyvault.KeyOperationsParameters{
		Algorithm: keyvault.JSONWebKeyEncryptionAlgorithm(d.Get("algorithm").(string)),
		Value:     utils.String(plaintext),
	}
	result, err := client.Encrypt(ctx, keyVaultKeyId.KeyVaultBaseUrl, keyVaultKeyId.Name, keyVaultKeyId.Version, params)
	if err != nil {
		return fmt.Errorf("failed to encrypt '%s' using key %s: %+v", plaintext, keyVaultKeyIdRaw, err)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("key_vault_key_id", result.Kid)
	d.Set("cipher_text", result.Result)

	return nil
}
