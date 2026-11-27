// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_shared_private_link_resource

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubsubSharedPrivateLinkServiceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	webPubSubId, err := webpubsub.ParseWebPubSubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubSubId, err)
	}

	id := webpubsub.NewSharedPrivateLinkResourceID(subscriptionId, webPubSubId.ResourceGroupName, webPubSubId.WebPubSubName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.SharedPrivateLinkResourcesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_web_pubsub_shared_private_link_resource", id.ID())
			}
		}
	}

	parameters := webpubsub.SharedPrivateLinkResource{
		Properties: &webpubsub.SharedPrivateLinkResourceProperties{
			GroupId:               d.Get("subresource_name").(string),
			PrivateLinkResourceId: d.Get("target_resource_id").(string),
		},
	}

	requestMessage := d.Get("request_message").(string)
	if requestMessage != "" {
		parameters.Properties.RequestMessage = pointer.To(requestMessage)
	}

	if d.IsNewResource() {
		if err := client.SharedPrivateLinkResourcesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %w", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %w", id, err)
		}
	}

	return resourceWebPubsubSharedPrivateLinkServiceRead(d, meta)
}
