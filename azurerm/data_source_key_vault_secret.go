package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVaultSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultSecretRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
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

			"value": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmKeyVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*ArmClient).keyvault.VaultsClient
	client := meta.(*ArmClient).keyvault.ManagementClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	keyVaultBaseUri := d.Get("vault_uri").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	if keyVaultBaseUri == "" {
		if keyVaultId == "" {
			return fmt.Errorf("one of `key_vault_id` or `vault_uri` must be set")
		}

		pKeyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
		if err != nil {
			return fmt.Errorf("Error looking up Secret %q vault url form id %q: %+v", name, keyVaultId, err)
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

	// we always want to get the latest version
	resp, err := client.GetSecret(ctx, keyVaultBaseUri, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Secret %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Secret %s: %+v", name, err)
	}

	// the version may have changed, so parse the updated id
	respID, err := azure.ParseKeyVaultChildID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	d.Set("name", respID.Name)
	d.Set("vault_uri", respID.KeyVaultBaseUrl)
	d.Set("value", resp.Value)
	d.Set("version", respID.Version)
	d.Set("content_type", resp.ContentType)

	return tags.FlattenAndSet(d, resp.Tags)
}
