package keyvault

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultSecrets() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultSecretsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.VaultID,
			},

			"names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceKeyVaultSecretsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultId, err := parse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("fetching base vault url from id %q: %+v", *keyVaultId, err)
	}

	secretList, err := client.GetSecretsComplete(ctx, *keyVaultBaseUri, utils.Int32(25))
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure KeyVault %q: %+v", *keyVaultId, err)
	}

	d.SetId(keyVaultId.ID())

	var names []string

	if secretList.Response().Value != nil {
		for secretList.NotDone() {
			for _, v := range *secretList.Response().Value {
				name, err := parseNameFromSecretUrl(*v.ID)
				if err != nil {
					return err
				}
				names = append(names, *name)
				err = secretList.NextWithContext(ctx)
				if err != nil {
					return fmt.Errorf("listing secrets on Azure KeyVault %q: %+v", *keyVaultId, err)
				}
			}
		}
	}

	d.Set("names", names)
	d.Set("key_vault_id", keyVaultId.ID())

	return nil
}

func parseNameFromSecretUrl(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	// https://favoretti-keyvault.vault.azure.net/secrets/secret-name
	segments := strings.Split(uri.Path, "/")
	if len(segments) != 3 {
		return nil, fmt.Errorf("expected a Path in the format `/secrets/secret-name` but got %q", uri.Path)
	}
	return &segments[2], nil
}
