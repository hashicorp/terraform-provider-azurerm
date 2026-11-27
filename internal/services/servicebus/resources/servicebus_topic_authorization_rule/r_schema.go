// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_topic_authorization_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceBusTopicAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Topic Authorization Rule resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},

		// lintignore: S013
		"topic_id": {
			Description:  "The ID of the ServiceBus Topic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},
	}
}
