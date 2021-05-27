package eventhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/eventhubs"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceEventHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEventHubRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"partition_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceEventHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := eventhubs.NewEventhubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil && model.Properties != nil {
		d.Set("partition_count", model.Properties.PartitionCount)
		d.Set("partition_ids", model.Properties.PartitionIds)
	}

	return nil
}
