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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		// lintignore: S013
		"subscription_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     subscriptions.ValidateSubscriptions2ID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"filter_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(rules.PossibleValuesForFilterType(), false),
		},

		"action": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"sql_filter": {
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"message_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"to": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"label": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"session_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"reply_to_session_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"content_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"correlation_filter.0.correlation_id", "correlation_filter.0.message_id", "correlation_filter.0.to",
							"correlation_filter.0.reply_to", "correlation_filter.0.label", "correlation_filter.0.session_id",
							"correlation_filter.0.reply_to_session_id", "correlation_filter.0.content_type", "correlation_filter.0.properties",
						},
					},
					"properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
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
