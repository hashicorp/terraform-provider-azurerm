package signalr

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceWebpubsubSharedPrivateLinkService() *pluginsdk.Resource {
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
			_, err := parse.WebPubsubSharedPrivateLinkResourceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"web_pubsub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WebPubsubID,
			},

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
	client := meta.(*clients.Client).SignalR.WebPubsubSharedPrivateLinkResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubsubID, err := parse.WebPubsubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubsubID, err)
	}

	id := parse.NewWebPubsubSharedPrivateLinkResourceID(subscriptionId, webPubsubID.ResourceGroup, webPubsubID.WebPubSubName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.SharedPrivateLinkResourceName, id.ResourceGroup, id.WebPubSubName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_web_pubsub_shared_private_link_resource", id.ID())
		}
	}

	parameters := webpubsub.SharedPrivateLinkResource{
		SharedPrivateLinkResourceProperties: &webpubsub.SharedPrivateLinkResourceProperties{
			GroupID:               utils.String(d.Get("subresource_name").(string)),
			PrivateLinkResourceID: utils.String(d.Get("target_resource_id").(string)),
		},
	}

	requestMessage := d.Get("request_message").(string)
	if requestMessage != "" {
		parameters.SharedPrivateLinkResourceProperties.RequestMessage = utils.String(requestMessage)
	}

	if _, err := client.CreateOrUpdate(ctx, id.SharedPrivateLinkResourceName, parameters, id.ResourceGroup, id.WebPubSubName); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceWebPubsubSharedPrivateLinkServiceRead(d, meta)
}
func resourceWebPubsubSharedPrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubSharedPrivateLinkResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.SharedPrivateLinkResourceName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return fmt.Errorf("%q was not found", id)
		}
		return fmt.Errorf("making request on %q: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("web_pubsub_id", parse.NewWebPubsubID(id.SubscriptionId, id.ResourceGroup, id.WebPubSubName).ID())

	if props := resp.SharedPrivateLinkResourceProperties; props != nil {
		if props.GroupID != nil {
			d.Set("subresource_name", props.GroupID)
		}

		if props.PrivateLinkResourceID != nil {
			d.Set("target_resource_id", props.PrivateLinkResourceID)
		}

		if props.RequestMessage != nil {
			d.Set("request_message", props.RequestMessage)
		}

		d.Set("status", props.Status)
	}

	return nil
}
func resourceWebPubsubSharedPrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubSharedPrivateLinkResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubSharedPrivateLinkResourceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.SharedPrivateLinkResourceName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting %q: %+v", id, err)
		}
	}
	return nil
}
