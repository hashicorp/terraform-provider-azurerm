// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusQueue() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceServiceBusQueueRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azValidate.QueueName(),
			},

			"namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"auto_delete_on_idle": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dead_lettering_on_message_expiration": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"default_message_ttl": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"duplicate_detection_history_time_window": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"batched_operations_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"express_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"partitioning_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"forward_dead_lettered_messages_to": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"forward_to": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"lock_duration": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"max_delivery_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_size_in_megabytes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"requires_duplicate_detection": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"requires_session": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["namespace_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
			AtLeastOneOf: []string{"namespace_id", "namespace_name", "resource_group_name"},
		}

		resource.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: resourcegroups.ValidateName,
			AtLeastOneOf: []string{"namespace_id", "namespace_name", "resource_group_name"},
			Deprecated:   "`resource_group_name` will be removed in favour of the property `namespace_id` in version 5.0 of the AzureRM Provider.",
		}

		resource.Schema["namespace_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azValidate.NamespaceName,
			AtLeastOneOf: []string{"namespace_id", "namespace_name", "resource_group_name"},
			Deprecated:   "`namespace_name` will be removed in favour of the property `namespace_id` in version 5.0 of the AzureRM Provider.",
		}

		resource.Schema["enable_batched_operations"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}

		resource.Schema["enable_express"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}

		resource.Schema["enable_partitioning"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}
	}

	return resource
}

func dataSourceServiceBusQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if !features.FivePointOh() && resourceGroup == "" {
		resourceGroup = d.Get("resource_group_name").(string)
		namespaceName = d.Get("namespace_name").(string)
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
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("requires_session", props.RequiresSession)
			d.Set("status", string(pointer.From(props.Status)))

			if !features.FivePointOh() {
				d.Set("enable_batched_operations", props.EnableBatchedOperations)
				d.Set("enable_express", props.EnableExpress)
				d.Set("enable_partitioning", props.EnablePartitioning)
			}

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
