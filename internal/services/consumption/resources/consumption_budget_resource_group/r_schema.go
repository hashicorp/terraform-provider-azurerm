// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_resource_group

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Resource Group Consumption Budget.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"resource_group_id": {
			Description:  "The ID of the Resource Group to create the consumption budget for in the form of /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ResourceGroupID,
		},
	}
	return r.base.Arguments(schema)
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.Attributes()
}
