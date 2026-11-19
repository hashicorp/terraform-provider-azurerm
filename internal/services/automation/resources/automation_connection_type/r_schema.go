// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the automation account in which the Connection is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"name": {
			Description:  "The name which should be used for this Automation Connection Type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConnectionTypeName,
		},

		"is_global": {
			Description: "Whether the connection type is global.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"field": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name which should be used for this connection field definition.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"is_encrypted": {
						Description: "Whether to set the isEncrypted flag of the connection field definition.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"is_optional": {
						Description: "Whether to set the isOptional flag of the connection field definition.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"type": {
						Description: "The type of the connection field definition.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
