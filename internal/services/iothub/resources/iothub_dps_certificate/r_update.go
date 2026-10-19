// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSCertificateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dpscertificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iot_dps_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, dpscertificate.GetOperationOptions{IfMatch: pointer.To("")})
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	etag := ""
	if existing.Model != nil && existing.Model.Etag != nil {
		etag = *existing.Model.Etag
	}

	if d.HasChange("is_verified") {
		existing.Model.Properties.IsVerified = pointer.To(d.Get("is_verified").(bool))
	}

	if d.HasChange("certificate_content") {
		existing.Model.Properties.Certificate = pointer.To(d.Get("certificate_content").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, *existing.Model, dpscertificate.CreateOrUpdateOperationOptions{IfMatch: pointer.To(etag)}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceIotHubDPSCertificateRead(d, meta)
}
