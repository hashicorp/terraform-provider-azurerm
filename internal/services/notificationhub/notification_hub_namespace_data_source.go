package notificationhub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceNotificationHubNamespace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: resourceArmDataSourceNotificationHubNamespaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"namespace_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

			"servicebus_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDataSourceNotificationHubNamespaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	sku := flattenNotificationHubDataSourceNamespacesSku(resp.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	if props := resp.NamespaceProperties; props != nil {
		d.Set("enabled", props.Enabled)
		d.Set("namespace_type", props.NamespaceType)
		d.Set("servicebus_endpoint", props.ServiceBusEndpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenNotificationHubDataSourceNamespacesSku(input *notificationhubs.Sku) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := map[string]interface{}{
		"name": string(input.Name),
	}
	outputs = append(outputs, output)
	return outputs
}
