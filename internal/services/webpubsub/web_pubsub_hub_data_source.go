package webpubsub

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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

			"web_pubsub_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"event_handler": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"url_template": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"user_event_pattern": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"system_events": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"auth": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"managed_identity_resource": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"anonymous_connect_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWebPubsubHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Webpubsub.WebPubsubHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWebPubsubHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("web_pubsub_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.HubName, id.ResourceGroup, id.WebPubsubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.HubName)
	d.Set("web_pubsub_name", id.WebPubsubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil && props.EventHandlers != nil {

		if err := d.Set("event_handler", flattenEventHandler(props.EventHandlers)); err != nil {
			return fmt.Errorf("setting `event_handler`: %+v", err)
		}
	}

	return nil
}
