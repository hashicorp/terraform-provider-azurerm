package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVaultKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultKeyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:          schema.TypeString,
				Optional:      true, //todo required in 2.0
				Computed:      true, //todo removed in 2.0
				ForceNew:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"vault_uri"},
			},

			//todo remove in 2.0
			"vault_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Deprecated:    "This property has been deprecated in favour of the key_vault_id property. This will prevent a class of bugs as described in https://github.com/terraform-providers/terraform-provider-azurerm/issues/2396 and will be removed in version 2.0 of the provider",
				ValidateFunc:  validate.URLIsHTTPS,
				ConflictsWith: []string{"key_vault_id"},
			},

			"key_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"key_opts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"n": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"e": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmKeyVaultKeyRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultBaseUri := d.Get("vault_uri").(string)
	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	if keyVaultBaseUri == "" {
		if keyVaultId == "" {
			return fmt.Errorf("one of `key_vault_id` or `vault_uri` must be set")
		}

		pKeyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
		}

		keyVaultBaseUri = pKeyVaultBaseUrl
		d.Set("vault_uri", keyVaultBaseUri)
	} else {
		id, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, keyVaultBaseUri)
		if err != nil {
			return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultBaseUri, err)
		}
		if id == nil {
			return fmt.Errorf("Unable to locate the Resource ID for the Key Vault at URL %q: %s", keyVaultBaseUri, err)
		}
		d.Set("key_vault_id", id)
	}

	resp, err := client.GetKey(ctx, keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Key %q was not found in Key Vault at URI %q", name, keyVaultBaseUri)
		}

		return err
	}

	id := *resp.Key.Kid
	parsedId, err := azure.ParseKeyVaultChildID(id)
	if err != nil {
		return err
	}

	d.SetId(id)
	if key := resp.Key; key != nil {
		d.Set("key_type", string(key.Kty))

		options := flattenKeyVaultKeyDataSourceOptions(key.KeyOps)
		if err := d.Set("key_opts", options); err != nil {
			return err
		}

		d.Set("n", key.N)
		d.Set("e", key.E)
	}

	d.Set("version", parsedId.Version)

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenKeyVaultKeyDataSourceOptions(input *[]string) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, option := range *input {
			results = append(results, option)
		}
	}

	return results
}
