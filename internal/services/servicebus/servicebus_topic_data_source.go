// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusTopic() *pluginsdk.Resource {
	d := &pluginsdk.Resource{
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

	if !features.FivePointOh() {
		d.Schema["namespace_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azValidate.NamespaceName,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
			Deprecated:   "`namespace_name` will be removed in favour of the `namespace_id` property in v5.0 of the AzureRM Provider.",
		}

		d.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: resourcegroups.ValidateName,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
			Deprecated:   "`resource_group_name` will be removed in favour of the `namespace_id` property in v5.0 of the AzureRM Provider.",
		}

		d.Schema["namespace_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
		}

		d.Schema["enable_batched_operations"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}

		d.Schema["enable_express"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}

		d.Schema["enable_partitioning"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}
	}

	return d
}

func dataSourceServiceBusTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	var id topics.TopicId // remove during `features.FivePointOh()` clean up
	if features.FivePointOh() {
		namespaceId, err := topics.ParseNamespaceID(d.Get("namespace_id").(string))
		if err != nil {
			return err
		}

		id = topics.NewTopicID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))
	} else {
		var resourceGroup string
		var namespaceName string
		if v, ok := d.Get("namespace_id").(string); ok && v != "" {
			namespaceId, err := topics.ParseNamespaceID(v)
			if err != nil {
				return err
			}
			resourceGroup = namespaceId.ResourceGroupName
			namespaceName = namespaceId.NamespaceName
		} else {
			resourceGroup = d.Get("resource_group_name").(string)
			namespaceName = d.Get("namespace_name").(string)
		}

		id = topics.NewTopicID(subscriptionId, resourceGroup, namespaceName, d.Get("name").(string))
	}

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

			if !features.FivePointOh() {
				d.Set("enable_batched_operations", props.EnableBatchedOperations)
				d.Set("enable_express", props.EnableExpress)
				d.Set("enable_partitioning", props.EnablePartitioning)
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
