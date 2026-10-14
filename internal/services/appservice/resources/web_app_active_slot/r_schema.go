// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"slot_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the Slot to swap with `Production`.",
			ValidateFunc: webapps.ValidateSlotID,
		},

		"overwrite_network_config": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "The swap action should overwrite the Production slot's network configuration with the configuration from this slot. Defaults to `true`.",
			ForceNew:    true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_successful_swap": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The timestamp of the last successful swap with `Production`",
		},
	}
}
