// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func servicebusNamespaceDisasterRecoveryConfigSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Disaster Recovery Config. This is the alias DNS name that will be created.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"primary_namespace_id": {
			Description: "The ID of the primary Service Bus Namespace to replicate.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"partner_namespace_id": {
			Description: "The ID of the Service Bus Namespace to replicate to.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// TODO: check to see if empty values should be allowed (the previous validator permitted them)
			ValidateFunc: validation.Any(validation.StringIsEmpty, namespaces.ValidateNamespaceID),
		},

		"alias_authorization_rule_id": {
			Description:  "The Shared access policies used to access the connection string for the alias.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"primary_connection_string_alias": {
			Description: "The alias Primary Connection String for the ServiceBus Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string_alias": {
			Description: "The alias Secondary Connection String for the ServiceBus Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_primary_key": {
			Description: "The primary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_secondary_key": {
			Description: "The secondary access key for the authorization rule `RootManageSharedAccessKey`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
