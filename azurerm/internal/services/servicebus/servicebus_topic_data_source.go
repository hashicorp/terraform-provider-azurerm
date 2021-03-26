package servicebus

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceServiceBusTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceBusTopicRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azValidate.TopicName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azValidate.NamespaceName,
			},

			"auto_delete_on_idle": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_message_ttl": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"duplicate_detection_history_time_window": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_batched_operations": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_express": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_partitioning": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"max_size_in_megabytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"requires_duplicate_detection": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"support_ordering": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceBusTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	id := parse.NewTopicID(subscriptionId, resourceGroup, namespaceName, name)

	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Service Bus Topic %q was not found in Resource Group %q Namespace %q", name, resourceGroup, namespaceName)
		}
		return fmt.Errorf("Error making Read request on Azure Service Bus Topic %q in Resource Group %q Namespace %q: %v", name, resourceGroup, namespaceName, err)
	}

	d.SetId(id.ID())

	if props := resp.SBTopicProperties; props != nil {
		d.Set("status", string(props.Status))
		d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
		d.Set("default_message_ttl", props.DefaultMessageTimeToLive)

		if window := props.DuplicateDetectionHistoryTimeWindow; window != nil && *window != "" {
			d.Set("duplicate_detection_history_time_window", window)
		}

		d.Set("enable_batched_operations", props.EnableBatchedOperations)
		d.Set("enable_express", props.EnableExpress)
		d.Set("enable_partitioning", props.EnablePartitioning)
		d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
		d.Set("support_ordering", props.SupportOrdering)

		if maxSizeMB := props.MaxSizeInMegabytes; maxSizeMB != nil {
			maxSize := int(*props.MaxSizeInMegabytes)

			// if the topic is in a premium namespace and partitioning is enabled then the
			// max size returned by the API will be 16 times greater than the value set
			if partitioning := props.EnablePartitioning; partitioning != nil && *partitioning {
				namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
				namespace, err := namespacesClient.Get(ctx, resourceGroup, namespaceName)
				if err != nil {
					return err
				}

				if namespace.Sku.Name != servicebus.Premium {
					const partitionCount = 16
					maxSize = int(*props.MaxSizeInMegabytes / partitionCount)
				}
			}

			d.Set("max_size_in_megabytes", maxSize)
		}
	}

	return nil
}
