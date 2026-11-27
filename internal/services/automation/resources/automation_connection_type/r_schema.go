// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Connection is created.",
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

		"global": {
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
					"encrypted": {
						Description: "Whether to set the isEncrypted flag of the connection field definition.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"optional": {
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

	if !features.SixPointOh() {
		s["is_global"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description:   "Whether the connection type is global.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"global"},
			Deprecated:    "`is_global` has been deprecated in favour of `global` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["global"] = &pluginsdk.Schema{
			Description:   "Whether the connection type is global.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"is_global"},
		}

		field := s["field"].Elem.(*pluginsdk.Resource).Schema

		field["is_encrypted"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description: "Whether to set the isEncrypted flag of the connection field definition.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
			Deprecated:  "`field.is_encrypted` has been deprecated in favour of `field.encrypted` and will be removed in v6.0 of the AzureRM Provider",
		}

		field["encrypted"] = &pluginsdk.Schema{
			Description: "Whether to set the isEncrypted flag of the connection field definition.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
		}

		field["is_optional"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description: "Whether to set the isOptional flag of the connection field definition.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
			Deprecated:  "`field.is_optional` has been deprecated in favour of `field.optional` and will be removed in v6.0 of the AzureRM Provider",
		}

		field["optional"] = &pluginsdk.Schema{
			Description: "Whether to set the isOptional flag of the connection field definition.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true,
		}
	}

	return s
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
