package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceKeyVaultSecret() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultSecretRead,

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

			"value": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_versionless_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceKeyVaultSecretRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, *keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Secret %q (KeyVault URI %q) does not exist", name, *keyVaultBaseUri)
		}
		return fmt.Errorf("making Read request on Azure KeyVault Secret %s: %+v", name, err)
	}

	// the version may have changed, so parse the updated id
	respID, err := parse.ParseNestedItemID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	d.Set("name", respID.Name)
	d.Set("key_vault_id", keyVaultId.ID())
	d.Set("value", resp.Value)
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)
	d.Set("versionless_id", respID.VersionlessID())

	d.Set("resource_id", parse.NewSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, respID.Name, respID.Version).ID())
	d.Set("resource_versionless_id", parse.NewSecretVersionlessID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, respID.Name).ID())

	return tags.FlattenAndSet(d, resp.Tags)
}
