// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_hybrid_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"web_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateWebAppID,
			Description:  "The ID of the Web App for this Hybrid Connection.",
		},

		"relay_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hybridconnections.ValidateHybridConnectionID,
			Description:  "The ID of the Relay Hybrid Connection to use.",
		},

		"hostname": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The hostname of the endpoint.",
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IsPortNumberOrZero,
			Description:  "The port to use for the endpoint",
		},

		"send_key_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "RootManageSharedAccessKey",
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The name of the Relay key with `Send` permission to use. Defaults to `RootManageSharedAccessKey`",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Relay Namespace.",
		},

		"relay_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the Relay in use.",
		},

		"service_bus_namespace": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Service Bus Namespace.",
		},

		"service_bus_suffix": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The suffix for the endpoint.",
		},

		"send_key_value": {
			Type:        pluginsdk.TypeString,
			Sensitive:   true,
			Computed:    true,
			Description: "The Primary Access Key for the `send_key_name`",
		},
	}
}
