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

var appServiceSlotCustomHostnameBindingResourceName = "azurerm_app_service_slot_custom_hostname_binding"

func resourceAppServiceSlotCustomHostnameBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotCustomHostnameBindingCreate,
		Read:   resourceAppServiceSlotCustomHostnameBindingRead,
		Delete: resourceAppServiceSlotCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceSlotCustomHostnameBindingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"app_service_slot_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceSlotID,
			},

			"hostname": {
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

func resourceAppServiceSlotCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Slot Hostname Binding creation.")

	slotId, err := parse.AppServiceSlotID(d.Get("app_service_slot_id").(string))
	if err != nil {
		return err
	}

	hostname := d.Get("hostname").(string)
	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)

	id := parse.NewAppServiceSlotCustomHostnameBindingID(slotId.SubscriptionId, slotId.ResourceGroup, slotId.SiteName, slotId.SlotName, hostname)

	locks.ByName(hostname, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(hostname, appServiceSlotCustomHostnameBindingResourceName)

	existing, err := client.GetHostNameBindingSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, id.HostNameBindingName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_app_service_slot_custom_hostname_binding", id.ID())
	}

	properties := web.HostNameBinding{
		HostNameBindingProperties: &web.HostNameBindingProperties{
			SiteName: utils.String(id.SiteName),
		},
	}

	if sslState != "" {
		properties.HostNameBindingProperties.SslState = web.SslState(sslState)
	}

	if thumbprint != "" {
		properties.HostNameBindingProperties.Thumbprint = utils.String(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBindingSlot(ctx, id.ResourceGroup, id.SiteName, id.HostNameBindingName, properties, id.SlotName); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceSlotCustomHostnameBindingRead(d, meta)
}

func resourceAppServiceSlotCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceSlotCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBindingSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, id.HostNameBindingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	slotId := parse.NewAppServiceSlotID(id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName)
	d.Set("app_service_slot_id", slotId.ID())
	d.Set("hostname", id.HostNameBindingName)

	if props := resp.HostNameBindingProperties; props != nil {
		d.Set("ssl_state", string(props.SslState))
		d.Set("thumbprint", props.Thumbprint)
		d.Set("virtual_ip", props.VirtualIP)
	}

	return nil
}

func resourceAppServiceSlotCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceSlotCustomHostnameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)

	log.Printf("[DEBUG] deleting %s", id)

	resp, err := client.DeleteHostNameBindingSlot(ctx, id.ResourceGroup, id.SiteName, id.SlotName, id.HostNameBindingName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
