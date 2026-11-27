// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_authorization_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceServiceBusNamespaceAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return helpers.AuthorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the ServiceBus Namespace Authorization Rule resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},
		// lintignore: S013
		"namespace_id": {
			Description:  "The ID of the ServiceBus Namespace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},
	})
}
