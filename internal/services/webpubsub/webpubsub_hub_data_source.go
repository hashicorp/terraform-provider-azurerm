package webpubsub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/webpubsub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceWebPubsubHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceWebPubsubHubRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebPubsubHubID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"web_pubsub_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceWebPubsubHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	webPubsubID, err := parse.WebPubsubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing ID of %q: %+v", webPubsubID, err)
	}
	id := parse.NewWebPubsubHubID(subscriptionId, webPubsubID.ResourceGroup, webPubsubID.WebPubSubName, d.Get("name").(string))

	resp, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%q was not found", id)
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HubName)
	d.Set("web_pubsub_id", parse.NewWebPubsubID(id.SubscriptionId, id.ResourceGroup, id.WebPubSubName).ID())

	return nil
}
