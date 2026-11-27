// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	keyvaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Key Vault Certificate Issuer.",
		Read:        dataSourceKeyVaultCertificateIssuerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),

			"name": {
				Description:  "The name of the Key Vault Certificate Issuer.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateIssuerName,
			},

			"provider_name": {
				Description: "The name of the third-party Certificate Issuer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"account_id": {
				Description: "The account number with the third-party Certificate Issuer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"org_id": {
				Description: "The organization ID with the third-party Certificate Issuer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"admin": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email_address": {
							Description: "E-mail address of the admin.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"first_name": {
							Description: "First name of the admin.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"last_name": {
							Description: "Last name of the admin.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"phone": {
							Description: "Phone number of the admin.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKeyVaultCertificateIssuerRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := keyvaultClient.From(meta)
	client := keyvaultClient.From(meta).ManagementClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Base URI for Certificate Issuer %q in %s: %+v", name, *keyVaultId, err)
	}

	id := parse.NewIssuerID(*keyVaultBaseUri, name)
	resp, err := client.GetCertificateIssuer(ctx, *keyVaultBaseUri, name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, *keyVaultBaseUri)
		}
		return fmt.Errorf("failed making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}
	d.SetId(id.ID())

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
