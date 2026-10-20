// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},

		"parent_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(connectedregistries.ValidateConnectedRegistryID, registries.ValidateRegistryID),
		},

		"sync_token_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tokens.ValidateTokenID,
		},

		"sync_schedule": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "* * * * *",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sync_message_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P1D",
			ValidateFunc: validation.ISO8601DurationBetween("P1D", "P90D"),
		},

		"sync_window": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601DurationBetween("PT3H", "P7D"),
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(connectedregistries.ConnectedRegistryModeReadWrite),
			ValidateFunc: validation.StringInSlice(
				connectedregistries.PossibleValuesForConnectedRegistryMode(),
				false,
			),
		},

		"notification": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tag": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"digest": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(parse.RepositoryNotificationActionPush),
							string(parse.RepositoryNotificationActionDelete),
							string(parse.RepositoryNotificationActionAny),
						}, false),
					},
				},
			},
		},

		"client_token_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: tokens.ValidateTokenID,
			},
		},

		"log_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  connectedregistries.LogLevelNone,
			ValidateFunc: validation.StringInSlice(
				connectedregistries.PossibleValuesForLogLevel(),
				false,
			),
		},

		"audit_log_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
