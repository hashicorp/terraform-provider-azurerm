// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var appServiceHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceAppServiceCertificateBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateBindingCreate,
		Read:   resourceAppServiceCertificateBindingRead,
		Delete: resourceAppServiceCertificateBindingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CertificateBindingID(id)
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
				ValidateFunc: validate.HostnameBindingID,
			},

			"certificate_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CertificateID,
			},

			"ssl_state": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.SslStateIPBasedEnabled),
					string(web.SslStateSniEnabled),
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
	client := meta.(*clients.Client).Web.AppServicesClient
	certClient := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	hostnameBindingID, err := parse.HostnameBindingID(d.Get("hostname_binding_id").(string))
	if err != nil {
		return err
	}

	certificateID, err := parse.CertificateID(d.Get("certificate_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewCertificateBindingId(*hostnameBindingID, *certificateID)
	if err != nil {
		return fmt.Errorf("could not parse ID: %+v", err)
	}

	certDetails, err := certClient.Get(ctx, id.CertificateId.ResourceGroup, id.CertificateId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(certDetails.Response) {
			return fmt.Errorf("retrieving App Service Certificate %q (Resource Group %q), not found", id.CertificateId.Name, id.CertificateId.ResourceGroup)
		}
		return fmt.Errorf("failed reading App Service Certificate %q (Resource Group %q): %+v", id.CertificateId.Name, id.CertificateId.ResourceGroup, err)
	}

	if certDetails.Thumbprint == nil {
		return fmt.Errorf("could not read thumbprint from certificate %q (resource group %q): %+v", id.CertificateId.Name, id.CertificateId.ResourceGroup, err)
	}
	thumbprint := certDetails.Thumbprint

	binding, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(binding.Response) {
			return fmt.Errorf("retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup, err)
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.SiteName, id.CertificateId.Name, id.HostnameBindingId.ResourceGroup, err)
	}

	props := binding.HostNameBindingProperties
	if props != nil {
		if props.Thumbprint != nil && *props.Thumbprint == *thumbprint {
			return tf.ImportAsExistsError("azurerm_app_service_certificate_binding", id.ID())
		}
	}

	locks.ByName(id.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceHostnameBindingResourceName)

	binding.SslState = web.SslState(d.Get("ssl_state").(string))
	binding.Thumbprint = thumbprint

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.SiteName, id.HostnameBindingId.Name, binding); err != nil {
		return fmt.Errorf("creating/updating Custom Hostname Certificate Binding %q with certificate name %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.CertificateId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateBindingRead(d, meta)
}

func resourceAppServiceCertificateBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	props := resp.HostNameBindingProperties
	if props == nil || props.Thumbprint == nil {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("hostname_binding_id", id.HostnameBindingId.ID())
	d.Set("certificate_id", id.CertificateId.ID())
	d.Set("ssl_state", string(props.SslState))
	d.Set("thumbprint", props.Thumbprint)
	d.Set("hostname", id.HostnameBindingId.Name)
	d.Set("app_service_name", id.SiteName)

	return nil
}

func resourceAppServiceCertificateBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateBindingID(d.Id())
	if err != nil {
		return err
	}

	binding, err := client.GetHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.SiteName, id.HostnameBindingId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(binding.Response) {
			log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	props := binding.HostNameBindingProperties
	if props == nil || props.Thumbprint == nil {
		log.Printf("[DEBUG] App Service Hostname Certificate Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup)
		d.SetId("")
		return nil
	}

	locks.ByName(id.SiteName, appServiceHostnameBindingResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceHostnameBindingResourceName)

	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q / Resource Group %q)", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup)

	binding.SslState = web.SslStateDisabled
	binding.Thumbprint = nil

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id.HostnameBindingId.ResourceGroup, id.SiteName, id.HostnameBindingId.Name, binding); err != nil {
		return fmt.Errorf("deleting Custom Hostname Certificate Binding %q (App Service %q / Resource Group %q): %+v", id.HostnameBindingId.Name, id.SiteName, id.HostnameBindingId.ResourceGroup, err)
	}

	return nil
}
