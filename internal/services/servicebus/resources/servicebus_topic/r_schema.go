// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceServiceBusTopicSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Topic resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.TopicName(),
		},

		// lintignore: S013
		"namespace_id": {
			Description:  "The ID of the ServiceBus Namespace to create this topic in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"status": {
			Description: "The Status of the Service Bus Topic.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(topics.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(topics.EntityStatusActive),
				string(topics.EntityStatusDisabled),
			}, false),
		},

		"auto_delete_on_idle": {
			Description:  "The ISO 8601 timespan duration of the idle interval after which the Topic is automatically deleted, minimum of 5 minutes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P10675199DT2H48M5.4775807S", // Never
			ValidateFunc: validation.ISO8601Duration,
		},

		"default_message_ttl": {
			Description:  "The ISO 8601 timespan duration of TTL of messages sent to this topic if no TTL value is set on the message itself.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P10675199DT2H48M5.4775807S", // Unbounded
			ValidateFunc: validation.ISO8601Duration,
		},

		"duplicate_detection_history_time_window": {
			Description:  "The ISO 8601 timespan duration during which duplicates can be detected.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT10M", // 10 minutes
			ValidateFunc: validation.ISO8601Duration,
		},

		"batched_operations_enabled": {
			Description: "Boolean flag which controls if server-side batched operations are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"express_enabled": {
			Description: "Boolean flag which controls whether Express Entities are enabled. An express topic holds a message in memory temporarily before writing it to persistent storage.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"partitioning_enabled": {
			Description: "Boolean flag which controls whether to enable the topic to be partitioned across multiple message brokers.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"max_message_size_in_kilobytes": {
			Description: "Integer value which controls the maximum size of a message allowed on the topic for Premium SKU. For supported values see the \"Large messages support\" section of [this document](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-premium-messaging#large-messages-support-preview).",
			// azignore:AZP003 - named `maximum_message_size_in_kb` in the data source to follow new naming conventions
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxMessageSizeInKilobytes(),
		},

		"max_size_in_megabytes": {
			Description: "Integer value which controls the size of memory allocated for the topic. For supported values see the \"Queue/topic size\" section of [this document](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-quotas).",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxSizeInMegabytes(),
		},

		"requires_duplicate_detection": {
			Description: "Boolean flag which controls whether the Topic requires duplicate detection.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"support_ordering": {
			Description: "Boolean flag which controls whether the Topic supports ordering.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}
