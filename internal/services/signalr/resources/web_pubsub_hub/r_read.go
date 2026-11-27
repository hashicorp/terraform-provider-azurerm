// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_hub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubSubHubRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return fmt.Errorf("%q was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HubName)
	d.Set("web_pubsub_id", webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID())

	if model := resp.Model; model != nil {
		if err := d.Set("event_handler", flattenEventHandler(model.Properties.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
		d.Set("anonymous_connections_enabled", strings.EqualFold(*model.Properties.AnonymousConnectPolicy, "Allow"))
		if err := d.Set("event_listener", flattenEventListener(model.Properties.EventListeners)); err != nil {
			return fmt.Errorf("setting `event_listener`: %+v", err)
		}
	}

	return nil
}
