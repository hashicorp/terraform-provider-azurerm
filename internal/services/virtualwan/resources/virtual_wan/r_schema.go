// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_wan

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualWanSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"disable_vpn_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"allow_branch_to_branch_traffic": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"office365_local_breakout_category": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForOfficeTrafficCategory(), false),
			Default:      string(virtualwans.OfficeTrafficCategoryNone),
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Standard",
		},

		"tags": commonschema.Tags(),
	}
}
