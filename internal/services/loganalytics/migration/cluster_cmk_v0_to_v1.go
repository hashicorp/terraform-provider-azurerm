// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ClusterCustomerManagedKeyV0ToV1 struct{}

func (c ClusterCustomerManagedKeyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"log_analytics_cluster_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"key_vault_key_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (c ClusterCustomerManagedKeyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		clusterId := rawState["log_analytics_cluster_id"].(string)

		// @tombuildsstuff: we could re-parse the `id` and trim off `/CMK`, however since the ID is `{ClusterID}/CMK`
		// we can instead save the hassle and just use the Cluster ID directly
		log.Printf("[DEBUG] Updating the ID from %q to %q..", oldId, clusterId)
		rawState["id"] = clusterId
		log.Printf("[DEBUG] Updated the ID from %q to %q.", oldId, clusterId)

		return rawState, nil
	}
}
