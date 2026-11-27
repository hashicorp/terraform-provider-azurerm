// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceServicebusSubscriptionRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Subscription Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		// lintignore: S013
		"subscription_id": {
			Description:      "The ID of the ServiceBus Subscription in which this Rule should be created.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     subscriptions.ValidateSubscriptions2ID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"filter_type": {
			Description:  "Type of filter to be applied to a BrokeredMessage.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(rules.PossibleValuesForFilterType(), false),
		},

		"action": {
			Description: "Represents a set of actions written in SQL language-based syntax that is performed against a BrokeredMessage.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"sql_filter": {
			Description:  "Represents a filter written in SQL language-based syntax that to be evaluated against a BrokeredMessage.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 1024),
		},

		// Reserved for future use, currently hard-coded to 20
		"sql_filter_compatibility_level": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"correlation_filter": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"sql_filter"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"correlation_id": {
						Description: "Identifier of the correlation.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"message_id": {
						Description: "Identifier of the message.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"to": {
						Description: "Address to send to.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to": {
						Description: "Address of the queue to reply to.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"label": {
						Description: "Application specific label.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"session_id": {
						Description: "Session identifier.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to_session_id": {
						Description: "Session identifier to reply to.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"content_type": {
						Description: "Content type of the message.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"properties": {
						Description: "A list of user defined properties to be included in the filter. Specified as a map of name/value pairs.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}
