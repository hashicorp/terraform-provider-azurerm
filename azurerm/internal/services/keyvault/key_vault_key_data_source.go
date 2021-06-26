package keyvault

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultKeyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.NestedItemName,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_size": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"key_opts": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"n": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"e": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Key %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	resp, err := client.GetKey(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Key %q was not found in Key Vault at URI %q", name, *keyVaultBaseUri)
		}

		return err
	}

	id := *resp.Key.Kid
	parsedId, err := parse.ParseNestedItemID(id)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("versionless_id", parsedId.VersionlessID())

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
