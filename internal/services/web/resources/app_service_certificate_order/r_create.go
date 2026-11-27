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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateOrderCreate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := appservicecertificateorders.NewCertificateOrderID(account.From(meta).SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_order", id.ID())
		}
	}

	certificateOrder := appservicecertificateorders.AppServiceCertificateOrder{
		Properties: &appservicecertificateorders.AppServiceCertificateOrderProperties{
			AutoRenew:         pointer.To(d.Get("auto_renew").(bool)),
			Csr:               pointer.To(d.Get("csr").(string)),
			DistinguishedName: pointer.To(d.Get("distinguished_name").(string)),
			KeySize:           pointer.To(int64(d.Get("key_size").(int))),
			ProductType:       expandProductType(d.Get("product_type").(string)),
			ValidityInYears:   pointer.To(int64(d.Get("validity_in_years").(int))),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, certificateOrder, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %w", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateOrderRead(d, meta)
}

func expandProductType(input string) appservicecertificateorders.CertificateProductType {
	if input == "WildCard" {
		return appservicecertificateorders.CertificateProductTypeStandardDomainValidatedWildCardSsl
	}
	return appservicecertificateorders.CertificateProductTypeStandardDomainValidatedSsl
}
