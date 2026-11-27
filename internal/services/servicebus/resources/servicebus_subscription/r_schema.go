// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceServicebusSubscriptionSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SubscriptionName(),
		},

		// lintignore: S013
		"topic_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},

		"auto_delete_on_idle": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "P10675199DT2H48M5.4775807S", // Never
		},

		"default_message_ttl": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "P10675199DT2H48M5.4775807S", // Unbounded
		},

		"lock_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "PT1M", // 1 minute
		},

		"dead_lettering_on_message_expiration": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"dead_lettering_on_filter_evaluation_error": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"batched_operations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"max_delivery_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"requires_session": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"forward_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"forward_dead_lettered_messages_to": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(subscriptions.EntityStatusActive),
			ValidateFunc: validation.StringInSlice([]string{
				string(subscriptions.EntityStatusActive),
				string(subscriptions.EntityStatusDisabled),
				string(subscriptions.EntityStatusReceiveDisabled),
			}, false),
		},

		"client_scoped_subscription_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_scoped_subscription": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"shareable": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  true,
					},

					"durable": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}

	if !features.SixPointOh() {
		client := s["client_scoped_subscription"].Elem.(*pluginsdk.Resource).Schema

		client["is_client_scoped_subscription_shareable"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"client_scoped_subscription.0.shareable"},
			Deprecated:    "`client_scoped_subscription.is_client_scoped_subscription_shareable` has been deprecated in favour of `client_scoped_subscription.shareable` and will be removed in v6.0 of the AzureRM Provider",
		}

		client["shareable"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"client_scoped_subscription.0.is_client_scoped_subscription_shareable"},
		}

		client["is_client_scoped_subscription_durable"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "`client_scoped_subscription.is_client_scoped_subscription_durable` has been deprecated in favour of `client_scoped_subscription.durable` and will be removed in v6.0 of the AzureRM Provider",
		}
	}

	return s
}
