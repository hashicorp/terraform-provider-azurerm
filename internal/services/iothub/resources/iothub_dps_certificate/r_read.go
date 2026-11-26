// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dpscertificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dpscertificate.GetOperationOptions{IfMatch: pointer.To("")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("iot_dps_name", id.ProvisioningServiceName)
	// We are unable to set `certificate_content` since it is not returned from the API
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("is_verified", pointer.From(props.IsVerified))
		}
	}

	return nil
}
