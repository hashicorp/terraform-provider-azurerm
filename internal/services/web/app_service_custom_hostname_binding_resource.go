// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const appServiceCustomHostnameBindingResourceName = "azurerm_app_service_custom_hostname_binding"

func resourceAppServiceCustomHostnameBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCustomHostnameBindingCreate,
		Read:   resourceAppServiceCustomHostnameBindingRead,
		Delete: resourceAppServiceCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapps.ParseHostNameBindingID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AppServiceCustomHostnameBindingV0ToV1{},
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
					string(webapps.SslStateIPBasedEnabled),
					string(webapps.SslStateSniEnabled),
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
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := webapps.NewHostNameBindingID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("app_service_name").(string), d.Get("hostname").(string))

	locks.ByName(id.HostNameBindingName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceCustomHostnameBindingResourceName)

	existing, err := client.GetHostNameBinding(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %w", id, err)
		}
		return tf.ImportAsExistsError("azurerm_app_service_custom_hostname_binding", id.ID())
	}

	payload := webapps.HostNameBinding{
		Properties: &webapps.HostNameBindingProperties{
			SiteName: pointer.To(id.SiteName),
		},
	}

	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)
	if sslState != "" {
		if thumbprint == "" {
			return fmt.Errorf("`thumbprint` must be specified when `ssl_state` is set")
		}

		payload.Properties.SslState = pointer.ToEnum[webapps.SslState](sslState)
	}

	if thumbprint != "" {
		if sslState == "" {
			return fmt.Errorf("`ssl_state` must be specified when `thumbprint` is set")
		}

		payload.Properties.Thumbprint = pointer.To(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBinding(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCustomHostnameBindingRead(d, meta)
}

func resourceAppServiceCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapps.ParseHostNameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBinding(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("hostname", id.HostNameBindingName)
	d.Set("app_service_name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("ssl_state", props.SslState)
			d.Set("thumbprint", props.Thumbprint)
			d.Set("virtual_ip", props.VirtualIP)
		}
	}

	return nil
}

func resourceAppServiceCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapps.ParseHostNameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostNameBindingName, appServiceCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceCustomHostnameBindingResourceName)

	resp, err := client.DeleteHostNameBinding(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
