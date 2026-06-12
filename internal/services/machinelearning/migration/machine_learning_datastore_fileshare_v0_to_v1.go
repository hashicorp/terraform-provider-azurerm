// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	storageparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = MachineLearningDataStoreFileShareV0ToV1{}

type MachineLearningDataStoreFileShareV0ToV1 struct{}

func (MachineLearningDataStoreFileShareV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {

		if v, ok := rawState["storage_fileshare_id"].(string); ok && v != "" {
			if id, err := storageparse.StorageShareResourceManagerID(v); err == nil {
				newID := fileshares.NewShareID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.FileshareName).ID()
				log.Printf("[DEBUG] Updating `storage_fileshare_id` from %q to %q", v, newID)
				rawState["storage_fileshare_id"] = newID
			}
		}

		return rawState, nil
	}
}

func (MachineLearningDataStoreFileShareV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"storage_fileshare_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"service_data_identity": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"account_key": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"shared_access_signature": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"is_default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
