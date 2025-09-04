// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesOrderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCertificateOrderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id.Name, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

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

		if productType := props.ProductType; productType == web.CertificateProductTypeStandardDomainValidatedSsl {
			d.Set("product_type", "Standard")
		} else if productType == web.CertificateProductTypeStandardDomainValidatedWildCardSsl {
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
