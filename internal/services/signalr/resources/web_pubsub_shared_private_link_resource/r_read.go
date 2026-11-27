// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_shared_private_link_resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubsubSharedPrivateLinkServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := webpubsub.ParseSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SharedPrivateLinkResourcesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return fmt.Errorf("%q was not found", id)
		}
		return fmt.Errorf("making request on %q: %+v", id, err)
	}

	d.Set("name", id.SharedPrivateLinkResourceName)
	d.Set("web_pubsub_id", webpubsub.NewWebPubSubID(id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("request_message", props.RequestMessage)
			d.Set("status", pointer.FromEnum(props.Status))
			d.Set("subresource_name", props.GroupId)
			d.Set("target_resource_id", props.PrivateLinkResourceId)
		}
	}

	return nil
}
