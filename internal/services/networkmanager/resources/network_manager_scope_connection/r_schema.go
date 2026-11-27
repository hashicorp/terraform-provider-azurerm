// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_scope_connection

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/scopeconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Network Manager Scope Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Description:  "The ID of the Network Manager Scope Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: scopeconnections.ValidateNetworkManagerID,
		},

		"target_scope_id": {
			Description:  "The Resource ID of the target scope which the Network Manager is connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_id": {
			Description:  "The Tenant ID of the Resource which the Network Manager is connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"description": {
			Description:  "A description of the Network Manager Scope Connection.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connection_state": {
			Description: "The connection state of the Network Manager Scope Connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
