// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StorageQueueV1ToV2{}

type StorageQueueV1ToV2 struct{}

func (StorageQueueV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (StorageQueueV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldID, ok := rawState["id"].(string)
		if !ok {
			return rawState, fmt.Errorf("expected `id` to be of type string, got %T", rawState["id"])
		}

		// Some instances of this resource may already be using the resource manager ID, in which case this is a no-op
		if _, err := storagequeues.ParseQueueID(oldID); err == nil {
			return rawState, nil
		}

		storageAccountID, err := resolveStorageAccountIDForStateUpgrade(ctx, meta, rawState, func(idstr string) (*commonids.StorageAccountId, error) {
			id, err := parse.StorageQueueResourceManagerID(idstr)
			if err != nil {
				return nil, err
			}
			return new(commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroup, id.StorageAccountName)), nil
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

		rawState["id"] = storagequeues.NewQueueID(storageAccountID.SubscriptionId, storageAccountID.ResourceGroupName, storageAccountID.StorageAccountName, name).ID()
		rawState["storage_account_id"] = storageAccountID.ID()
		delete(rawState, "storage_account_name")

		return rawState, nil
	}
}
