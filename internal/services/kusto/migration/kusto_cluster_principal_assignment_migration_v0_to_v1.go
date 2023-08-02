// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoAttachedClusterPrincipalAssignmentV0ToV1 struct{}

func (s KustoAttachedClusterPrincipalAssignmentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"tenant_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"principal_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"principal_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"principal_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"role": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (s KustoAttachedClusterPrincipalAssignmentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.ClusterPrincipalAssignmentIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
