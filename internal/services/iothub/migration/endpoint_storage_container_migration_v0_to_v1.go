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

type IoTHubEndPointStorageContainerV0ToV1 struct{}

func (s IoTHubEndPointStorageContainerV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"container_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"file_name_format": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
		},

		"batch_frequency_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  300,
		},

		"max_chunk_size_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  314572800,
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
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"connection_string": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},

		"encoding": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
	}
}

func (s IoTHubEndPointStorageContainerV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.EndpointStorageContainerIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
