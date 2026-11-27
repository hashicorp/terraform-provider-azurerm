// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_order

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateOrderRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appservicecertificateorders.ParseCertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CertificateOrderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("tags", tags.Flatten(model.Tags))

		if props := model.Properties; props != nil {
			d.Set("auto_renew", props.AutoRenew)
			d.Set("csr", props.Csr)
			d.Set("distinguished_name", props.DistinguishedName)
			d.Set("key_size", props.KeySize)
			d.Set("validity_in_years", props.ValidityInYears)
			d.Set("domain_verification_token", props.DomainVerificationToken)
			d.Set("status", pointer.FromEnum(props.Status))
			d.Set("private_key_external", props.IsPrivateKeyExternal)
			if !features.SixPointOh() {
				d.Set("is_private_key_external", props.IsPrivateKeyExternal)
			}

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
