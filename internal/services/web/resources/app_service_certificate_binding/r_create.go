// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	webClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceCertificateBindingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := webClient.From(meta).WebAppsClient
	certClient := webClient.From(meta).CertificatesClient

	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	hostnameBindingID, err := webapps.ParseHostNameBindingID(d.Get("hostname_binding_id").(string))
	if err != nil {
		return err
	}

	certificateID, err := certificates.ParseCertificateID(d.Get("certificate_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewCompositeResourceID(hostnameBindingID, certificateID)

	certificate, err := certClient.Get(ctx, *id.Second)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id.Second, err)
	}

	if certificate.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.Second)
	}

	if certificate.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.Second)
	}

	if certificate.Model.Properties.Thumbprint == nil {
		return fmt.Errorf("retrieving %s: `thumbprint` was nil", id.Second)
	}
	thumbprint := certificate.Model.Properties.Thumbprint

	binding, err := client.GetHostNameBinding(ctx, *id.First)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if binding.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}

	if binding.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	props := binding.Model.Properties
	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if props.Thumbprint != nil && *props.Thumbprint == *thumbprint {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_binding", id.ID())
		}
	}

	locks.ByName(id.First.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.First.SiteName, appServiceHostnameBindingResourceName)

	props.SslState = pointer.ToEnum[webapps.SslState](d.Get("ssl_state").(string))
	props.Thumbprint = thumbprint

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, *id.First, *binding.Model); err != nil {
		return fmt.Errorf("updating certificate for %s: %+v", id.First, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateBindingRead(d, meta)
}
