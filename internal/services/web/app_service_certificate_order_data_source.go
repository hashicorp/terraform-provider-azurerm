// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := appservicecertificateorders.NewCertificateOrderID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.CertificateOrderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", model.Tags)

		if props := model.Properties; props != nil {
			d.Set("auto_renew", props.AutoRenew)
			d.Set("csr", props.Csr)
			d.Set("distinguished_name", props.DistinguishedName)
			d.Set("key_size", props.KeySize)
			d.Set("validity_in_years", props.ValidityInYears)
			d.Set("domain_verification_token", props.DomainVerificationToken)
			d.Set("status", pointer.FromEnum(props.Status))
			d.Set("is_private_key_external", props.IsPrivateKeyExternal)
			d.Set("certificates", flattenArmCertificateOrderCertificate(props.Certificates))
			d.Set("app_service_certificate_not_renewable_reasons", flattenAppServiceCertificateNotRenewableReasons(props.AppServiceCertificateNotRenewableReasons))
			d.Set("expiration_time", props.ExpirationTime)
			d.Set("product_type", flattenProductType(props.ProductType))

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
	}

	return nil
}
