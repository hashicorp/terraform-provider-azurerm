// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_certificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubCertificateDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.IotHubCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotHubCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Etag == nil {
		return fmt.Errorf("deleting  %s because Etag is nil", *id)
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IotHubName, id.CertificateName, *pointer.To(*resp.Etag)); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}
