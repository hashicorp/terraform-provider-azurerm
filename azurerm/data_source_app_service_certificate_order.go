package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAppServiceCertificateOrder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAppServiceCertificateOrderRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"auto_renew": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_vault_secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"csr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"distinguished_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"validity_in_years": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"domain_verification_token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_certificate_issuance_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"next_auto_renewal_time_stamp": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_private_key_external": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"app_service_certificate_not_renewable_reasons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"signed_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"root": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"intermediate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: appServiceCertificateDetailsSchema(),
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmAppServiceCertificateOrderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Web.CertificatesOrderClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	ctx := meta.(*ArmClient).StopContext
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
		d.Set("serial_number", props.SerialNumber)
		d.Set("is_private_key_external", props.IsPrivateKeyExternal)
		d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
		d.Set("app_service_certificate_not_renewable_reasons", utils.FlattenStringSlice(props.AppServiceCertificateNotRenewableReasons))

		if productType := props.ProductType; productType == web.StandardDomainValidatedSsl {
			d.Set("product_type", "Standard")
		} else if productType == web.StandardDomainValidatedWildCardSsl {
			d.Set("product_type", "WildCard")
		}

		if lastCertificateIssuanceTime := props.LastCertificateIssuanceTime; lastCertificateIssuanceTime != nil {
			d.Set("last_certificate_issuance_time", lastCertificateIssuanceTime.Format(time.RFC3339))
		}

		if nextAutoRenewalTimeStamp := props.NextAutoRenewalTimeStamp; nextAutoRenewalTimeStamp != nil {
			d.Set("next_auto_renewal_time_stamp", nextAutoRenewalTimeStamp.Format(time.RFC3339))
		}

		if expirationTime := props.ExpirationTime; expirationTime != nil {
			d.Set("expiration_time", expirationTime.Format(time.RFC3339))
		}

		d.Set("signed_certificate", flattenArmCertificateOrderCertificateDetails(props.SignedCertificate))
		d.Set("root", flattenArmCertificateOrderCertificateDetails(props.Root))
		d.Set("intermediate", flattenArmCertificateOrderCertificateDetails(props.Intermediate))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
