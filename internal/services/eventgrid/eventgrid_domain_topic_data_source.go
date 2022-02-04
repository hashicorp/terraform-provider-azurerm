package eventgrid

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceEventGridDomainTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceEventGridDomainTopicRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"domain_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
		},
	}
}

func dataSourceEventGridDomainTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopicsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDomainTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("domain_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.DomainName, id.TopicName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found: %+v", id, err)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.TopicName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("domain_name", id.DomainName)

	return nil
}
