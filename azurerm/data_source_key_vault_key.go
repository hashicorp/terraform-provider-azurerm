package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVaultKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"vault_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.URLIsHTTPS,
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

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmKeyVaultKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).keyVaultManagementClient
	ctx := meta.(*ArmClient).StopContext

	vaultUri := d.Get("vault_uri").(string)
	name := d.Get("name").(string)

	resp, err := client.GetKey(ctx, vaultUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Key %q was not found in Key Vault at URI %q", name, vaultUri)
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

	flattenAndSetTags(d, resp.Tags)

	return nil
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
