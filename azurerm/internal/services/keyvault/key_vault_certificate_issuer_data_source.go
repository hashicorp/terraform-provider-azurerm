package keyvault

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceKeyVaultCertificateIssuer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyVaultCertificateIssuerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateIssuerName,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
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

func dataSourceKeyVaultCertificateIssuerRead(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUri, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("failed to look up Certificate Issuer %q, vault id %q: %+v", name, keyVaultId, err)
	}

	resp, err := client.GetCertificateIssuer(ctx, keyVaultBaseUri, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, keyVaultBaseUri)
		}
		return fmt.Errorf("failed making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("failure reading Key Vault Certificate Issuer ID for %q", name)
	}
	d.SetId(*resp.ID)

	d.Set("provider_name", resp.Provider)
	if resp.OrganizationDetails != nil {
		if resp.OrganizationDetails.ID != nil {
			d.Set("org_id", resp.OrganizationDetails.ID)
		}
		d.Set("admin", flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails))
	}
	if resp.Credentials != nil && resp.Credentials.AccountID != nil {
		d.Set("account_id", resp.Credentials.AccountID)
	}

	return nil
}
