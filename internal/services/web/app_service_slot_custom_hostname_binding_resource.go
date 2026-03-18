// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const appServiceSlotCustomHostnameBindingResourceName = "azurerm_app_service_slot_custom_hostname_binding"

func resourceAppServiceSlotCustomHostnameBinding() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceSlotCustomHostnameBindingCreate,
		Read:   resourceAppServiceSlotCustomHostnameBindingRead,
		Delete: resourceAppServiceSlotCustomHostnameBindingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webapps.ParseSlotHostNameBindingID(id)
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
				ValidateFunc: webapps.ValidateSlotID,
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

func resourceAppServiceSlotCustomHostnameBindingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	slotId, err := webapps.ParseSlotID(d.Get("app_service_slot_id").(string))
	if err != nil {
		return err
	}

	sslState := d.Get("ssl_state").(string)
	thumbprint := d.Get("thumbprint").(string)

	id := webapps.NewSlotHostNameBindingID(slotId.SubscriptionId, slotId.ResourceGroupName, slotId.SiteName, slotId.SlotName, d.Get("hostname").(string))

	locks.ByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)

	existing, err := client.GetHostNameBindingSlot(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_app_service_slot_custom_hostname_binding", id.ID())
	}

	payload := webapps.HostNameBinding{
		Properties: &webapps.HostNameBindingProperties{
			SiteName: pointer.To(id.SiteName),
		},
	}

	if sslState != "" {
		payload.Properties.SslState = pointer.ToEnum[webapps.SslState](sslState)
	}

	if thumbprint != "" {
		payload.Properties.Thumbprint = pointer.To(thumbprint)
	}

	if _, err := client.CreateOrUpdateHostNameBindingSlot(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceSlotCustomHostnameBindingRead(d, meta)
}

func resourceAppServiceSlotCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapps.ParseSlotHostNameBindingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHostNameBindingSlot(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("app_service_slot_id", webapps.NewSlotID(id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName).ID())
	d.Set("hostname", id.HostNameBindingName)

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("ssl_state", pointer.FromEnum(props.SslState))
			d.Set("thumbprint", props.Thumbprint)
			d.Set("virtual_ip", props.VirtualIP)
		}
	}

	return nil
}

func resourceAppServiceSlotCustomHostnameBindingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.WebAppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webapps.ParseSlotHostNameBindingID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)
	defer locks.UnlockByName(id.HostNameBindingName, appServiceSlotCustomHostnameBindingResourceName)

	if _, err := client.DeleteHostNameBindingSlot(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
