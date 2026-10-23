// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubCertificateRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("iothub_name", id.IotHubName)

	if props := resp.Properties; props != nil {
		d.Set("is_verified", props.IsVerified)
		d.Set("certificate_content", props.Certificate)
	}

	return nil
}
