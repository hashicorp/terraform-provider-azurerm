// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var appServiceCustomHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceAppServiceCustomHostnameBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCustomHostnameBindingCreate,
		Read:   resourceAppServiceCustomHostnameBindingRead,
		Delete: resourceAppServiceCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceCustomHostnameBindingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"hostname": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"app_service_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ssl_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(web.SslStateIPBasedEnabled),
					string(web.SslStateSniEnabled),
				}, false),
			},

			"thumbprint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_ip": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Hostname Binding creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	hostname := d.Get("hostname").(string)
	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)

	locks.ByName(appServiceName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(appServiceName, appServiceCustomHostnameBindingResourceName)

	if d.IsNewResource() {
		existing, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Custom Hostname Binding %q (App Service %q / Resource Group %q): %s", hostname, appServiceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_binding", *existing.ID)
		}
	}

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SiteName: utils.String(appServiceName),
		},
	}

	if sslState != "" {
		if thumbprint == "" {
			return fmt.Errorf("`thumbprint` must be specified when `ssl_state` is set")
		}

		properties.HostNameBindingProperties.SslState = web.SslState(sslState)
	}

	if thumbprint != "" {
		if sslState == "" {
			return fmt.Errorf("`ssl_state` must be specified when `thumbprint` is set")
		}

		properties.HostNameBindingProperties.Thumbprint = utils.String(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, resourceGroup, appServiceName, hostname, properties); err != nil {
		return fmt.Errorf("creating/updating Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", hostname, appServiceName, resourceGroup, err)
	}

	read, err := client.GetHostNameBinding(ctx, resourceGroup, appServiceName, hostname)
	if err != nil {
		return fmt.Errorf("retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", hostname, appServiceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Hostname Binding %q (App Service %q / Resource Group %q) ID", hostname, appServiceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceAppServiceCustomHostnameBindingRead(d, meta)
}

func resourceAppServiceCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Hostname Binding %q (App Service %q / Resource Group %q) was not found - removing from state", id.Name, id.AppServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
	}

	d.Set("hostname", id.Name)
	d.Set("app_service_name", id.AppServiceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.HostNameBindingProperties; props != nil {
		d.Set("ssl_state", props.SslState)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("virtual_ip", props.VirtualIP)
	}

	return nil
}

func resourceAppServiceCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.AppServiceName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.AppServiceName, appServiceCustomHostnameBindingResourceName)

	log.Printf("[DEBUG] Deleting App Service Hostname Binding %q (App Service %q / Resource Group %q)", id.Name, id.AppServiceName, id.ResourceGroup)

	resp, err := client.DeleteHostNameBinding(ctx, id.ResourceGroup, id.AppServiceName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Custom Hostname Binding %q (App Service %q / Resource Group %q): %+v", id.Name, id.AppServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
