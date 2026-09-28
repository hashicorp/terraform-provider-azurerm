// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	legacyTables "github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

var _ pluginsdk.StateUpgrade = StorageTableEntityV0ToV1{}

type StorageTableEntityV0ToV1 struct{}

func (StorageTableEntityV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_table_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"partition_key": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"row_key": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"entity": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (StorageTableEntityV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawStorageTableID, ok := rawState["storage_table_id"].(string)
		if !ok {
			return rawState, fmt.Errorf("expected `storage_table_id` to be of type string, got %T", rawState["storage_table_id"])
		}

		// Some instances of this resource may already be using the resource manager ID, in which case this is a no-op
		if _, err := tables.ParseTableID(rawStorageTableID); err == nil {
			return rawState, nil
		}

		client := meta.(*clients.Client).Storage
		subscriptionID := meta.(*clients.Client).Account.SubscriptionId

		legacyStorageTableID, err := legacyTables.ParseTableID(rawStorageTableID, client.StorageDomainSuffix)
		if err != nil {
			return rawState, err
		}

		// This may seem like an excessive timeout, however, populating the accounts via the list API could take a significant amount of time
		// when the subscription contains a large number of accounts and the account cache is not already populated.
		findCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
		defer cancel()

		log.Printf("[DEBUG] searching for a storage account by name (`%s`) in subscription (`%s`)", legacyStorageTableID.AccountId.AccountName, subscriptionID)
		account, err := client.FindAccount(findCtx, subscriptionID, legacyStorageTableID.AccountId.AccountName)
		if err != nil || account == nil {
			return rawState, fmt.Errorf("locating a storage account by name (`%s`) in subscription (`%s`): %w", legacyStorageTableID.AccountId.AccountName, subscriptionID, err)
		}

		rawState["storage_table_id"] = tables.NewTableID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, legacyStorageTableID.TableName).ID()

		return rawState, nil
	}
}
