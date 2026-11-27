// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_active_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Description:  "The ID of the Slot to swap with `Production`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: webapps.ValidateSlotID,
		},

		"overwrite_network_config": {
			Description: "The swap action should overwrite the Production slot's network configuration with the configuration from this slot.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_successful_swap": {
			Description: "The timestamp of the last successful swap with `Production`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
