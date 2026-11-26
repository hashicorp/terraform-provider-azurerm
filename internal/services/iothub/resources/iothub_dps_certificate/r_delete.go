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

func resourceIotHubDPSCertificateDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dpscertificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dpscertificate.GetOperationOptions{IfMatch: pointer.To("")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model != nil && resp.Model.Etag == nil {
		return fmt.Errorf("deleting %s because Etag is nil", id)
	}

	deleteOptions := dpscertificate.DeleteOperationOptions{
		IfMatch:         resp.Model.Etag,
		CertificateName: pointer.To(id.CertificateName),
	}
	if _, err := client.Delete(ctx, *id, deleteOptions); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}
