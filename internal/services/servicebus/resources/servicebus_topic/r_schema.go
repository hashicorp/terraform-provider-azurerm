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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.TopicName(),
		},

		// lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(topics.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(topics.EntityStatusActive),
				string(topics.EntityStatusDisabled),
			}, false),
		},

		"auto_delete_on_idle": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P10675199DT2H48M5.4775807S", // Never
			ValidateFunc: validation.ISO8601Duration,
		},

		"default_message_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P10675199DT2H48M5.4775807S", // Unbounded
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
		},

		"express_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"partitioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
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
		},

		"support_ordering": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}
