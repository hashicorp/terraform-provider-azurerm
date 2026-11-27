// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_hub

import (
	"fmt"

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

func resourceWebPubSubHubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	webPubSubIdRaw := d.Get("web_pubsub_id").(string)
	webPubSubId, err := webpubsub.ParseWebPubSubID(webPubSubIdRaw)
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubSubIdRaw, err)
	}

	id := webpubsub.NewHubID(subscriptionId, webPubSubId.ResourceGroupName, webPubSubId.WebPubSubName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.HubsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_web_pubsub_hub", id.ID())
			}
		}
	}

	anonymousPolicyEnabled := "Deny"
	if d.Get("anonymous_connections_enabled").(bool) {
		anonymousPolicyEnabled = "Allow"
	}

	parameters := webpubsub.WebPubSubHub{
		Properties: webpubsub.WebPubSubHubProperties{
			EventHandlers:          expandEventHandler(d.Get("event_handler").([]any)),
			AnonymousConnectPolicy: &anonymousPolicyEnabled,
		},
	}

	eventListener, err := expandEventListener(d.Get("event_listener").([]any))
	if err != nil {
		return fmt.Errorf("expanding event listener for web pubsub %s: %+v", id, err)
	}

	parameters.Properties.EventListeners = eventListener

	if d.IsNewResource() {
		if err := client.HubsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.HubsCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceWebPubSubHubRead(d, meta)
}
