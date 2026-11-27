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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.QueueName(),
		},

		// lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"auto_delete_on_idle": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C this gets a default except when using basic sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"default_message_ttl": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C this gets a default of "P10675199DT2H48M5.4775807S" (Unbounded) and "P14D" in Basic sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"duplicate_detection_history_time_window": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT10M", // 10 minutes
			ValidateFunc: validation.ISO8601Duration,
		},

		"batched_operations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"express_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"partitioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"forward_dead_lettered_messages_to": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
		},

		"forward_to": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
		},

		"lock_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "PT1M", // 1 minute
		},

		"max_delivery_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"max_message_size_in_kilobytes": { // azignore:AZP003 - named `maximum_message_size_in_kb` in the data source to follow new naming conventions
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxMessageSizeInKilobytes(),
		},

		"max_size_in_megabytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C this gets a variable default based on the sku and can be updated without issues
			Computed:     true,
			ValidateFunc: validate.ServiceBusMaxSizeInMegabytes(),
		},

		"requires_duplicate_detection": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"requires_session": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(queues.EntityStatusActive),
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
