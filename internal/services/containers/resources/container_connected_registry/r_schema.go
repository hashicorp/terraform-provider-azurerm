// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name which should be used for this Container Connected Registry.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"container_registry_id": {
			Description:  "The ID of the Container Registry that this Connected Registry will reside in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},

		"parent_registry_id": {
			Description:  "The ID of the parent registry. This can be either a Container Registry ID or a Connected Registry ID.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(connectedregistries.ValidateConnectedRegistryID, registries.ValidateRegistryID),
		},

		"sync_token_id": {
			Description:  "The ID of the Container Registry Token which is used for synchronizing the Connected Registry.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tokens.ValidateTokenID,
		},

		"sync_schedule": {
			Description:  "The cron expression indicating the schedule that the Connected Registry will sync with its parent.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "* * * * *",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sync_message_ttl": {
			Description:  "The period of time (in form of ISO8601) for which a message is available to sync before it is expired.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "P1D",
			ValidateFunc: validation.ISO8601DurationBetween("P1D", "P90D"),
		},

		"sync_window": {
			Description:  "The time window (in form of ISO8601) during which sync is enabled for each schedule occurrence.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.ISO8601DurationBetween("PT3H", "P7D"),
		},

		"mode": {
			Description: "The mode of the Connected Registry.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(connectedregistries.ConnectedRegistryModeReadWrite),
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
						Description:  "The name of the artifact that wants to be subscribed for the Connected Registry.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tag": {
						Description:  "The tag of the artifact that wants to be subscribed for the Connected Registry.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"digest": {
						Description:  "The digest of the artifact that wants to be subscribed for the Connected Registry.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"action": {
						Description: "The action of the artifact that wants to be subscribed for the Connected Registry.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
			Description: "A list of IDs of Container Registry Tokens, which are meant to be used by the clients to connect to the Connected Registry.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: tokens.ValidateTokenID,
			},
		},

		"log_level": {
			Description: "The verbosity of the logs.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     connectedregistries.LogLevelNone,
			ValidateFunc: validation.StringInSlice(
				connectedregistries.PossibleValuesForLogLevel(),
				false,
			),
		},

		"audit_log_enabled": {
			Description: "Whether the log auditing should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
