// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_shared_private_link_resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubsubSharedPrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webpubsub.ParseSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.SharedPrivateLinkResourcesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
