// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceServicebusSubscriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Subscription resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SubscriptionName(),
		},

		// lintignore: S013
		"topic_id": {
			Description:  "The ID of the ServiceBus Topic to create this Subscription in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},

		"auto_delete_on_idle": {
			Description: "The idle interval after which the subscription is automatically deleted as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The minimum duration is `5` minutes or `PT5M`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "P10675199DT2H48M5.4775807S", // Never
		},

		"default_message_ttl": {
			Description: "This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the value used when TimeToLive is not set on a message itself.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "P10675199DT2H48M5.4775807S", // Unbounded
		},

		"lock_duration": {
			Description: "The lock duration for the subscription as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "PT1M", // 1 minute
		},

		"dead_lettering_on_message_expiration": {
			Description: "Boolean flag which controls whether the Subscription has dead letter support when a message expires.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"dead_lettering_on_filter_evaluation_error": {
			Description: "Boolean flag which controls whether the Subscription has dead letter support on filter evaluation exceptions.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"batched_operations_enabled": {
			Description: "Boolean flag which controls whether the Subscription supports batched operations.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"max_delivery_count": {
			Description: "The maximum number of deliveries.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"requires_session": {
			Description: "Boolean flag which controls whether this Subscription supports the concept of a session.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"forward_to": {
			Description: "The name of a Queue or Topic to automatically forward messages to.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"forward_dead_lettered_messages_to": {
			Description: "The name of a Queue or Topic to automatically forward Dead Letter messages to.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"status": {
			Description: "The status of the Subscription.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(subscriptions.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(subscriptions.EntityStatusActive),
				string(subscriptions.EntityStatusDisabled),
				string(subscriptions.EntityStatusReceiveDisabled),
			}, false),
		},

		"client_scoped_subscription_enabled": {
			Description: "Whether the subscription is scoped to a client ID.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_scoped_subscription": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description: "The Client ID of the application that created the client-scoped subscription.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ForceNew:    true,
					},

					"is_client_scoped_subscription_shareable": {
						Description: "Whether the client scoped subscription is shareable.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     true,
					},

					"is_client_scoped_subscription_durable": {
						Description: "Whether the client scoped subscription is durable. This property can only be controlled from the application side.",
						Type:        pluginsdk.TypeBool,
						Computed:    true,
					},
				},
			},
		},
	}
}
