// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceServicebusQueueSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Queue resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.QueueName(),
		},

		// lintignore: S013
		"namespace_id": {
			Description:  "The ID of the ServiceBus Namespace to create this queue in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"auto_delete_on_idle": {
			Description: "The ISO 8601 timespan duration of the idle interval after which the Queue is automatically deleted, minimum of 5 minutes.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C this gets a default except when using basic sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"dead_lettering_on_message_expiration": {
			Description: "Boolean flag which controls whether the Queue has dead letter support when a message expires.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"default_message_ttl": {
			Description: "The ISO 8601 timespan duration of the TTL of messages sent to this queue. This is the default value used when TTL is not set on a message itself.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C this gets a default of "P10675199DT2H48M5.4775807S" (Unbounded) and "P14D" in Basic sku and can be updated without issues
			Computed:     true,
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
			Description: "Boolean flag which controls whether server-side batched operations are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"express_enabled": {
			Description: "Boolean flag which controls whether Express Entities are enabled. An express queue holds a message in memory temporarily before writing it to persistent storage. For Premium, it MUST be set to `false`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"partitioning_enabled": {
			Description: "Boolean flag which controls whether to enable the queue to be partitioned across multiple message brokers.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"forward_dead_lettered_messages_to": {
			Description:  "The name of a Queue or Topic to automatically forward dead lettered messages to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
		},

		"forward_to": {
			Description:  "The name of a Queue or Topic to automatically forward messages to. [See the documentation](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-auto-forwarding) for more information.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
		},

		"lock_duration": {
			Description: "The ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. Maximum value is 5 minutes.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "PT1M", // 1 minute
		},

		"max_delivery_count": {
			Description:  "Integer value which controls when a message is automatically dead lettered.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"max_message_size_in_kilobytes": {
			Description: "Integer value which controls the maximum size of a message allowed on the queue for Premium SKU. For supported values see the \"Large messages support\" section of [this document](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-premium-messaging#large-messages-support-preview).",
			// azignore:AZP003 - named `maximum_message_size_in_kb` in the data source to follow new naming conventions
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxMessageSizeInKilobytes(),
		},

		"max_size_in_megabytes": {
			Description: "Integer value which controls the size of memory allocated for the queue. For supported values see the \"Queue or topic size\" section of [Service Bus Quotas](https://learn.microsoft.com/azure/service-bus-messaging/service-bus-quotas).",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxSizeInMegabytes(),
		},

		"requires_duplicate_detection": {
			Description: "Boolean flag which controls whether the Queue requires duplicate detection.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"requires_session": {
			Description: "Boolean flag which controls whether the Queue requires sessions. This will allow ordered handling of unbounded sequences of related messages. With sessions enabled a queue can guarantee first-in-first-out delivery of messages.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"status": {
			Description: "The status of the Queue. `Restoring` is not accepted.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(queues.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(queues.EntityStatusActive),
				string(queues.EntityStatusCreating),
				string(queues.EntityStatusDeleting),
				string(queues.EntityStatusDisabled),
				string(queues.EntityStatusReceiveDisabled),
				string(queues.EntityStatusRenaming),
				string(queues.EntityStatusSendDisabled),
				string(queues.EntityStatusUnknown),
			}, false),
		},
	}
}
