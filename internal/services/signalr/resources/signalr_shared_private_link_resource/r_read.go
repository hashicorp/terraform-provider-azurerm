// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_shared_private_link_resource

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalRSharedPrivateLinkRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SharedPrivateLinkResourcesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving shared private link %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("signalr_service_id", signalr.NewSignalRID(id.SubscriptionId, id.ResourceGroupName, id.SignalRName).ID())

		if props := model.Properties; props != nil {
			d.Set("sub_resource_name", props.GroupId)
			d.Set("target_resource_id", props.PrivateLinkResourceId)

			if props.RequestMessage != nil {
				d.Set("request_message", props.RequestMessage)
			}

			d.Set("status", pointer.FromEnum(props.Status))
		}
	}
	return nil
}
