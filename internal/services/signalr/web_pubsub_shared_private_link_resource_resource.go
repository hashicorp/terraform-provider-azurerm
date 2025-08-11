// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebPubSubSharedPrivateLinkService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebPubsubSharedPrivateLinkServiceCreateUpdate,
		Read:   resourceWebPubsubSharedPrivateLinkServiceRead,
		Update: resourceWebPubsubSharedPrivateLinkServiceCreateUpdate,
		Delete: resourceWebPubsubSharedPrivateLinkServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webpubsub.ParseSharedPrivateLinkResourceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(&webpubsub.WebPubSubId{}),

			"subresource_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"request_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWebPubsubSharedPrivateLinkServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubSubId, err := webpubsub.ParseWebPubSubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubSubId, err)
	}

	id := webpubsub.NewSharedPrivateLinkResourceID(subscriptionId, webPubSubId.ResourceGroupName, webPubSubId.WebPubSubName, d.Get("name").(string))

	if d.IsNewResource() {
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

	parameters := webpubsub.SharedPrivateLinkResource{
		Properties: &webpubsub.SharedPrivateLinkResourceProperties{
			GroupId:               d.Get("subresource_name").(string),
			PrivateLinkResourceId: d.Get("target_resource_id").(string),
		},
	}

	requestMessage := d.Get("request_message").(string)
	if requestMessage != "" {
		parameters.Properties.RequestMessage = utils.String(requestMessage)
	}

	if err := client.SharedPrivateLinkResourcesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceWebPubsubSharedPrivateLinkServiceRead(d, meta)
}

func resourceWebPubsubSharedPrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
			d.Set("status", string(pointer.From(props.Status)))
			d.Set("subresource_name", props.GroupId)
			d.Set("target_resource_id", props.PrivateLinkResourceId)
		}
	}

	return nil
}

func resourceWebPubsubSharedPrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
