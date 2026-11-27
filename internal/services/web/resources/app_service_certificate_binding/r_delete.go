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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateBindingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient

	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &webapps.HostNameBindingId{}, &certificates.CertificateId{})
	if err != nil {
		return err
	}

	binding, err := client.GetHostNameBinding(ctx, *id.First)
	if err != nil {
		if response.WasNotFound(binding.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if binding.Model == nil || binding.Model.Properties == nil || binding.Model.Properties.Thumbprint == nil {
		d.SetId("")
		return nil
	}
	props := binding.Model.Properties

	locks.ByName(id.First.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.First.SiteName, appServiceHostnameBindingResourceName)

	props.SslState = pointer.To(webapps.SslStateDisabled)
	props.Thumbprint = nil

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, *id.First, *binding.Model); err != nil {
		return fmt.Errorf("deleting certificate from %s: %+v", id.First, err)
	}

	return nil
}
