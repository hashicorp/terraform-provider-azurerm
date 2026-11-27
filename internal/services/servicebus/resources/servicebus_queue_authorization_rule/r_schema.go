// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue_authorization_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceBusqueueAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return helpers.AuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Authorization Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},

		// lintignore: S013
		"queue_id": {
			Description:  "The ID of the ServiceBus Queue.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: queues.ValidateQueueID,
		},
	})
}
