// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const appServiceHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceAppServiceCertificateBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateBindingCreate,
		Read:   resourceAppServiceCertificateBindingRead,
		Delete: resourceAppServiceCertificateBindingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &webapps.HostNameBindingId{}, &certificates.CertificateId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"hostname_binding_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: webapps.ValidateHostNameBindingID,
			},

			"certificate_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: certificates.ValidateCertificateID,
			},

			"ssl_state": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(webapps.SslStateIPBasedEnabled),
					string(webapps.SslStateSniEnabled),
				}, false),
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"app_service_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceCertificateBindingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient
	certClient := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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

	props := binding.Model.Properties
	if props != nil && props.Thumbprint != nil && *props.Thumbprint == *thumbprint {
		return tf.ImportAsExistsError("azurerm_app_service_certificate_binding", id.ID())
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

func resourceAppServiceCertificateBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceAppServiceCertificateBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
