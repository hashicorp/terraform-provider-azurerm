package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKeyVaultCertificateIssuer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKeyVaultCertificateIssuerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateKeyVaultCertificateIssuerName,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmKeyVaultCertificateIssuerRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Certificate Issuer %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	// we always want to get the latest version
	resp, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("provider_name", resp.Provider)
	d.Set("org_id", resp.OrganizationDetails.ID)
	d.Set("account_id", resp.Credentials.AccountID)
	d.Set("password", resp.Credentials.Password)
	if resp.OrganizationDetails.AdminDetails != nil {
		d.Set("admins", flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails))
	}

	return nil
}
