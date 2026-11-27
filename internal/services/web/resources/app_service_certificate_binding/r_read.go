// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateBindingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &webapps.HostNameBindingId{}, &certificates.CertificateId{})
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, *id.First)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Thumbprint == nil {
		d.SetId("")
		return nil
	}
	props := resp.Model.Properties

	d.Set("hostname_binding_id", id.First.ID())
	d.Set("certificate_id", id.Second.ID())
	d.Set("ssl_state", pointer.FromEnum(props.SslState))
	d.Set("thumbprint", props.Thumbprint)
	d.Set("hostname", id.First.HostNameBindingName)
	d.Set("app_service_name", id.First.SiteName)

	return nil
}
