// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceServiceBusTopicRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azValidate.TopicName(),
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

			"max_size_in_megabytes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"requires_duplicate_detection": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"support_ordering": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceBusTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceId, err := topics.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	id := topics.NewTopicID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))

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
			status := ""
			if v := props.Status; v != nil {
				status = string(*v)
			}
			d.Set("status", status)
			d.Set("auto_delete_on_idle", props.AutoDeleteOnIdle)
			d.Set("default_message_ttl", props.DefaultMessageTimeToLive)

			if window := props.DuplicateDetectionHistoryTimeWindow; window != nil && *window != "" {
				d.Set("duplicate_detection_history_time_window", window)
			}

			d.Set("batched_operations_enabled", props.EnableBatchedOperations)
			d.Set("express_enabled", props.EnableExpress)
			d.Set("partitioning_enabled", props.EnablePartitioning)
			d.Set("requires_duplicate_detection", props.RequiresDuplicateDetection)
			d.Set("support_ordering", props.SupportOrdering)

			if maxSizeMB := props.MaxSizeInMegabytes; maxSizeMB != nil {
				maxSize := int(*props.MaxSizeInMegabytes)

				// if the topic is in a premium namespace and partitioning is enabled then the
				// max size returned by the API will be 16 times greater than the value set
				if partitioning := props.EnablePartitioning; partitioning != nil && *partitioning {
					namespacesClient := meta.(*clients.Client).ServiceBus.NamespacesClient
					namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
					namespace, err := namespacesClient.Get(ctx, namespaceId)
					if err != nil {
						return err
					}

					if namespaceModel := namespace.Model; namespaceModel != nil {
						if namespaceModel.Sku.Name != namespaces.SkuNamePremium {
							const partitionCount = 16
							maxSize = int(*props.MaxSizeInMegabytes / partitionCount)
						}
					}
				}

				d.Set("max_size_in_megabytes", maxSize)
			}
		}
	}

	return nil
}
