// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceBusQueueFlatten(ctx context.Context, namespacesClient *namespaces.NamespacesClient, d *pluginsdk.ResourceData, id *queues.QueueId, model *queues.SBQueue, includeResource bool) error {
	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	d.Set("name", id.QueueName)
	d.Set("namespace_id", namespaceId.ID())

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
			d.Set("dead_lettering_on_message_expiration", props.DeadLetteringOnMessageExpiration)
			d.Set("default_message_ttl", props.DefaultMessageTimeToLive)
			d.Set("duplicate_detection_history_time_window", props.DuplicateDetectionHistoryTimeWindow)
			d.Set("forward_dead_lettered_messages_to", props.ForwardDeadLetteredMessagesTo)
			d.Set("forward_to", props.ForwardTo)
			d.Set("lock_duration", props.LockDuration)
			d.Set("max_delivery_count", props.MaxDeliveryCount)
			d.Set("max_message_size_in_kilobytes", props.MaxMessageSizeInKilobytes)
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("requires_session", props.RequiresSession)
			d.Set("status", pointer.FromEnum(props.Status))

			d.Set("batched_operations_enabled", props.EnableBatchedOperations)
			d.Set("express_enabled", props.EnableExpress)
			d.Set("partitioning_enabled", props.EnablePartitioning)

			// If the queue is NOT in a premium namespace (ie. it is Basic or Standard) and partitioning is enabled
			// then the max size returned by the API will be 16 times greater than the value set.
			maxSizeInMegabytes := int(pointer.From(props.MaxSizeInMegabytes))
			if includeResource && pointer.From(props.EnablePartitioning) {
				namespace, err := namespacesClient.Get(ctx, namespaceId)
				if err != nil {
					return err
				}

				if nsModel := namespace.Model; nsModel != nil && nsModel.Sku.Name != namespaces.SkuNamePremium {
					const partitionCount = 16
					maxSizeInMegabytes = int(pointer.From(props.MaxSizeInMegabytes) / partitionCount)
				}
			}
			d.Set("max_size_in_megabytes", maxSizeInMegabytes)
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
