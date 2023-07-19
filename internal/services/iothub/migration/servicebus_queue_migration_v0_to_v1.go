// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IoTHubServiceBusQueueV0ToV1 struct{}

func (s IoTHubServiceBusQueueV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		//lintignore: S013
		"iothub_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"authentication_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"connection_string"},
		},

		"endpoint_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"entity_path"},
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"entity_path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"endpoint_uri"},
		},

		"connection_string": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},
	}
}

func (s IoTHubServiceBusQueueV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.EndpointServiceBusQueueIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
