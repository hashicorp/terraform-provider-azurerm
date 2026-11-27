// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_order

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateOrderUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.AppServiceCertificateOrdersClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appservicecertificateorders.ParseCertificateOrderID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("auto_renew") {
		props.AutoRenew = pointer.To(d.Get("auto_renew").(bool))
	}

	if d.HasChange("csr") {
		props.Csr = pointer.To(d.Get("csr").(string))
	}

	if d.HasChange("distinguished_name") {
		props.DistinguishedName = pointer.To(d.Get("distinguished_name").(string))
	}

	if d.HasChange("key_size") {
		props.KeySize = pointer.To(int64(d.Get("key_size").(int)))
	}

	if d.HasChange("product_type") {
		props.ProductType = expandProductType(d.Get("product_type").(string))
	}

	if d.HasChange("validity_in_years") {
		props.ValidityInYears = pointer.To(int64(d.Get("validity_in_years").(int)))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %w", id, err)
	}

	return resourceAppServiceCertificateOrderRead(d, meta)
}
