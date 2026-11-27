// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Service Bus Queue.",
		Read:        dataSourceServiceBusQueueRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Service Bus Queue.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.QueueName(),
			},

			"namespace_id": {
				Description:  "The ID of the ServiceBus Namespace where the Service Bus Queue exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"auto_delete_on_idle": {
				Description: "The ISO 8601 timespan duration of the idle interval after which the Queue is automatically deleted, minimum of 5 minutes.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"dead_lettering_on_message_expiration": {
				Description: "Boolean flag which controls whether the Queue has dead letter support when a message expires.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"default_message_ttl": {
				Description: "The ISO 8601 timespan duration of the TTL of messages sent to this queue. This is the default value used when TTL is not set on a message itself.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"duplicate_detection_history_time_window": {
				Description: "The ISO 8601 timespan duration during which duplicates can be detected.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"batched_operations_enabled": {
				Description: "Boolean flag which controls whether server-side batched operations are enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"express_enabled": {
				Description: "Boolean flag which controls whether Express Entities are enabled. An express queue holds a message in memory temporarily before writing it to persistent storage.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"partitioning_enabled": {
				Description: "Boolean flag which controls whether to enable the queue to be partitioned across multiple message brokers.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"forward_dead_lettered_messages_to": {
				Description: "The name of a Queue or Topic to automatically forward dead lettered messages to.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"forward_to": {
				Description: "The name of a Queue or Topic to automatically forward messages to. [See the documentation](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-auto-forwarding) for more information.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"lock_duration": {
				Description: "The ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"max_delivery_count": {
				Description: "Integer value which controls when a message is automatically dead lettered.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"maximum_message_size_in_kb": {
				Description: "Integer value which controls the maximum size of a message allowed on the queue for Premium SKU.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"max_size_in_megabytes": {
				Description: "Integer value which controls the size of memory allocated for the queue. For supported values see the \"Queue or topic size\" section of [Service Bus Quotas](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-quotas).",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"requires_duplicate_detection": {
				Description: "Boolean flag which controls whether the Queue requires duplicate detection.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"requires_session": {
				Description: "Boolean flag which controls whether the Queue requires sessions. This will allow ordered handling of unbounded sequences of related messages. With sessions enabled a queue can guarantee first-in-first-out delivery of messages.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"status": {
				Description: "The status of the Queue.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceServiceBusQueueRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.QueuesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	var resourceGroup string
	var namespaceName string
	if v, ok := d.Get("namespace_id").(string); ok && v != "" {
		namespaceId, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
		if err != nil {
			return err
		}
		resourceGroup = namespaceId.ResourceGroupName
		namespaceName = namespaceId.NamespaceName
	}

	id := queues.NewQueueID(subscriptionId, resourceGroup, namespaceName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
			d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
			d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
			d.Set("duplicate_detection_history_time_window", props.DuplicateDetectionHistoryTimeWindow)
			d.Set("batched_operations_enabled", props.EnableBatchedOperations)
			d.Set("express_enabled", props.EnableExpress)
			d.Set("partitioning_enabled", props.EnablePartitioning)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("forward_to", props.ForwardTo)
			d.Set("lock_duration", props.LockDuration)
			d.Set("max_delivery_count", props.MaxDeliveryCount)
			d.Set("maximum_message_size_in_kb", props.MaxMessageSizeInKilobytes)
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("requires_session", props.RequiresSession)
			d.Set("status", pointer.FromEnum(props.Status))

			if apiMaxSizeInMegabytes := props.MaxSizeInMegabytes; apiMaxSizeInMegabytes != nil {
				maxSizeInMegabytes := int(*apiMaxSizeInMegabytes)

				// If the queue is NOT in a premium namespace (ie. it is Basic or Standard) and partitioning is enabled
				// then the max size returned by the API will be 16 times greater than the value set.
				if *props.EnablePartitioning {
					namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
					namespaceId := namespaces.NewNamespaceID(subscriptionId, resourceGroup, namespaceName)
					namespace, err := namespacesClient.Get(ctx, namespaceId)
					if err != nil {
						return err
					}

					if nsModel := namespace.Model; nsModel != nil && nsModel.Sku.Name != namespaces.SkuNamePremium {
						const partitionCount = 16
						maxSizeInMegabytes = int(*apiMaxSizeInMegabytes / partitionCount)
					}
				}

				d.Set("max_size_in_megabytes", maxSizeInMegabytes)
			}
		}
	}

	return nil
}
