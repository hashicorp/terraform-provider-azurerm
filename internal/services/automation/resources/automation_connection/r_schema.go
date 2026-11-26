// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Connection is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"type": {
			Description:  "The type of the Connection - can be either builtin type such as `Azure`, `AzureClassicCertificate`, and `AzureServicePrincipal`, or a user defined types.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"values": {
			Description: "A mapping of key value pairs passed to the connection. Different `type` needs different parameters in the `values`.",
			Type:        pluginsdk.TypeMap,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description: "A description for this Connection.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}
