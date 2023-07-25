// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoClusterCustomerManagedKeyV0ToV1 struct{}

func (s KustoClusterCustomerManagedKeyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"key_name": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"key_vault_id": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"key_version": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},

		"user_identity": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (s KustoClusterCustomerManagedKeyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.ClusterIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
