// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_hybrid_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"function_app_id": {
			Description:  "The ID of the Function App for this Hybrid Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
		},

		"relay_id": {
			Description:  "The ID of the Relay Hybrid Connection to use.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hybridconnections.ValidateHybridConnectionID,
		},

		"hostname": {
			Description:  "The hostname of the endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Description:  "The port to use for the endpoint.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
		},

		"send_key_name": {
			Description:  "The name of the Relay key with `Send` permission to use.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "RootManageSharedAccessKey",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_name": {
			Description: "The name of the Relay Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"relay_name": {
			Description: "The name of the Relay in use.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"service_bus_namespace": {
			Description: "The Service Bus Namespace.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"service_bus_suffix": {
			Description: "The suffix for the endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"send_key_value": {
			Description: "The Primary Access Key for the `send_key_name`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
