// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StorageContainerV1ToV2{}

type StorageContainerV1ToV2 struct{}

func (StorageContainerV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"container_access_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"default_encryption_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"encryption_scope_override_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"has_immutability_policy": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"has_legal_hold": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (StorageContainerV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldID, ok := rawState["id"].(string)
		if !ok {
			return rawState, fmt.Errorf("expected `id` to be of type string, got %T", rawState["id"])
		}

		// Some instances of this resource may already be using the resource manager ID, in which case this is a no-op
		if _, err := commonids.ParseStorageContainerID(oldID); err == nil {
			return rawState, nil
		}

		storageAccountID, err := resolveStorageAccountIDForStateUpgrade(ctx, meta, rawState, func(input string) (*commonids.StorageAccountId, error) {
			parsed, err := commonids.ParseStorageContainerID(input)
			if err != nil {
				return nil, err
			}
			id := commonids.NewStorageAccountID(parsed.SubscriptionId, parsed.ResourceGroupName, parsed.StorageAccountName)
			return &id, nil
		})
		if err != nil {
			return rawState, err
		}

		nameRaw, ok := rawState["name"]
		if !ok {
			return rawState, errors.New("expected a `name` attribute to be present")
		}

		name, ok := nameRaw.(string)
		if !ok {
			return rawState, fmt.Errorf("expected `name` to be of type string, got %T", nameRaw)
		}

		rawState["id"] = commonids.NewStorageContainerID(storageAccountID.SubscriptionId, storageAccountID.ResourceGroupName, storageAccountID.StorageAccountName, name).ID()
		rawState["storage_account_id"] = storageAccountID.ID()
		delete(rawState, "storage_account_name")

		return rawState, nil
	}
}
