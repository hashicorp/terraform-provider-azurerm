// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUri)
	resp, err := client.GetCertificateIssuer(ctx, certificates.NewIssuernameID(*keyVaultBaseUri, name))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist", name, *keyVaultBaseUri)
		}
		return fmt.Errorf("failed making Read request on Azure KeyVault Certificate Issuer %s: %+v", name, err)
	}
	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("provider_name", model.Provider)
		if model.OrgDetails != nil {
			if model.OrgDetails.Id != nil {
				d.Set("org_id", model.OrgDetails.Id)
			}
			d.Set("admin", flattenKeyVaultCertificateIssuerAdmins(model.OrgDetails.AdminDetails))
		}
		if model.Credentials != nil && model.Credentials.AccountId != nil {
			d.Set("account_id", model.Credentials.AccountId)
		}
	}

	return nil
}
