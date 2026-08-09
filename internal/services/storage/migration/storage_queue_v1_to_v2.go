// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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

		findAccount := false
		resourceManagerID, ok := rawState["resource_manager_id"].(string)
		if !ok {
			findAccount = true
		}

		storageAccountID := commonids.StorageAccountId{}
		if !findAccount {
			// The `resource_manager_id` can be malformed (see: #32950), in which case fallbacks to find account.
			if parsed, err := storagequeues.ParseQueueIDInsensitively(resourceManagerID); err != nil {
				findAccount = true
			} else {
				storageAccountID = commonids.NewStorageAccountID(parsed.SubscriptionId, parsed.ResourceGroupName, parsed.StorageAccountName)
			}
		}

		if findAccount {
			// `resource_manager_id` was introduced in v3.39.0
			// The find account logic is only here for edge cases where users upgrade from < v3.39.0 directly to >= 5.0.0
			client := meta.(*clients.Client).Storage
			subscriptionID := meta.(*clients.Client).Account.SubscriptionId

			storageAccountNameRaw, ok := rawState["storage_account_name"]
			if !ok {
				return rawState, errors.New("expected a `storage_account_name` attribute to be present in state")
			}

			storageAccountName, ok := storageAccountNameRaw.(string)
			if !ok {
				return rawState, fmt.Errorf("expected `storage_account_name` to be of type string, got %T", storageAccountNameRaw)
			}

			// This may seem like an excessive timeout, however, populating the accounts via the list API could take a significant amount of time
			// when the subscription contains a large number of accounts and the account cache is not already populated.
			findCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
			defer cancel()

			log.Printf("[DEBUG] searching for a storage account by name (`%s`) in subscription (`%s`)", storageAccountName, subscriptionID)
			account, err := client.FindAccount(findCtx, subscriptionID, storageAccountName)
			if err != nil || account == nil {
				return rawState, fmt.Errorf("locating a storage account by name (`%s`) in subscription (`%s`): %w", storageAccountName, subscriptionID, err)
			}

			storageAccountID = account.StorageAccountId
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
