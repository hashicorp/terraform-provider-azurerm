// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_custom_hostname_binding

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAppServiceSlotCustomHostnameBindingRead(d *pluginsdk.ResourceData, meta any) error {
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
