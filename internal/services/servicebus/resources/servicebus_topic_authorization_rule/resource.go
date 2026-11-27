// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a ServiceBus Topic authorization Rule within a ServiceBus Topic.",
		Create:      resourceServiceBusTopicAuthorizationRuleCreateUpdate,
		Read:        resourceServiceBusTopicAuthorizationRuleRead,
		Update:      resourceServiceBusTopicAuthorizationRuleCreateUpdate,
		Delete:      resourceServiceBusTopicAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := topics.ParseTopicAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: helpers.AuthorizationRuleSchemaFrom(resourceServiceBusTopicAuthorizationRuleSchema()),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(helpers.AuthorizationRuleCustomizeDiff),
	}
}
