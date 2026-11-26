// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_module

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func automationModuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Module.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Module is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"module_link": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"uri": {
						Description: "The URI of the module content (zip or nupkg).",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"hash": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"algorithm": {
									Description: "The algorithm used for the hash content.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"value": {
									Description: "The hash value of the content.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}
