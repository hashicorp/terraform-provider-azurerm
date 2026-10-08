// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DataProtectionBackupVaultV0ToV1{}

type DataProtectionBackupVaultV0ToV1 struct{}

func (DataProtectionBackupVaultV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"datastore_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"redundancy": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"cross_region_restore_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"retention_duration_in_days": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},
		"soft_delete": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"immutability": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}

func (DataProtectionBackupVaultV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		datastoreType, ok := rawState["datastore_type"].(string)
		if !ok {
			return nil, fmt.Errorf("`datastore_type` is missing or is not a string")
		}
		redundancy, ok := rawState["redundancy"].(string)
		if !ok {
			return nil, fmt.Errorf("`redundancy` is missing or is not a string")
		}

		settings := []any{
			map[string]any{
				"datastore_type": datastoreType,
				"redundancy":     redundancy,
			},
		}
		// The 5.x provider also configured VaultStore for ArchiveStore vaults.
		if strings.EqualFold(datastoreType, string(backupvaultresources.StorageSettingStoreTypesArchiveStore)) {
			settings = append(settings, map[string]any{
				"datastore_type": string(backupvaultresources.StorageSettingStoreTypesVaultStore),
				"redundancy":     redundancy,
			})
		}

		rawState["storage_setting"] = settings
		delete(rawState, "datastore_type")
		delete(rawState, "redundancy")
		return rawState, nil
	}
}
