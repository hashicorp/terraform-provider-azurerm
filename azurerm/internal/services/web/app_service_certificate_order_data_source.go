package web

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAppServiceCertificateOrder() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAppServiceCertificateOrderRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"auto_renew": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"certificates": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_secret_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"provisioning_state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"csr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"distinguished_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_size": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"product_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validity_in_years": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"domain_verification_token": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expiration_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"is_private_key_external": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"app_service_certificate_not_renewable_reasons": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"signed_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"root_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"intermediate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: App Service Certificate Order %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM App Service Certificate Order %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AppServiceCertificateOrderProperties; props != nil {
		d.Set("auto_renew", props.AutoRenew)
		d.Set("csr", props.Csr)
		d.Set("distinguished_name", props.DistinguishedName)
		d.Set("key_size", props.KeySize)
		d.Set("validity_in_years", props.ValidityInYears)
		d.Set("domain_verification_token", props.DomainVerificationToken)
		d.Set("status", string(props.Status))
		d.Set("is_private_key_external", props.IsPrivateKeyExternal)
		d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
		d.Set("app_service_certificate_not_renewable_reasons", utils.FlattenStringSlice(props.AppServiceCertificateNotRenewableReasons))

		if productType := props.ProductType; productType == web.StandardDomainValidatedSsl {
			d.Set("product_type", "Standard")
		} else if productType == web.StandardDomainValidatedWildCardSsl {
			d.Set("product_type", "WildCard")
		}

		if expirationTime := props.ExpirationTime; expirationTime != nil {
			d.Set("expiration_time", expirationTime.Format(time.RFC3339))
		}

		if signedCertificate := props.SignedCertificate; signedCertificate != nil {
			d.Set("signed_certificate_thumbprint", signedCertificate.Thumbprint)
		}

		if root := props.Root; root != nil {
			d.Set("root_thumbprint", root.Thumbprint)
		}

		if intermediate := props.Intermediate; intermediate != nil {
			d.Set("intermediate_thumbprint", intermediate.Thumbprint)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
