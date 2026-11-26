// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoClusterCustomerManagedKeyV0ToV1 struct{}

func (s KustoClusterCustomerManagedKeyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"key_vault_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"key_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"user_identity": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s KustoClusterCustomerManagedKeyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId := rawState["id"].(string)
		newId, err := commonids.ParseKustoClusterIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
