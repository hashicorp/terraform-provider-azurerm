// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func eventhubNamespaceDisasterRecoveryConfigSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Disaster Recovery Config.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
		},

		"namespace_name": {
			Description:  "The name of the primary EventHub Namespace to replicate.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateEventHubNamespaceName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"partner_namespace_id": {
			Description:  "The ID of the EventHub Namespace to replicate to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}
